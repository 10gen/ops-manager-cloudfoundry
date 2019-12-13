package opsmanager

import (
	"bytes"
	"encoding/json"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
)

// UpdateAutomationConfig
// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/#update-the-automation-configuration
func (client opsManagerClient) UpdateAutomationConfig(projectID string, config AutomationConfig) (AutomationConfig, error) {
	var result AutomationConfig

	bodyBytes, err := json.Marshal(config)
	if err != nil {
		return result, err
	}

	url := client.resolver.Of("/groups/%s/automationConfig", projectID)
	resp := client.PutJSON(url, bytes.NewReader(bodyBytes))
	if resp.IsError() {
		return result, resp.Err
	}
	defer httpclient.CloseResponseBodyIfNotNil(resp)

	decoder := json.NewDecoder(resp.Response.Body)
	err2 := decoder.Decode(&result)
	useful.PanicOnUnrecoverableError(err2)

	return result, nil
}
