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

// Package opsmanager hosts a HTTP client which abstracts communication with Ops Manager instances.
//
// To create a new client, you have to call the following code:
//
//		resolver := httpclient.NewURLResolverWithPrefix("http://OPS-MANAGER-INSTANCE", opsmanager.PublicApiPrefix)
// 		client := opsmanager.NewClient(WithResolver(resolver))
//
// The client can then be used to issue requests such as:
//
//		user := User{Username: ..., Password: ..., ...}
// 		globalOwner := client.CreateFirstUser(user, WhitelistAllowAll)
//
// See the Client interface below for a list of all the support operations.
// If however, you need one that is not currently supported, the _opsManagerApi_ struct extends
// _httpclient.BasicHTTPOperation_, allowing you to issue raw HTTP requests to the specified Ops Manager instance.
//
// 		url := resolver.Of("/path/to/a/resource/%s", id)
//		resp:= client.GetJSON(url)
//		useful.PanicOnUnrecoverableError(resp.Err)
//		defer useful.LogError(resp.Response.Body.Close)
//		var data SomeType
//		decoder := json.NewDecoder(resp.Response.Body)
//		err := decoder.Decode(&result)
//		useful.PanicOnUnrecoverableError(err)
//		// do something with _data_
//
// You can create an authenticated client as follows:
//		resolver := httpclient.NewURLResolverWithPrefix("http://OPS-MANAGER-INSTANCE", PublicAPIPrefix)
//		client:= opsmanager.NewClientWithDigestAuth(resolver, "username", "password")
//
// The code above is a simplification of the following (which can be used for a more configurable set-up):
//		withResolver := opsmanager.WithResolver(httpclient.NewURLResolverWithPrefix("http://OPS-MANAGER-INSTANCE", PublicAPIPrefix))
//		withDigestAuth := httpclient.WithDigestAuthentication(publicKey, privateKey)
//		withHTTPClient := opsmanager.WithHTTPClient(httpclient.NewClient(withDigestAuth))
//		client := opsmanager.NewClient(withResolver, withHTTPClient)
//
// The following can be used for authentication:
//		- Ops Manager user credentials: (username, password) - only works with some APIs and should not be used
//		- Programmatic API keys: (publicKey, privateKey) - preferred credentials pair
//		- Ops Manager user and a Personal API Key: (username, personalAPIKey) - deprecated
// You can read more about this topic here: https://docs.opsmanager.mongodb.com/master/tutorial/configure-public-api-access/#configure-public-api-access
//
package opsmanager

import (
	"errors"
	"io"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
)

type opsManagerClient struct {
	httpclient.BasicClient

	resolver httpclient.URLResolver
}

// Client defines the API actions implemented in this client
type Client interface {
	httpclient.BasicClient

	// https://docs.opsmanager.mongodb.com/master/reference/api/user-create-first/
	CreateFirstUser(user User, whitelistIP string) (CreateFirstUserResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/groups/get-all-groups-for-current-user/
	GetAllProjects() (ProjectsResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/groups/create-one-group/
	CreateOneProject(name string, orgID string) (ProjectResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/#get-the-automation-configuration
	GetAutomationConfig(projectID string) (AutomationConfig, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/#update-the-automation-configuration
	UpdateAutomationConfig(projectID string, config AutomationConfig) (AutomationConfig, error)
	// GET /agents/api/automation/conf/v1/{projectID}
	GetRawAutomationConfig(projectID string) (RawAutomationConfig, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/automation-status/
	GetAutomationStatus(projectID string) (AutomationStatusResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/agents-get-by-type/
	GetAgentsByType(projectID string, agentType string) (GetAgentsByTypeResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/#update-the-automation-configuration
	UpdateDeployments(projectID string, body io.Reader) (map[string]interface{}, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/agentapikeys/create-one-agent-api-key/
	CreateAgentAPIKEY(projectID string, desc string) (CreateAgentAPIKEYResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/groups/get-one-group-by-id/
	GetProjectByID(projectID string) (ProjectResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/groups/get-one-group-by-name/
	GetProjectByName(name string) (ProjectResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/groups/delete-one-group/
	DeleteProject(projectID string) error
	// https://docs.opsmanager.mongodb.com/master/reference/api/groups/add-or-remove-tags-from-one-group/
	SetProjectTags(projectID string, tags []string) (ProjectResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/hosts/get-all-hosts-in-group/
	GetHosts(projectID string) (HostsResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/index.html#update-the-monitoring-or-backup
	UpdateMonitoringConfig(projectID string, config AgentAttributes) error
	// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/index.html#update-the-monitoring-or-backup
	UpdateBackupConfig(projectID string, config AgentAttributes) error
}

// NewClient builds a new API client for connecting to Ops Manager
func NewClient(configs ...func(*opsManagerClient)) Client {
	// initialize a bare client
	client := &opsManagerClient{}

	// apply all configurations
	for _, configure := range configs {
		configure(client)
	}

	// validations
	if client.resolver == nil {
		useful.PanicOnUnrecoverableError(errors.New("the client requires a URLResolver with the appropriate Ops Manager URL configured"))
	}
	if client.BasicClient == nil {
		useful.PanicOnUnrecoverableError(errors.New("the client requires an underlying basic HTTP client to be configured"))
	}

	return client
}

// WithResolver configures an Ops Manager client which relies on the specified resolver
func WithResolver(resolver httpclient.URLResolver) func(*opsManagerClient) {
	return func(client *opsManagerClient) {
		client.resolver = resolver
	}
}

// WithHTTPClient configures an Ops Manager which delegates basic HTTP operations to the specified client
func WithHTTPClient(basicClient httpclient.BasicClient) func(*opsManagerClient) {
	return func(client *opsManagerClient) {
		client.BasicClient = basicClient
	}
}

// NewDefaultClient builds a new, unauthenticated, API client with default configurations
func NewDefaultClient(resolver httpclient.URLResolver) Client {
	return NewClient(WithHTTPClient(httpclient.NewClient()), WithResolver(resolver))
}

// NewClientWithDigestAuth builds a new API client with default configurations, which uses digest authentication
func NewClientWithDigestAuth(resolver httpclient.URLResolver, username string, password string) Client {
	return NewClient(WithHTTPClient(httpclient.NewClient(httpclient.WithDigestAuthentication(username, password))), WithResolver(resolver))
}
