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

func dockerBuild(workDir string, dockerfile string, branch string, arch string, reCommit string) bool {
	repo := "iceflower/redeclipse-server"

	success := runCmd("docker", []string{"build", "--build-arg", "BRANCH=" + branch, "--build-arg", "RECOMMIT=" + reCommit, "-t", repo + ":" + arch + "-" + branch, "-f", dockerfile, workDir}, nil)
	if !success {
		return false
	}
	success = runCmd("docker", []string{"push", repo + ":" + arch + "-" + branch}, nil)
	if !success {
		return false
	}
	repoBranch := repo + ":" + branch
	runCmd("docker", []string{"manifest", "create", repoBranch, repo + ":amd64-" + branch, repo + ":arm64-" + branch}, nil)
	runCmd("docker", []string{"manifest", "annotate", repoBranch, repo + ":arm64-" + branch, "--variant", "v8"}, nil)
	runCmd("docker", []string{"manifest", "push", "--purge", repoBranch}, nil)
	return true
}
