// Copyright 2019 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
