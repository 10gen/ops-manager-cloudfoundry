package opsmanager

import (
	"bytes"
	"encoding/json"

	"github.com/mongodb-labs/pcgc/pkg/useful"
)

func (client opsManagerClient) SetProjectTags(projectID string, tags []string) (ProjectResponse, error) {
	var result ProjectResponse

	url := client.resolver.Of("/groups/%s", projectID)
	request := map[string][]string{
		"tags": tags,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return result, err
	}

	resp := client.PatchJSON(url, bytes.NewReader(body))
	decoder := json.NewDecoder(resp.Response.Body)
	err = decoder.Decode(&result)
	useful.PanicOnUnrecoverableError(err)

	return result, nil
}
