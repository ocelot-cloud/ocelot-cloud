package apps_new

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

var expectedSampleTagSizeInByte = 903

func TestHubClientReal(t *testing.T) {
	hubClient := NewHubClientReal().(HubClient)
	conductApiChecks(t, hubClient)
}

func conductApiChecks(t *testing.T, hubClient HubClient) {
	apps, err := hubClient.SearchApps("sample")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*apps))
	app := (*apps)[0]
	assert.Equal(t, "sampleuser", app.Maintainer)
	assert.Equal(t, "nginxdefault", app.Name)

	tagList, err := hubClient.GetTags(app.Id)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*tagList))
	tag := (*tagList)[0]
	assert.Equal(t, "0.0.1", tag.Name)

	tagContent, err := hubClient.DownloadTag(tag.Id)
	assert.Nil(t, err)
	assert.Equal(t, expectedSampleTagSizeInByte, len(*tagContent))
}

func TestHubClientMock(t *testing.T) {
	hubClient := NewHubClientMock().(HubClient)
	conductApiChecks(t, hubClient)
}
