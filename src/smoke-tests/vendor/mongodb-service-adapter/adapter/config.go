package adapter

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	URL           string `json:"url"`
	Username      string `json:"username"`
	APIKey        string `json:"api_key"`
	AuthKey       string `json:"auth_key"`
	GroupID       string `json:"group"`
	NodeAddresses string `json:"nodes"`
	OrgID         string `json:"org"`
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

	if c.NodeAddresses == "" {
		return errors.New("Must provide a non-empty Node Addresses")
	}

	if c.OrgID == "" {
		return errors.New("Must provide a non-empty Org ID")
	}

	if c.AuthKey == "" {
		return errors.New("Must provide a non-empty Auth Key")
	}

	return nil
}
