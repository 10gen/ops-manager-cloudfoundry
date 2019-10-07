package opsmanager

import (
	"encoding/json"
	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
)

// Process process status

// AutomationStatusResponse automation status
type AutomationStatusResponse struct {
	Processes   []Process `json:"processes"`
	GoalVersion int       `json:"goalVersion"`
}

// GetAutomationStatus
// https://docs.opsmanager.mongodb.com/master/reference/api/automation-status/
func (client opsManagerClient) GetAutomationStatus(projectID string) (AutomationStatusResponse, error) {
	var result AutomationStatusResponse

	url := client.resolver.Of("/groups/%s/automationStatus", projectID)
	resp := client.GetJSON(url)
	if resp.IsError() {
		return result, resp.Err
	}
	defer httpclient.CloseResponseBodyIfNotNil(resp)

	decoder := json.NewDecoder(resp.Response.Body)
	err := decoder.Decode(&result)
	useful.PanicOnUnrecoverableError(err)

	return result, nil
}
