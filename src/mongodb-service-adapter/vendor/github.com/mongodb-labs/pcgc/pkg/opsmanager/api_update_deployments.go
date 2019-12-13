package opsmanager

import (
	"bytes"
	"encoding/json"
	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
	"io"
	"log"
)

// UpdateAutomationConfig
// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/#update-the-automation-configuration
func (client opsManagerClient) UpdateDeployments(projectID string, body io.Reader) (map[string]interface{}, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := client.resolver.Of("/groups/%s/automationConfig", projectID)
	resp := client.PutJSON(url, bytes.NewReader(bodyBytes))
	if resp.IsError() {
		return nil, resp.Err
	}
	defer httpclient.CloseResponseBodyIfNotNil(resp)

	result := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Response.Body)
	err2 := decoder.Decode(&result)
	useful.PanicOnUnrecoverableError(err2)

	return result, nil
}

// AddProcessToDeployment adds the specified Process to the processes list of the given deployment
func AddProcessToDeployment(process Process, automationConfig map[string]interface{}) map[string]interface{} {
	// TODO(mihaibojin): use mergo to merge these together
	log.Fatalln("not implemented")
	return nil
}
