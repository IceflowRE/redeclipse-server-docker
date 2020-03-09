package updater

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/structs"
)


func GetConfigs(workDir string, configFile string, repo string, user string, password string, ref string, arch string, dockerfile string, hashUrl string, apiKey string, dryRun bool) *Config {
	if configFile != "" {
		configFile = filepath.Join(workDir, configFile)
	}
	switch {
	case !structs.DirectoryExists(workDir):
		fmt.Println("working directory '" + workDir + "' does not exist")
		return nil
	case configFile != "" && !structs.FileExists(configFile):
		fmt.Println("config file '" + configFile + "' does not exist")
		return nil
	}

	var config *Config
	if configFile == "" {
		config = newConfig(repo, user, password, ref, arch, dockerfile, hashUrl, apiKey, dryRun)
	} else {
		var err error
		if config, err = LoadConfig(configFile); err != nil {
			log.Println(err.Error())
			return nil
		}
	}
	// if specified override config values
	if repo != "" {
		config.Docker.Repo = repo
	}
	if user != "" {
		config.Docker.User = user
	}
	if password != "" {
		config.Docker.Password = password
	}
	if apiKey != "" {
		config.HashApi.ApiKey = apiKey
	}
	if err := config.CheckDockerfiles(workDir); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return config
}

func GetHashStorage(config *Config, workDir string) *structs.HashStorage {
	hashStorageFile := filepath.Join(workDir, "hash-storage.json")
	var hashStorage *structs.HashStorage
	if config.HashApi == nil || config.HashApi.Url == "" {
		if !structs.FileExists(hashStorageFile) {
			fmt.Println("hash storage '" + hashStorageFile + "' does not exist, will use empty one")
		}
		hashStorage = structs.NewLocalStorage(hashStorageFile)
		if err := hashStorage.LoadFromFile(); err != nil {
			fmt.Println(err.Error())
			return nil
		}
	} else {
		hashStorage = structs.NewOnlineStorage(config.HashApi.Url, config.HashApi.ApiKey)
	}
	return hashStorage
}

func EntryPoint() (config *Config, hashStorage *structs.HashStorage, workDir string) {
	var dryRun bool
	flag.BoolVar(&dryRun, "dry", false, "dry run, just shows which images are outdated")
	var configFile string
	flag.StringVar(&configFile, "config", "", "config file")
	var dockerFile string
	flag.StringVar(&dockerFile, "dockerfile", "", "name of the dockerfile in the working directory")
	var ref string
	flag.StringVar(&ref, "ref", "", "git reference")
	var arch string
	flag.StringVar(&arch, "arch", "", "cpu architecture")
	var repo string
	flag.StringVar(&repo, "repo", "", "docker repository")
	var user string
	flag.StringVar(&user, "user", "", "docker user, overrides config value")
	var password string
	flag.StringVar(&password, "password", "", "docker password, overrides config value")
	var hashUrl string
	flag.StringVar(&hashUrl, "hashUrl", "", "docker repository")
	var apiKey string
	flag.StringVar(&apiKey, "apiKey", "", "apiKey of hash url, overrides config value")

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Missing working directory.")
		flag.Usage()
		return nil, nil, ""
	}
	if configFile == "" && (ref == "" || arch == "" || dockerFile == "" || repo == "" || user == "" || password == "") && (hashUrl != "" && apiKey == "" || hashUrl == "" && apiKey != "") {
		flag.Usage()
		return nil, nil, ""
	}
	workDir = flag.Arg(0)

	config = GetConfigs(workDir, configFile, repo, user, password, ref, arch, dockerFile, hashUrl, apiKey, dryRun)
	hashStorage = GetHashStorage(config, workDir)
	if hashStorage.IsLocal() {
		fmt.Println("hash storage: local")
	} else {
		fmt.Println("ash storage:", config.HashApi.Url)
	}
	return config, hashStorage, workDir
}
