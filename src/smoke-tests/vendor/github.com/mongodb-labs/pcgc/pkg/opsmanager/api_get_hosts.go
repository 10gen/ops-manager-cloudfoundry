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

// HostResponse represents a host information
type HostResponse struct {
	ID                 string        `json:"id"`
	Aliases            []string      `json:"aliases,omitempty"`
	AlertsEnabled      bool          `json:"alertsEnabled"`
	AuthMechanismName  string        `json:"authMechanismName,omitempty"`
	ClusterID          string        `json:"clusterId,omitempty"`
	Created            string        `json:"created,omitempty"`
	Deactivated        bool          `json:"deactivated"`
	GroupID            string        `json:"groupId,omitempty"`
	HasStartupWarnings bool          `json:"hasStartupWarnings"`
	Hidden             bool          `json:"hidden"`
	HiddenSecondary    bool          `json:"hiddenSecondary"`
	HostEnabled        bool          `json:"hostEnabled"`
	Hostname           string        `json:"hostname,omitempty"`
	IPAddress          string        `json:"ipAddress,omitempty"`
	JournalingEnabled  bool          `json:"journalingEnabled"`
	LastDataSizeBytes  int           `json:"lastDataSizeBytes,omitempty"`
	LastIndexSizeBytes int           `json:"lastIndexSizeBytes,omitempty"`
	LastPing           string        `json:"lastPing,omitempty"`
	LastRestart        string        `json:"lastRestart,omitempty"`
	Links              []interface{} `json:"links,omitempty"`
	LogsEnabled        bool          `json:"logsEnabled"`
	LowULimit          bool          `json:"lowULimit"`
	MuninEnabled       bool          `json:"muninEnabled"`
	MuninPort          int           `json:"muninPort,omitempty"`
	Port               int           `json:"port,omitempty"`
	ProfilerEnabled    bool          `json:"profilerEnabled"`
	ReplicaSetName     string        `json:"replicaSetName,omitempty"`
	ReplicaStateName   string        `json:"replicaStateName,omitempty"`
	ShardName          string        `json:"shardName,omitempty"`
	SlaveDelaySec      int           `json:"slaveDelaySec,omitempty"`
	SSLEnabled         bool          `json:"sslEnabled"`
	TypeName           string        `json:"typeName,omitempty"`
	UptimeMSec         int           `json:"uptimeMsec,omitempty"`
	Version            string        `json:"version,omitempty"`
}

// HostsResponse an array of hosts
type HostsResponse struct {
	Results    []HostResponse `json:"results,omitempty"`
	Links      []interface{}  `json:"links,omitempty"`
	TotalCount int            `json:"totalCount,omitempty"`
}

func (client opsManagerClient) GetHosts(projectID string) (HostsResponse, error) {
	var result HostsResponse

	url := client.resolver.Of("/groups/%s/hosts", projectID)
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
