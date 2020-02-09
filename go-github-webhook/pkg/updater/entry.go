package updater

import (
	"flag"
	"fmt"
)

func EntryPoint() (*AppConfig, *HashStorage, *BuildContext) {
	var useDryRun bool
	flag.BoolVar(&useDryRun, "dry", false, "dry run, just shows which images are outdated")
	var configFile string
	flag.StringVar(&configFile, "config", "", "use of config file")
	var dockerFile string
	flag.StringVar(&dockerFile, "dockerfile", "", "name of the dockerfile in the working directory, optional, default to \"Dockerfile_<branch>\"")
	var branch string
	flag.StringVar(&branch, "branch", "", "git branch")
	var arch string
	flag.StringVar(&arch, "arch", "", "cpu architecture")
	var user string
	flag.StringVar(&user, "user", "", "docker user")
	var password string
	flag.StringVar(&password, "password", "", "docker password")

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Missing working directory.")
		flag.Usage()
		return nil, nil, nil
	}
	if configFile == "" && (branch == "" || arch == "" || user == "" || password == "") {
		flag.Usage()
		return nil, nil, nil
	}
	workDir := flag.Arg(0)

	return GetConfigs(workDir, configFile, user, password, branch, arch, dockerFile, useDryRun)
}
