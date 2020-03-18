package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

//TestMatrix matrix for generation test-cases environments for service suite
type TestMatrix struct {
	ServiceName    string         `json:"service_name"`
	PlanNames      []string       `json:"plan_names"`
	Retry          retryConfig    `json:"retry"`
	Backup         []string       `json:"backup_enabled"`
	SSL            []string       `json:"ssl_enabled"`
	MongoDBVersion []string       `json:"mongodb_version"`
	OpsMan         OpsManagerCred `json:"mongo_ops"`
}

//SetDefaultForNonDefinedParameters test if all parameters in pipeline was provided
func (testConfig *TestMatrix) SetDefaultForNonDefinedParameters() { //TODO recheck
	if testConfig.ServiceName == "" {
		testConfig.ServiceName = "mongodb-odb"
	}
	if testConfig.PlanNames == nil || len(testConfig.PlanNames) == 0 {
		testConfig.PlanNames = []string{"standalone_small", "replica_set_small", "sharded_cluster_small"}
	}
	if testConfig.Backup == nil || len(testConfig.Backup) == 0 {
		testConfig.Backup = []string{"false", "true"}
	}
	if testConfig.SSL == nil || len(testConfig.SSL) == 0 {
		testConfig.SSL = []string{"false"}
	}
	if testConfig.MongoDBVersion == nil || len(testConfig.MongoDBVersion) == 0 {
		testConfig.MongoDBVersion = []string{"4.0.9-ent"}
	}
}

//GenerateTestServiceParameters generate cases of service parameters from testMatrix
func GenerateTestServiceParameters(testConfig TestMatrix) []ServiceParameters {
	var (
		testServiceParameters []ServiceParameters
		names                 InstanceNames
	)
	names.GenerateInstanceNames()
	for _, planName := range testConfig.PlanNames {
		for _, backup := range testConfig.Backup {
			for _, ssl := range testConfig.SSL {
				for _, version := range testConfig.MongoDBVersion {
					testServiceParameters = append(testServiceParameters,
						ServiceParameters{names, testConfig.ServiceName, planName, "", backup, ssl, version, "", "", "", "", ""})
				}
			}
		}
	}
	return testServiceParameters
}

//PrintGeneratedServiceParameters print result of generateServiceParameters
func PrintGeneratedServiceParameters(testConfig []ServiceParameters) {
	for i, element := range testConfig {
		fmt.Printf("==== Case %d ====\n", i+1)
		fmt.Printf("Service Name: %s \n", element.MarketPlaceServiceName)
		fmt.Printf("PlanName: %s \n", element.PlanName)
		fmt.Printf("Backup: %s \n", element.BackupEnable)
		fmt.Printf("SSL enabled: %s \n", element.SSLEnable)
		fmt.Printf("Version MongoDB: %s \n", element.MongoDBVersion)
		fmt.Println("=================")
	}
}

//LoadMongodbTestConfig load test matrix from file
func LoadMongodbTestConfig(path string) TestMatrix {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	config := TestMatrix{}
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		panic(err)
	}

	return config
}
