package opsmanager

import (
	"encoding/json"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
)

// HostCounts for the project
type HostCounts struct {
	Arbiter   int `json:"arbiter"`
	Config    int `json:"config"`
	Primary   int `json:"primary"`
	Secondary int `json:"secondary"`
	Mongos    int `json:"mongos"`
	Master    int `json:"master"`
	Slave     int `json:"slave"`
}

// LDAPGroupMapping holds a single mapping of role to LDAP groups
type LDAPGroupMapping struct {
	RoleName   string   `json:"roleName,omitempty"`
	LdapGroups []string `json:"ldapGroups,omitempty"`
}

// ProjectResponse represents the structure of a project
type ProjectResponse struct {
	ID                string             `json:"id"`
	OrgID             string             `json:"orgId"`
	Name              string             `json:"name"`
	LastActiveAgent   string             `json:"lastActiveAgent,omitempty"`
	AgentAPIKey       string             `json:"agentApiKey,omitempty"`
	ActiveAgentCount  int                `json:"activeAgentCount"`
	HostCounts        HostCounts         `json:"hostCounts,omitempty"`
	PublicAPIEnabled  bool               `json:"publicApiEnabled"`
	ReplicaSetCount   int                `json:"replicaSetCount"`
	ShardCount        int                `json:"shardCount"`
	Tags              []string           `json:"tags,omitempty"`
	Links             []Link             `json:"links,omitempty"`
	LDAPGroupMappings []LDAPGroupMapping `json:"ldapGroupMappings,omitempty"`
}

// ProjectsResponse represents a array of project
type ProjectsResponse struct {
	Links      []Link            `json:"links"`
	Results    []ProjectResponse `json:"results"`
	TotalCount int               `json:"totalCount"`
}

// Result is part of TeamsAssigned structure
type Result struct {
	Links     []Link   `json:"links"`
	RoleNames []string `json:"roleNames"`
	TeamID    string   `json:"teamId"`
}

// GetAllProjects registers the first ever Ops Manager user (global owner)
// https://docs.opsmanager.mongodb.com/master/reference/api/groups/get-all-groups-for-current-user/
func (client opsManagerClient) GetAllProjects() (ProjectsResponse, error) {
	var result ProjectsResponse

	url := client.resolver.Of("/groups")
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
