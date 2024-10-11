package apps_new

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/repo"
	"testing"
)

func TestDownloadTag(t *testing.T) {
	hubClient = NewHubClientMock().(HubClient)
	repo.InitializeDatabaseWithSource(":memory:")
	tagInfo := TagInfo{"sampleuser", "nginxdefault", "0.0.1"}

	err := DownloadTag(tagInfo)
	assert.Nil(t, err)
	apps, err := repo.AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
	MaintainerAndApp := apps[0]
	assert.Equal(t, tagInfo.User, MaintainerAndApp.Maintainer)
	assert.Equal(t, tagInfo.App, MaintainerAndApp.App)
	tags, err := repo.AppRepo.ListTagsOfApp(tagInfo.User, tagInfo.App)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, tagInfo.Tag, tags[0])

	blob, err := repo.AppRepo.LoadTagBlob(tagInfo.User, tagInfo.App, tagInfo.Tag)
	assert.Nil(t, err)
	assert.Equal(t, expectedSampleTagSizeInByte, len(blob))

	// TODO Create network before starting and delete container afterwards
	err = StartContainer(tagInfo)
	assert.Nil(t, err)
}

// TODO Can MaintainerAndApp be merged with UserAndApp?
