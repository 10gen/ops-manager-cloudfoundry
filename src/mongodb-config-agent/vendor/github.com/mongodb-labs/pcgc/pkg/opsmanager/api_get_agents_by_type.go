package opsmanager

import (
	"encoding/json"
	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
)

// AgentResult represents an agent returned by the GetAgentsByType() call
type AgentResult struct {
	ConfCount int    `json:"confCount,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	LastConf  string `json:"lastConf,omitempty"`
	StateName string `json:"stateName,omitempty"`
	TypeName  string `json:"typeName,omitempty"`
}

// GetAgentsByTypeResponse API response for the GetAgentsByType() call
type GetAgentsByTypeResponse struct {
	Links      []Link        `json:"links,omitempty"`
	Results    []AgentResult `json:"results,omitempty"`
	TotalCount int           `json:"totalCount,omitempty"`
}

// GetAllHostsInProject
// https://docs.opsmanager.mongodb.com/current/reference/api/hosts/get-all-hosts-in-group/
func (client opsManagerClient) GetAgentsByType(projectID string, agentType string) (GetAgentsByTypeResponse, error) {
	var result GetAgentsByTypeResponse

	url := client.resolver.Of("/groups/%s/agents/%s", projectID, agentType)
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
