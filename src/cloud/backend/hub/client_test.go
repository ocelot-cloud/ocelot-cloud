package hub

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"testing"
)

type UserAndApp struct {
	User string `json:"user"`
	App  string `json:"app"`
}

func TestHubClient(t *testing.T) {
	client := utils.ComponentClient{
		RootUrl: "http://localhost:8082",
	}

	responseBody, err := client.DoRequest("/apps/search", utils.SingleString{"sample"}, "")
	assert.Nil(t, err)

	apps, err := unpackResponse[[]UserAndApp](responseBody)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(*apps))
	app := (*apps)[0]
	assert.Equal(t, "sampleuser", app.User)
	assert.Equal(t, "sampleapp", app.App)
}

func unpackResponse[T any](object interface{}) (*T, error) {
	respBody, ok := object.([]byte)
	if !ok {
		return nil, fmt.Errorf("Failed to assert result to []byte")
	}

	var result T
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v", err)
	}
	return &result, nil
}
