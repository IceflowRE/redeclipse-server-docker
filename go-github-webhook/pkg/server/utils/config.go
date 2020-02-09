package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type AppConfig struct {
	Port          *int    `json:"port,omitempty"`
	WebhookSecret *string `json:"webhookSecret,omitempty"`
	UpdaterConfig *string `json:"updaterConfig,omitempty"`
}

func LoadConfig(file string) (*AppConfig, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config AppConfig
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

func (config AppConfig) check() error {
	switch {
	case config.Port == nil:
		return errors.New("config: 'port' is missing")
	case config.WebhookSecret == nil:
		return errors.New("config: 'webhookSecret' is missing")
	case config.UpdaterConfig == nil:
		return errors.New("config: 'updaterConfig' is missing")
	}
	return nil
}
