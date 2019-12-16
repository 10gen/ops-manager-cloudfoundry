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
	"bytes"
	"encoding/json"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
)

// CreateAgentAPIKEYResponse API response for the CreateAgentAPIKEY() call
type CreateAgentAPIKEYResponse struct {
	ID            string `json:"_id,omitempty"`
	CreatedBy     string `json:"createdBy,omitempty"`
	CreatedIPAddr string `json:"createdIpAddr,omitempty"`
	CreatedTime   int    `json:"createdTime,omitempty"`
	CreatedUserID string `json:"createdUserId,omitempty"`
	Desc          string `json:"desc,omitempty"`
	Key           string `json:"key,omitempty"`
}

// CreateAgentAPIKEY creates a new Agent API key in the specified project
// https://docs.opsmanager.mongodb.com/master/reference/api/agentapikeys/create-one-agent-api-key/
func (client opsManagerClient) CreateAgentAPIKEY(projectID string, desc string) (CreateAgentAPIKEYResponse, error) {
	var result CreateAgentAPIKEYResponse

	// create request
	request := make(map[string]string)
	request["desc"] = desc
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return result, err
	}

	url := client.resolver.Of("/groups/%s/agentapikeys", projectID)
	resp := client.PostJSON(url, bytes.NewReader(bodyBytes))
	if resp.IsError() {
		return result, resp.Err
	}
	defer httpclient.CloseResponseBodyIfNotNil(resp)

	decoder := json.NewDecoder(resp.Response.Body)
	err2 := decoder.Decode(&result)
	useful.PanicOnUnrecoverableError(err2)

	return result, nil
}
