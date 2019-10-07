package opsmanager

import (
	"encoding/json"
	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
)

// RawAutomationConfig represents a raw automation config
// NOTE: this struct is mutable
type RawAutomationConfig struct {
	IDCounter                        int                    `json:"idCounter,omitempty"`
	State                            string                 `json:"state,omitempty"`
	Version                          int                    `json:"version,omitempty"`
	GroupID                          string                 `json:"groupId,omitempty"`
	UserID                           string                 `json:"userId,omitempty"`
	PublishTimestamp                 int                    `json:"publishTimestamp,omitempty"`
	PublishedVersion                 string                 `json:"publishedVersion,omitempty"`
	SaveTimestamp                    int                    `json:"saveTimestamp,omitempty"`
	LatestAutomationAgentVersionName string                 `json:"latestAutomationAgentVersionName,omitempty"`
	LatestMonitoringAgentVersionName string                 `json:"latestMonitoringAgentVersionName,omitempty"`
	LatestBackupAgentVersionName     string                 `json:"latestBackupAgentVersionName,omitempty"`
	LatestBiConnectorVersionName     string                 `json:"latestBiConnectorVersionName,omitempty"`
	Cluster                          *AutomationConfig      `json:"cluster,omitempty"`
	VersionConfig                    map[string]interface{} `json:"versionConfig,omitempty"`
	LogRotate                        map[string]interface{} `json:"logRotate,omitempty"`
	MonitoringAgentTemplate          map[string]interface{} `json:"monitoringAgentTemplate,omitempty"`
	BackupAgentTemplate              map[string]interface{} `json:"backupAgentTemplate,omitempty"`
	CPSModuleTemplate                map[string]interface{} `json:"cpsModuleTemplate,omitempty"`
	DeploymentJobStatuses            []interface{}
}

// GetRawAutomationConfig returns the RAW automation config, just like the Automation Agent sees it
// /agents/api/automation/conf/v1/{projectID}
func (client opsManagerClient) GetRawAutomationConfig(projectID string) (RawAutomationConfig, error) {
	var result RawAutomationConfig

	url := client.resolver.OfUnprefixed("/agents/api/automation/conf/v1/%s", projectID)
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
