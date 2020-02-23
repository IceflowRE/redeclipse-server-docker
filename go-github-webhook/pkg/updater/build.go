package updater

import (
	"fmt"
)

type BuildContext struct {
	WorkDir         string
	HashStorageFile string
	DryRun          bool
}

func BuildStep(config *AppConfig, storage *HashStorage, buildCtx *BuildContext, build *buildConfig) bool {
	fmt.Println("Update step: " + build.Ref + " - " + build.Arch + " - " + build.Os + " - " + build.Dockerfile)
	curHash := storage.Get(build.Ref, build.Arch, build.Os)
	newHash := getNewHashes(build.Dockerfile, build.Ref, build.Arch, build.Os)
	if newHash == nil {
		fmt.Println("failed to get new hashes")
		return false
	}
	fmt.Println("Current:", "alpine:", curHash.Alpine, "- dockerfile:", curHash.Dockerfile, "- recommit:", curHash.ReCommit)
	fmt.Println("New:    ", "alpine:", newHash.Alpine, "- dockerfile:", newHash.Dockerfile, "- recommit:", newHash.ReCommit)

	if *curHash == *newHash {
		fmt.Println("No update required.")
		return true
	}
	fmt.Println("Update required.")

	if buildCtx.DryRun {
		fmt.Println("Dry run, stop here.")
		return true
	}

	if !dockerLogin(config.Docker.User, config.Docker.Password) {
		fmt.Println("failed to docker login")
		return false
	}
	defer dockerLogout()
	if !dockerBuild(buildCtx.WorkDir, config.Docker.Repo, build.Dockerfile, build.Ref, build.Arch, newHash.ReCommit) {
		return false
	}
	storage.Update(build.Ref, build.Arch, build.Os, newHash)
	err := saveHashStorage(buildCtx.HashStorageFile, storage)
	if err != nil {
		fmt.Println("FAILED to save new hash values.")
		return false
	}
	fmt.Println("Saved new hash values.")
	return true
}

func BuildLoop(config *AppConfig, storage *HashStorage, buildCtx *BuildContext) bool {
	success := true
	for _, build := range config.Build {
		if !BuildStep(config, storage, buildCtx, build) {
			success = false
		}
	}
	return success
}
