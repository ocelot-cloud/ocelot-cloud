package hub

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestHubClient(t *testing.T) {
	hubClient := NewHubClient()
	userAndAppList, err := hubClient.SearchApps("sample")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*userAndAppList))
	userAndApp := (*userAndAppList)[0]
	assert.Equal(t, "sampleuser", userAndApp.User)
	assert.Equal(t, "nginxdefault", userAndApp.App)

	responseBody, err := client.DoRequest("/tags/get-tags", userAndApp, "")
	assert.Nil(t, err)

	tagList, err := unpackResponse[[]string](responseBody)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(*tagList))
	tag := (*tagList)[0]
	assert.Equal(t, "0.0.1", tag)

	tagInfo := TagInfo{userAndApp.User, userAndApp.App, tag}
	result, err := client.DoRequest("/tags/download", tagInfo, "")
	assert.Nil(t, err)
	downloadedContent, ok := result.([]byte)
	assert.True(t, ok)
	assert.Equal(t, 1260, len(downloadedContent))
}
