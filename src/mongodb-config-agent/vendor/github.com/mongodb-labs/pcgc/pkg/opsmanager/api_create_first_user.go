package opsmanager

import (
	"bytes"
	"encoding/json"
	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
)

// User request object which identifies a user
type User struct {
	Username     string `json:"username"`
	Password     string `json:"password,omitempty"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	EmailAddress string `json:"emailAddress,omitempty"`
}

// UserRole denotes a single user role
type UserRole struct {
	RoleName string `json:"roleName"`
	GroupID  string `json:"groupId,omitempty"`
	OrgID    string `json:"orgId,omitempty"`
}

// UserResponse wrapper for a user response, augmented with a few extra fields
type UserResponse struct {
	User

	ID    string     `json:"id"`
	Links []Link     `json:"links,omitempty"`
	Roles []UserRole `json:"roles,omitempty"`
}

// CreateFirstUserResponse API response for the CreateFirstUser() call
type CreateFirstUserResponse struct {
	APIKey string       `json:"apiKey"`
	User   UserResponse `json:"user"`
}

// CreateFirstUser registers the first ever Ops Manager user (global owner)
// pass a whitelist of 0.0.0.1/0 if you want to whitelist all IPv4 addresses
// https://docs.opsmanager.mongodb.com/master/reference/api/user-create-first/
func (client opsManagerClient) CreateFirstUser(user User, whitelistIP string) (CreateFirstUserResponse, error) {
	var result CreateFirstUserResponse

	bodyBytes, err := json.Marshal(user)
	if err != nil {
		return result, err
	}

	// if a whitelist was not specified, do not pass the parameter
	var url string
	if whitelistIP == "" {
		url = client.resolver.Of("/unauth/users")
	} else {
		url = client.resolver.Of("/unauth/users?whitelist=%s", whitelistIP)
	}

	resp := client.PostJSON(url, bytes.NewReader(bodyBytes))
	if resp.IsError() {
		return result, resp.Err
	}
	defer httpclient.CloseResponseBodyIfNotNil(resp)

	decoder := json.NewDecoder(resp.Response.Body)
	err = decoder.Decode(&result)
	useful.PanicOnUnrecoverableError(err)

	return result, nil
}
