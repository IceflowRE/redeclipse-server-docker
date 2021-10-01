package main

import (
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
	"os"
)

func main() {
	config, workDir := updater.EntryPoint()
	if config == nil {
		os.Exit(1)
	}
	if !updater.BuildLoop(config, workDir) {
		os.Exit(2)
	}
}
