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
	"io"
	"log"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
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
