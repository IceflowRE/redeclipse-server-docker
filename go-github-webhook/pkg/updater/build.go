package updater

import (
	"fmt"
)

type BuildContext struct {
	WorkDir         string
	HashStorageFile string
	DryRun          bool
}

func BuildStep(config *AppConfig, storage *HashStorage, buildCtx *BuildContext, build *buildConfig) {
	fmt.Println("Update step: " + build.Branch + " - " + build.Arch + " - " + build.Os + " - " + build.Dockerfile)
	curHash := storage.Get(build.Branch, build.Arch, build.Os)
	newHash := getNewHashes(build.Dockerfile, build.Branch, build.Arch, build.Os)
	if newHash == nil {
		fmt.Println("failed to get new hashes")
		return
	}
	fmt.Println("Current:", "alpine:", curHash.Alpine, "- dockerfile:", curHash.Dockerfile, "- recommit:", curHash.ReCommit)
	fmt.Println("New:    ", "alpine:", newHash.Alpine, "- dockerfile:", newHash.Dockerfile, "- recommit:", newHash.ReCommit)

	if curHash == newHash {
		fmt.Println("No update required.")
		return
	}
	fmt.Println("Update required.")

	if buildCtx.DryRun {
		fmt.Println("Dry run, stop here.")
		return
	}

	if !dockerLogin(config.Docker.User, config.Docker.Password) {
		return
	}
	defer dockerLogout()
	if !dockerBuild(buildCtx.WorkDir, build.Dockerfile, build.Branch, build.Arch, newHash.ReCommit) {
		return
	}
	storage.Update(build.Branch, build.Arch, build.Os, newHash)
	err := saveHashStorage(buildCtx.HashStorageFile, storage)
	if err != nil {
		fmt.Println("FAILED to save new hash values.")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Saved new hash values.")
}

func BuildLoop(config *AppConfig, storage *HashStorage, buildCtx *BuildContext) {
	for _, build := range config.Build {
		BuildStep(config, storage, buildCtx, build)
	}
}
