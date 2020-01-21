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

// BackupConfigs response structure
type BackupConfigs struct {
	Links         []Link         `json:"links"`
	BackupResults []BackupResult `json:"results"`
	TotalCount    int            `json:"totalCount"`
}

// BackupResult response structure
type BackupResult struct {
	AuthMechanismName  string   `json:"authMechanismName,omitempty"`
	ClusterID          string   `json:"clusterId,omitempty"`
	EncryptionEnabled  bool     `json:"encryptionEnabled,omitempty"`
	ExcludedNamespaces []string `json:"excludedNamespaces,omitempty"`
	IncludedNamespaces []string `json:"includednamespaces,omitempty"`
	GroupID            string   `json:"groupId,omitempty"`
	Links              []Link   `json:"links,omitempty"`
	SSLEnabled         bool     `json:"sslEnabled,omitempty"`
	StatusName         string   `json:"statusName,omitempty"`
	StorageEngineName  string   `json:"storageEngineName,omitempty"`
	Username           string   `json:"username,omitempty"`
}

// GetBackupConfigs
// https://docs.opsmanager.mongodb.com/master/reference/api/backup/get-all-backup-configs-for-group/
func (client opsManagerClient) GetBackupConfigs(projectID string) (BackupConfigs, error) {
	var result BackupConfigs

	url := client.resolver.Of("/groups/%s/backupConfigs", projectID)
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
