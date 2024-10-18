package apps_new

import (
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

type HubApp struct {
}

type HubClient interface {
	SearchApps(searchTerm string) (*[]tools.App, error)
	GetTags(appId int) (*[]tools.Tag, error)
	DownloadTag(tagId int) (*tools.FullTagInfo, error)
}

type hubClientReal struct{}

func NewHubClientReal() HubClient {
	return &hubClientReal{}
}

func (h hubClientReal) SearchApps(searchTerm string) (*[]tools.App, error) {
	responseBody, err := client.DoRequest("/apps/search", utils.SingleString{searchTerm}, "")
	if err != nil {
		return nil, err
	}
	userAndAppList, err := utils.UnpackResponse[[]tools.App](responseBody)
	if err != nil {
		return nil, err
	}
	return userAndAppList, nil
}

// TODO Duplication tools.SingleInt and utils.SingleInteger

func (h hubClientReal) GetTags(appId int) (*[]tools.Tag, error) {
	responseBody, err := client.DoRequest("/tags/get-tags", tools.SingleInt{appId}, "")
	if err != nil {
		return nil, err
	}

	tagList, err := utils.UnpackResponse[[]tools.Tag](responseBody)
	if err != nil {
		return nil, err
	}
	return tagList, nil
}

func (h hubClientReal) DownloadTag(tagId int) (*tools.FullTagInfo, error) {
	result, err := client.DoRequest("/tags/download", tools.SingleInt{tagId}, "")
	if err != nil {
		return nil, err
	}
	fullTagInfo, err := utils.UnpackResponse[tools.FullTagInfo](result)
	if err != nil {
		return nil, err
	}
	return fullTagInfo, nil
}

type hubClientMock struct{}

func (h hubClientMock) SearchApps(searchTerm string) (*[]tools.App, error) {
	return &[]tools.App{
		{"sampleuser", "nginxdefault", -1},
	}, nil
}

func (h hubClientMock) GetTags(appId int) (*[]tools.Tag, error) {
	return &[]tools.Tag{{"0.0.1", -1}}, nil
}

func (h hubClientMock) DownloadTag(tagId int) (*tools.FullTagInfo, error) {
	data, err := utils.ZipDirectoryToBytes(getSampleAppFolder())
	if err != nil {
		return nil, err
	}
	return &tools.FullTagInfo{
		Maintainer: "sampleuser",
		AppName:    "nginxdefault",
		TagName:    "0.0.1",
		Content:    data,
		Id:         -1,
	}, nil
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
