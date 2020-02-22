package updater

import (
	"flag"
	"fmt"
)

func EntryPoint() (*AppConfig, *HashStorage, *BuildContext) {
	var useDryRun bool
	flag.BoolVar(&useDryRun, "dry", false, "dry run, just shows which images are outdated")
	var configFile string
	flag.StringVar(&configFile, "config", "", "config file")
	var dockerFile string
	flag.StringVar(&dockerFile, "dockerfile", "", "name of the dockerfile in the working directory")
	var ref string
	flag.StringVar(&ref, "ref", "", "git reference")
	var arch string
	flag.StringVar(&arch, "arch", "", "cpu architecture")
	var repo string
	flag.StringVar(&repo, "repo", "", "docker repository")
	var user string
	flag.StringVar(&user, "user", "", "docker user, overrides config value")
	var password string
	flag.StringVar(&password, "password", "", "docker password, overrides config value")

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Missing working directory.")
		flag.Usage()
		return nil, nil, nil
	}
	if configFile == "" && (ref == "" || arch == "" || dockerFile == "" || repo == "" || user == "" || password == "") {
		flag.Usage()
		return nil, nil, nil
	}
	workDir := flag.Arg(0)

	return GetConfigs(workDir, configFile, repo, user, password, ref, arch, dockerFile, useDryRun)
}
