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

// CreateOneProject
// https://docs.opsmanager.mongodb.com/master/reference/api/groups/create-one-group/
func (client opsManagerClient) CreateOneProject(name string, orgID string) (ProjectResponse, error) {
	var result ProjectResponse

	// create request object
	request := make(map[string]string)
	request["name"] = name
	if orgID != "" {
		// optional parameter
		request["orgId"] = orgID
	}

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return result, err
	}

	url := client.resolver.Of("/groups")
	resp := client.PostJSON(url, bytes.NewReader(bodyBytes))
	if resp.IsError() {
		return result, resp.Err
	}
	defer httpclient.CloseResponseBodyIfNotNil(resp)

	decoder := json.NewDecoder(resp.Response.Body)
	err = decoder.Decode(&result)
	useful.PanicOnUnrecoverableError(err)

	return result, nil
}
