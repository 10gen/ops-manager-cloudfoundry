package service_test

import (
	"fmt"
)

//ServiceParameters is the model of req. parameters for creating services
type ServiceParameters struct {
	ServiceName  string
	PlanName     string
	BackupEnable string
	SSLEnable    string
}

//generate different service config for test
func generateTestServiceConfigs(testConfig mongodbTestConfig) []ServiceParameters {
	var testServiceConfigs []ServiceParameters
	for _, planName := range testConfig.PlanNames {
		for _, backup := range testConfig.Backup {
			for _, ssl := range testConfig.SSL {
				testServiceConfigs = append(testServiceConfigs,
					ServiceParameters{testConfig.ServiceName, planName, backup, ssl})
			}
		}
	}
	return testServiceConfigs
}

//print result of generateServiceConfigs
func printGeneratedServiceConfigs(testConfig []ServiceParameters) {
	for i, element := range testConfig {
		fmt.Printf("==== Case %d ====\n", i+1)
		fmt.Printf("Service Name %s \n", element.ServiceName)
		fmt.Printf("PlanName is %s \n", element.PlanName)
		fmt.Printf("Backup %s \n", element.BackupEnable)
		fmt.Printf("SSL enabled %s \n", element.SSLEnable)
		fmt.Println("=================")
	}
}
