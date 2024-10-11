package apps_new

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/repo"
	"os/exec"
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

	// TODO Duplication
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	err = StartContainer(tagInfo)
	assert.Nil(t, err)

	err = exec.Command("/bin/sh", "-c", "docker ps | grep -q nginx-default").Run()
	assert.Nil(t, err)

	err = StopContainer(tagInfo)
	assert.Nil(t, err)

	err = exec.Command("/bin/sh", "-c", "docker ps -a | grep -q nginx-default").Run()
	assert.NotNil(t, err)
}

// TODO Make an integration test similar to the test above, but which does a http request to the nginx container.
// TODO Can MaintainerAndApp be merged with UserAndApp?
// TODO New network approach should be added. E.g. Starting stack "gitea" should create a network "gitea-net". This is the only network "gitea" is member of. Ocelot joins the network.
// TODO remove the "ocelot-net" from the docker-compose.yml files.
