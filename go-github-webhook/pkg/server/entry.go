package server

import (
	"flag"
	"fmt"
	"log"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/structs"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
)

func EntryPoint() (*utils.Config, *updater.Config, *structs.HashStorage, string) {
	var configFile string
	flag.StringVar(&configFile, "config", "./config.json", "config file")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Missing working directory.")
		flag.Usage()
		return nil, nil, nil, ""
	}
	if configFile == "" {
		flag.Usage()
		return nil, nil, nil, ""
	}
	workDir := flag.Arg(0)

	// load config
	config, err := utils.LoadConfig(configFile)
	if err != nil {
		log.Fatalln("ERROR", "loading config", err)
		return nil, nil, nil, ""
	}

	updaterConfig := updater.GetConfigs(workDir, *config.UpdaterConfig, "", "", "", "", "", "", "", "", false)
	if updaterConfig == nil {
		log.Fatalln("ERROR", "loading updater config", err)
		return nil, nil, nil, ""
	}
	hashStorage := updater.GetHashStorage(updaterConfig, workDir)
	if hashStorage == nil {
		log.Fatalln("ERROR", "loading hash storage", err)
		return nil, nil, nil, ""
	}
	if !hashStorage.IsLocal() {
		log.Println("ignoring online hash storage, use local one instead")
		hashStorage = structs.NewLocalStorage(*config.UpdaterConfig)
		if err := hashStorage.LoadFromFile(); err != nil {
			log.Fatalln("ERROR", "could not load local hash storage", err)
			return nil, nil, nil, ""
		}
	}
	return config, updaterConfig, hashStorage, workDir
}
