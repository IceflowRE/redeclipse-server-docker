package updater

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
)

func GetConfigs(workDir string, configFile string, repo string, user string, password string, ref string, arch string, dockerfile string, dryRun bool) *Config {
	if configFile != "" {
		configFile = filepath.Join(workDir, configFile)
	}
	switch {
	case !DirectoryExists(workDir):
		fmt.Println("working directory '" + workDir + "' does not exist")
		return nil
	case configFile != "" && !FileExists(configFile):
		fmt.Println("config file '" + configFile + "' does not exist")
		return nil
	}

	var config *Config
	if configFile == "" {
		config = newConfig(repo, user, password, ref, arch, dockerfile, dryRun)
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
	if err := config.CheckDockerfiles(workDir); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return config
}

func EntryPoint() (config *Config, workDir string) {
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

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Missing working directory.")
		flag.Usage()
		return nil, ""
	}
	if configFile == "" && (ref == "" || arch == "" || dockerFile == "" || repo == "" || user == "" || password == "") {
		flag.Usage()
		return nil, ""
	}
	workDir = flag.Arg(0)

	config = GetConfigs(workDir, configFile, repo, user, password, ref, arch, dockerFile, dryRun)
	return config, workDir
}
