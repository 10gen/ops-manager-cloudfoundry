package mongodb

import (

	// helpersCF "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	// "github.com/onsi/gomega/gexec"
	//CF "smoke-tests/cf"
	"mongodb-service-adapter/adapter"
)

type OpsManager struct {
	BaseUrl string
}
type Service struct {
	Name          string
	Service       string
	Tags          string
	Plan          string
	Description   string
	Documentation string
	Dashboard     string
	ServiceBroker string
	Status        string
	Message       string
	Started       string
	Update        string
}

func GetConfig() *adapter.Config {
	config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")
	if err != nil {
		return nil
	}
	return config
}

var (
	config = GetConfig()
	c      = &adapter.OMClient{
		Url:      config.URL,
		Username: config.Username,
		ApiKey:   config.APIKey,
	}
)
