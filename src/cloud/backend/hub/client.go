package hub

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/utils"
)

var client = utils.ComponentClient{
	RootUrl: "http://localhost:8082",
}

type UserAndApp struct {
	User string `json:"user"`
	App  string `json:"app"`
}

type TagInfo struct {
	User string `json:"user"`
	App  string `json:"app"`
	Tag  string `json:"tag"`
}

type HubClient interface {
	SearchApps(searchTerm string) (*[]UserAndApp, error)
	GetTags(userAndApp UserAndApp) (*[]string, error)
	DownloadTag(tagInfo TagInfo) (*[]byte, error)
}

type hubClientReal struct{}

func NewHubClientReal() HubClient {
	return &hubClientReal{}
}

func (h hubClientReal) SearchApps(searchTerm string) (*[]UserAndApp, error) {
	responseBody, err := client.DoRequest("/apps/search", utils.SingleString{searchTerm}, "")
	if err != nil {
		return nil, err
	}
	userAndAppList, err := unpackResponse[[]UserAndApp](responseBody)
	if err != nil {
		return nil, err
	}
	return userAndAppList, nil
}

func (h hubClientReal) GetTags(userAndApp UserAndApp) (*[]string, error) {
	responseBody, err := client.DoRequest("/tags/get-tags", userAndApp, "")
	if err != nil {
		return nil, err
	}

	tagList, err := unpackResponse[[]string](responseBody)
	if err != nil {
		return nil, err
	}
	return tagList, nil
}

func (h hubClientReal) DownloadTag(tagInfo TagInfo) (*[]byte, error) {
	result, err := client.DoRequest("/tags/download", tagInfo, "")
	if err != nil {
		return nil, err
	}
	downloadedContent, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("Failed to assert result to []byte")
	}
	return &downloadedContent, nil
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

type hubClientMock struct{}

func (h hubClientMock) SearchApps(searchTerm string) (*[]UserAndApp, error) {
	return &[]UserAndApp{
		{"sampleuser", "nginxdefault"},
	}, nil
}

func (h hubClientMock) GetTags(userAndApp UserAndApp) (*[]string, error) {
	return &[]string{"0.0.1"}, nil
}

func (h hubClientMock) DownloadTag(tagInfo TagInfo) (*[]byte, error) {
	data := []byte("sample content")
	return &data, nil
}

func NewHubClientMock() HubClient {
	return &hubClientMock{}
}
