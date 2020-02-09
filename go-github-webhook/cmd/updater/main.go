package main

import (
	"os"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
)

func main() {
	config, storage, buildCtx := updater.EntryPoint()
	if config == nil {
		os.Exit(1)
	}
	updater.BuildLoop(config, storage, buildCtx)
}
