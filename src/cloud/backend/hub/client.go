package hub

import "github.com/ocelot-cloud/shared/utils"

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
	SearchApps(searchTerm string) ([]UserAndApp, error)
	GetTags(userAndApp UserAndApp) ([]string, error)
	DownloadTag(tagInfo TagInfo) ([]byte, error)
}

type hubClientReal struct{}

func NewHubClient() HubClient {
	return &hubClientReal{}
}

func (h hubClientReal) SearchApps(searchTerm string) ([]UserAndApp, error) {
	//TODO implement me
	panic("implement me")
}

func (h hubClientReal) GetTags(userAndApp UserAndApp) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (h hubClientReal) DownloadTag(tagInfo TagInfo) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
