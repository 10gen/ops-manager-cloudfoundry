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

	"github.com/mongodb-labs/pcgc/pkg/useful"
)

func (client opsManagerClient) SetProjectTags(projectID string, tags []string) (ProjectResponse, error) {
	var result ProjectResponse

	url := client.resolver.Of("/groups/%s", projectID)
	request := map[string][]string{"tags": tags}

	body, err := json.Marshal(request)
	if err != nil {
		return result, err
	}

	resp := client.PatchJSON(url, bytes.NewReader(body))
	if resp.IsError() {
		return result, resp.Err
	}

	decoder := json.NewDecoder(resp.Response.Body)
	err = decoder.Decode(&result)
	useful.PanicOnUnrecoverableError(err)

	return result, nil
}
