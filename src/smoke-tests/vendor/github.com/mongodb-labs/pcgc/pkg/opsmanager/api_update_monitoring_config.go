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
	SSLPEMKeyfile           string     `json:"sslPEMKeyfile,omitempty"`
	SSLPEMKeyfileWindows    string     `json:"sslPEMKeyfileWindows,omitempty"`
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
