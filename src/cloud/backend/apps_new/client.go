package apps_new

import (
	"fmt"
	"github.com/ocelot-cloud/shared/utils"
	"ocelot/backend/tools"
	"os"
	"path/filepath"
	"strings"
)

// TODO When errors occur, add logs

var client = utils.ComponentClient{
	RootUrl: "http://localhost:8082",
}

type HubClient interface {
	SearchApps(searchTerm string) (*[]tools.UserAndApp, error)
	GetTags(userAndApp tools.UserAndApp) (*[]string, error)
	DownloadTag(tagInfo tools.TagInfo) (*[]byte, error)
}

type hubClientReal struct{}

func NewHubClientReal() HubClient {
	return &hubClientReal{}
}

func (h hubClientReal) SearchApps(searchTerm string) (*[]tools.UserAndApp, error) {
	responseBody, err := client.DoRequest("/apps/search", utils.SingleString{searchTerm}, "")
	if err != nil {
		return nil, err
	}
	userAndAppList, err := utils.UnpackResponse[[]tools.UserAndApp](responseBody)
	if err != nil {
		return nil, err
	}
	return userAndAppList, nil
}

func (h hubClientReal) GetTags(userAndApp tools.UserAndApp) (*[]string, error) {
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

func (h hubClientReal) DownloadTag(tagInfo tools.TagInfo) (*[]byte, error) {
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

func (h hubClientMock) SearchApps(searchTerm string) (*[]tools.UserAndApp, error) {
	return &[]tools.UserAndApp{
		{"sampleuser", "nginxdefault"},
	}, nil
}

func (h hubClientMock) GetTags(userAndApp tools.UserAndApp) (*[]string, error) {
	return &[]string{"0.0.1"}, nil
}

func (h hubClientMock) DownloadTag(tagInfo tools.TagInfo) (*[]byte, error) {
	data, err := utils.ZipDirectoryToBytes(getSampleAppFolder())
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func NewHubClientMock() HubClient {
	return &hubClientMock{}
}

func getSampleAppFolder() string {
	currentDir, err := os.Getwd()
	if err != nil {
		Logger.Fatal("Failed to get current dir: %v", err)
	}
	parentDir := filepath.Base(currentDir)
	Logger.Debug("Current dir is backend")
	if strings.EqualFold(parentDir, "backend") {
		return "apps_new/sampleuser_nginxdefault"
	}
	return "sampleuser_nginxdefault"
}
