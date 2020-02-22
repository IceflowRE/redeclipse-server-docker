package updater

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

func dockerBuild(workDir string, repo string, dockerfile string, ref string, arch string, reCommit string) bool {
	dockerTag := ref[strings.LastIndex(ref, "/")+1:]
	dockerName := repo + ":" + dockerTag
	dockerArchName := repo + ":" + arch + "-" + dockerTag


	success := runCmd("docker", []string{"build", "--build-arg", "TAG=" + dockerTag, "--build-arg", "RECOMMIT=" + reCommit, "-t", dockerArchName, "-f", dockerfile, workDir}, nil)
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
