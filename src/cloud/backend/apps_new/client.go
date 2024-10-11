package apps_new

import (
	"fmt"
	"github.com/ocelot-cloud/shared/utils"
)

// TODO When errors occur, add logs

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
	userAndAppList, err := utils.UnpackResponse[[]UserAndApp](responseBody)
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

	tagList, err := utils.UnpackResponse[[]string](responseBody)
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
	data, err := utils.ZipDirectoryToBytes("sampleuser_nginxdefault")
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func NewHubClientMock() HubClient {
	return &hubClientMock{}
}
