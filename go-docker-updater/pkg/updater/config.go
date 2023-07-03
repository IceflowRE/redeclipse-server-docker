package updater

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Docker *dockerConfig  `json:"docker"`
	Build  []*BuildConfig `json:"build"`
	DryRun bool           `json:"dryRun"`
}

type dockerConfig struct {
	Repo     string `json:"repo"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type BuildConfig struct {
	Ref        string `json:"ref"`
	Arch       string `json:"arch"`
	Os         string `json:"os"`
	Dockerfile string `json:"dockerfile"`
}

func newConfig(repo string, user string, password string, ref string, arch string, dockerfile string, dryRun bool) *Config {
	return &Config{
		Docker: &dockerConfig{
			Repo:     repo,
			User:     user,
			Password: password,
		},
		Build: []*BuildConfig{
			{
				Ref:        ref,
				Arch:       arch,
				Os:         "linux",
				Dockerfile: dockerfile,
			},
		},
		DryRun: dryRun,
	}
}

func (config *Config) check() error {
	switch {
	case config.Docker == nil:
		return errors.New("config: 'docker' is missing")
	case config.Docker.Repo == "":
		return errors.New("config: 'docker > repo' is empty")
	case config.Docker.User == "":
		return errors.New("config: 'docker > user' is empty")
	case config.Docker.Password == "":
		return errors.New("config: 'docker > password' is empty")
	case config.Build == nil:
		return errors.New("config: 'build' is missing")
	case len(config.Build) == 0:
		return errors.New("config: 'build' is empty")
	}
	for idx, build := range config.Build {
		switch {
		case build.Ref == "":
			return errors.New("config: 'build[" + string(idx) + "] > ref' is empty")
		case build.Arch == "":
			return errors.New("config: 'build[" + string(idx) + "] > arch' is empty")
		case build.Os == "":
			return errors.New("config: 'build[" + string(idx) + "] > os' is empty")
		case build.Dockerfile == "":
			return errors.New("config: 'build[" + string(idx) + "] > dockerfile' is empty")
		}
	}
	return nil
}

func (config *Config) CheckDockerfiles(workDir string) error {
	for _, build := range config.Build {
		build.Dockerfile = filepath.Join(workDir, build.Dockerfile)
		if !FileExists(build.Dockerfile) {
			return errors.New("Dockerfile '" + build.Dockerfile + "' does not exist.")
		}
	}
	return nil
}

func LoadConfig(file string) (*Config, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(raw, &config)
	if err != nil {
		return nil, err
	}
	err = config.check()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (config *Config) Get(ref string) *BuildConfig {
	for _, build := range config.Build {
		if build.Ref == ref {
			return build
		}
	}
	return nil
}
