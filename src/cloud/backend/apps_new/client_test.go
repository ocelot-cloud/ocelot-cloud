package apps_new

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestHubClientReal(t *testing.T) {
	hubClient := NewHubClientReal().(HubClient)
	conductApiChecks(t, hubClient)
}

func conductApiChecks(t *testing.T, hubClient HubClient) {
	userAndAppList, err := hubClient.SearchApps("sample")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*userAndAppList))
	userAndApp := (*userAndAppList)[0]
	assert.Equal(t, "sampleuser", userAndApp.User)
	assert.Equal(t, "nginxdefault", userAndApp.App)

	tagList, err := hubClient.GetTags(userAndApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*tagList))
	tag := (*tagList)[0]
	assert.Equal(t, "0.0.1", tag)

	tagInfo := TagInfo{userAndApp.User, userAndApp.App, tag}
	tagContent, err := hubClient.DownloadTag(tagInfo)
	assert.Nil(t, err)
	assert.Equal(t, 1260, len(*tagContent))
}

func TestHubClientMock(t *testing.T) {
	hubClient := NewHubClientMock().(HubClient)
	conductApiChecks(t, hubClient)
}
