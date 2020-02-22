package server

import (
	"flag"
	"fmt"
	"log"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
)

func EntryPoint() (*utils.AppConfig, *updater.AppConfig, *updater.HashStorage, *updater.BuildContext) {
	var configFile string
	flag.StringVar(&configFile, "config", "./config.json", "config file")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Missing working directory.")
		flag.Usage()
		return nil, nil, nil, nil
	}
	if configFile == "" {
		flag.Usage()
		return nil, nil, nil, nil
	}
	workDir := flag.Arg(0)

	// load config
	config, err := utils.LoadConfig(configFile)
	if err != nil {
		log.Fatalln("ERROR", "loading config", err)
		return nil, nil, nil, nil
	}

	updaterConfig, hashStorage, buildCtx := updater.GetConfigs(workDir, *config.UpdaterConfig, "", "", "", "", "", "", false)
	if updaterConfig == nil {
		return nil, nil, nil, nil
	}
	return config, updaterConfig, hashStorage, buildCtx
}
