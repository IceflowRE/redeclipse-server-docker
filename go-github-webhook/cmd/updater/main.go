package main

import (
	"os"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
)

func main() {
	config, storage, workDir := updater.EntryPoint()
	if config == nil {
		os.Exit(1)
	}
	if !updater.BuildLoop(config, storage, workDir) {
		os.Exit(2)
	}
}
