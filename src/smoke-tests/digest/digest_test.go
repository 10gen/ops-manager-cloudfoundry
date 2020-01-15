package digest_test

import (
	"fmt"
	"net/http"
	"testing"

	"mongodb-service-adapter/adapter"
	"mongodb-service-adapter/digest"
)

func TestApplyDigestAuth(t *testing.T) {
	config, err := adapter.LoadConfig("../../mongodb-service-adapter/testdata/manifest.json")

	if err != nil {
		fmt.Print("Error opening manifest file ")
	}

	method := "GET"
	path := fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig", config.GroupID)
	uri := fmt.Sprintf("%s%s", config.URL, path)
	req, err := http.NewRequest(method, uri, nil)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	err = digest.ApplyDigestAuth(config.Username, config.APIKey, uri, req)

	if err != nil {
		t.Fatal(err)
	}
}
