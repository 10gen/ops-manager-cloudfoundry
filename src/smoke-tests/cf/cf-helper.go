package cf

import(
	"github.com/pivotal-cf-experimental/cf-test-helpers/services"
)

//LoadCFTestConfig use pivotal functions for getting configuration
func LoadCFTestConfig(path string) services.Config {
	config := services.Config{}

	if err := services.LoadConfig(path, &config); err != nil {
		panic(err)
	}

	if err := services.ValidateConfig(&config); err != nil {
		panic(err)
	}

	config.TimeoutScale = 3

	return config
}