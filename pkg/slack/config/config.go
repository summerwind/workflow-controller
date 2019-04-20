package config

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

const (
	defaultConfigPath = "/etc/slack-controller/config.yaml"
)

type Config struct {
	Channels map[string]string `json:"channels"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("SLACK_CONTROLLER_CONFIG")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
