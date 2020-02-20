package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	ID                      string `json:"id"`
	URL                     string `json:"url"`
	Username                string `json:"username"`
	APIKey                  string `json:"api_key"`
	AuthKey                 string `json:"auth_key"`
	GroupID                 string `json:"group"`
	PlanID                  string `json:"plan"`
	NodeAddresses           string `json:"nodes"`
	AdminPassword           string `json:"admin_password"`
	AutomationAgentPassword string `json:"auth_pwd"`
	EngineVersion           string `json:"engine_version"`
	Routers                 int    `json:"routers"`
	ConfigServers           int    `json:"config_servers"`
	Replicas                int    `json:"replicas"`
	RequireSSL              bool   `json:"require_ssl"`
}

func LoadConfig(configFile string) (config *Config, err error) {
	if configFile == "" {
		return config, errors.New("Must provide a config file")
	}

	file, err := os.Open(configFile)
	if err != nil {
		return config, errors.Wrap(err, "cannot open config")
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return config, errors.Wrap(err, "cannot read config")
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return config, errors.Wrap(err, "cannot unmarshal config JSON")
	}

	if err = config.Validate(); err != nil {
		return config, errors.Wrap(err, "cannot validate config")
	}

	return config, nil
}

func (c Config) Validate() error {
	if c.ID == "" {
		return errors.New("Must provide a non-empty ID")
	}

	if c.URL == "" {
		return errors.New("Must provide a non-empty URL")
	}

	if c.Username == "" {
		return errors.New("Must provide a non-empty Username")
	}

	if c.APIKey == "" {
		return errors.New("Must provide a non-empty API Key")
	}

	if c.GroupID == "" {
		return errors.New("Must provide a non-empty Group ID")
	}

	if c.PlanID == "" {
		return errors.New("Must provide a non-empty Plan ID")
	}

	if c.NodeAddresses == "" {
		return errors.New("Must provide a non-empty Node Addresses")
	}

	if c.AdminPassword == "" {
		return errors.New("Must provide a non-empty Admin Password")
	}

	if c.EngineVersion == "" {
		return errors.New("Must provide a non-empty Engine Version")
	}

	return nil
}
