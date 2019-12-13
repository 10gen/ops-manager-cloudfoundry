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

// Build a MongoDB build
type Build struct {
	Architecture       string   `json:"architecture"`
	Bits               int      `json:"bits"`
	Flavor             string   `json:"flavor,omitempty"`
	GitVersion         string   `json:"gitVersion,omitempty"`
	MaxOsVersion       string   `json:"maxOsVersion,omitempty"`
	MinOsVersion       string   `json:"minOsVersion,omitempty"`
	Platform           string   `json:"platform,omitempty"`
	URL                string   `json:"url,omitempty"`
	Modules            []string `json:"modules,omitempty"`
	Win2008plus        bool     `json:"win2008plus,omitempty"`
	WinVCRedistDll     string   `json:"winVCRedistDll,omitempty"`
	WinVCRedistOptions []string `json:"winVCRedistOptions,omitempty"`
	WinVCRedistURL     string   `json:"winVCRedistUrl,omitempty"`
	WinVCRedistVersion string   `json:"winVCRedistVersion,omitempty"`
}

// MongoDBVersion ways to install MongoDB
type MongoDBVersion struct {
	Name   string  `json:"name,omitempty"`
	Builds []Build `json:"builds,omitempty"`
}

// GetAutomationConfig
// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/#get-the-automation-configuration
func (client opsManagerClient) GetAutomationConfig(projectID string) (AutomationConfig, error) {
	var result AutomationConfig

	url := client.resolver.Of("/groups/%s/automationConfig", projectID)
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
