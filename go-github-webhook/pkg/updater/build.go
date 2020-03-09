package updater

import (
	"fmt"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/structs"
)

func Build(config *Config, storage *structs.HashStorage, build *BuildConfig, workDir string) bool {
	fmt.Println("Update step: " + build.Ref + " - " + build.Arch + " - " + build.Os + " - " + build.Dockerfile)
	curHash, err := storage.Get(build.Ref, build.Arch, build.Os)
	if err != nil {
		fmt.Println("failed to get current hashes:", err.Error())
		return false
	}
	newHash := getNewHashes(build.Dockerfile, build.Ref, build.Arch, build.Os)
	if newHash == nil {
		fmt.Println("failed to get new hashes")
		return false
	}

	if curHash == nil {
		fmt.Println("Current: alpine: - | dockerfile: - | recommit: -")
	} else {
		fmt.Println("Current:", "alpine:", curHash.Alpine, "| dockerfile:", curHash.Dockerfile, "| recommit:", curHash.ReCommit)
	}
	fmt.Println("New:    ", "alpine:", newHash.Alpine, "| dockerfile:", newHash.Dockerfile, "| recommit:", newHash.ReCommit)
	if curHash != nil && *curHash == *newHash {
		fmt.Println("No update required.")
		return true
	}
	fmt.Println("Update required.")

	if config.DryRun {
		fmt.Println("Dry run, stop here.")
		return true
	}

	if !dockerLogin(config.Docker.User, config.Docker.Password) {
		fmt.Println("failed to docker login")
		return false
	}
	defer dockerLogout()
	if !dockerBuild(workDir, config.Docker.Repo, build.Dockerfile, build.Ref, build.Arch, newHash.ReCommit) {
		return false
	}

	if err = storage.Update(build.Ref, build.Arch, build.Os, newHash); err != nil {
		fmt.Println("FAILED to save new hash values:", err.Error())
		return false
	}
	if err = storage.SaveToFile(); err != nil {
		fmt.Println("FAILED to save new hash values:", err.Error())
		return false
	}
	fmt.Println("Saved new hash values.")
	return true
}

func BuildLoop(config *Config, storage *structs.HashStorage, workDir string) bool {
	success := true
	for _, build := range config.Build {
		if !Build(config, storage, build, workDir) {
			success = false
		}
	}
	return success
}
