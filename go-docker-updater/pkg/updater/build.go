package updater

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Build(config *Config, build *BuildConfig, workDir string) bool {
	fmt.Println("Update step: " + build.Ref + " - " + build.Arch + " - " + build.Os + " - " + build.Dockerfile)
	curHash, err := GetCurrentHashes(config.Docker.Repo, build.Ref, build.Arch)
	if err != nil {
		fmt.Println("failed to get current hashes")
		return false
	}
	newHash := GetNewHashes(build.Dockerfile, build.Ref, build.Arch, build.Os)
	if newHash == nil {
		fmt.Println("failed to get new hashes")
		return false
	}
	
	fmt.Println("Current:", "alpine:", curHash.Alpine, "| dockerfile:", curHash.Dockerfile, "| re-commit:", curHash.ReCommit)
	fmt.Println("New:    ", "alpine:", newHash.Alpine, "| dockerfile:", newHash.Dockerfile, "| re-commit:", newHash.ReCommit)

	if *curHash == *newHash {
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
	if !dockerBuild(workDir, config.Docker.Repo, build, newHash) {
		return false
	}
	return true
}

func BuildLoop(config *Config, workDir string) bool {
	success := true
	for _, build := range config.Build {
		if !Build(config, build, workDir) {
			success = false
		}
	}
	return success
}

func runCmd(name string, arg []string, stdin *string) bool {
	fmt.Println("RUN: " + name + " " + strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if stdin != nil {
		cmd.Stdin = strings.NewReader(*stdin)
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("FAILED: " + name + " " + strings.Join(arg, " "))
		fmt.Println(err.Error())
		return false
	}
	return true
}

func dockerLogin(user string, password string) bool {
	return runCmd("docker", []string{"login", "--username", user, "--password-stdin"}, &password)
}

func dockerLogout() bool {
	return runCmd("docker", []string{"logout"}, nil)
}

func dockerBuild(workDir string, repo string, config *BuildConfig, hash *Hash) bool {
	dockerTag := DockerTagFromRef(config.Ref)
	dockerName := repo + ":" + dockerTag
	dockerArchName := repo + ":" + config.Arch + "-" + dockerTag

	success := runCmd("docker", []string{"build",
		"--build-arg", "TAG=" + dockerTag,
		"--build-arg", "RE_COMMIT=" + hash.ReCommit,
		"--build-arg", "ALPINE_SHA=" + hash.Alpine,
		"--build-arg", "DOCKERFILE_SHA=" + hash.Dockerfile,
		"-t", dockerArchName, "-f", config.Dockerfile, workDir}, nil)
	if !success {
		return false
	}
	success = runCmd("docker", []string{"push", dockerArchName}, nil)
	if !success {
		return false
	}
	runCmd("docker", []string{"manifest", "create", dockerName, repo + ":amd64-" + dockerTag, repo + ":arm64-" + dockerTag}, nil)
	runCmd("docker", []string{"manifest", "annotate", dockerName, repo + ":arm64-" + dockerTag, "--variant", "v8"}, nil)
	runCmd("docker", []string{"manifest", "push", "--purge", dockerName}, nil)
	return true
}
