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
)

// AgentAttributes represent an agent properties
type AgentAttributes struct {
	LogPath                 string     `json:"logPath,omitempty"`
	LogPathWindows          string     `json:"logPathWindows,omitempty"`
	LogRotate               *LogRotate `json:"logRotate"`
	Username                string     `json:"username,omitempty"`
	Password                string     `json:"password,omitempty"`
	KerberosPrincipal       string     `json:"kerberosPrincipal,omitempty"`
	KerberosKeytab          string     `json:"kerberosKeytab,omitempty"`
	KerberowWindowsUsername string     `json:"kerberosWindowsUsername,omitempty"`
	KerberowWindowsPassword string     `json:"kerberosWindowsPassword,omitempty"`
	SSLPEMKeyFile           string     `json:"sslPEMKeyFile,omitempty"`
	SSLPEMKeyFileWindows    string     `json:"sslPEMKeyFileWindows,omitempty"`
	SSLPEMKeyPwd            string     `json:"sslPEMKeyPwd,omitempty"`
}

func (client opsManagerClient) UpdateMonitoringConfig(projectID string, config AgentAttributes) error {
	bodyBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	url := client.resolver.Of("/groups/%s/automationConfig/monitoringAgentConfig", projectID)
	resp := client.PutJSON(url, bytes.NewReader(bodyBytes))
	return resp.Err
}

func (client opsManagerClient) UpdateBackupConfig(projectID string, config AgentAttributes) error {
	bodyBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	url := client.resolver.Of("/groups/%s/automationConfig/backupAgentConfig", projectID)
	resp := client.PutJSON(url, bytes.NewReader(bodyBytes))
	return resp.Err
}
