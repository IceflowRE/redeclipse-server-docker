package updater

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CheckPaths(workDir string, hashStorageFile string, configFile string) (bool, bool, bool) {
	_, err := os.Stat(workDir)
	workDirExist := !os.IsNotExist(err)
	hashStoreFileExist := fileExists(hashStorageFile)
	configFileExist := configFile != "" && fileExists(configFile)
	return workDirExist, hashStoreFileExist, configFileExist
}

func GetConfigs(workDir string, configFile string, repo string, user string, password string, ref string, arch string, dockerFile string, useDryRun bool) (*AppConfig, *HashStorage, *BuildContext) {
	if configFile != "" {
		configFile = filepath.Join(workDir, configFile)
	}
	hashStorageFile := filepath.Join(workDir, "hash-storage.json")
	workDirExist, hashStorageFileExist, configFileExist := CheckPaths(workDir, hashStorageFile, configFile)
	switch {
	case !workDirExist:
		fmt.Println("working directory '" + workDir + "' does not exist")
		return nil, nil, nil
	case configFile != "" && !configFileExist:
		fmt.Println("config file '" + configFile + "' does not exist")
		return nil, nil, nil
	}
	if !hashStorageFileExist {
		fmt.Println("hash storage '" + hashStorageFile + "' does not exist, will use empty one")
	}

	var config *AppConfig
	if configFile != "" {
		var err error
		if config, err = parseConfig(configFile); err != nil {
			log.Println(err.Error())
			return nil, nil, nil
		}
	} else {
		config = &AppConfig{
			Docker: &dockerConfig{
				Repo:     repo,
				User:     user,
				Password: password,
			},
			Build: []*buildConfig{
				{
					Ref:     ref,
					Arch:       arch,
					Os:         "linux",
					Dockerfile: dockerFile,
				},
			},
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

	for _, build := range config.Build {
		if build.Dockerfile == "" {
			build.Dockerfile = "Dockerfile_" + build.Ref
		}
		build.Dockerfile = filepath.Join(workDir, build.Dockerfile)
		if !fileExists(build.Dockerfile) {
			fmt.Println("Dockerfile '" + build.Dockerfile + "' does not exist.")
			return nil, nil, nil
		}
	}

	hashStorage := NewHashStorage()
	if hashStorageFileExist {
		var err error
		if hashStorage, err = loadHashStorage(hashStorageFile); err != nil {
			fmt.Println(err.Error())
			return nil, nil, nil
		}
	}
	return config, hashStorage, &BuildContext{
		WorkDir:         workDir,
		HashStorageFile: hashStorageFile,
		DryRun:          useDryRun,
	}
}
