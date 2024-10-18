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
	searchApps, err := hubClient.SearchApps("nginxdefault")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*searchApps))
	app := (*searchApps)[0]
	assert.Equal(t, "sampleuser", app.Maintainer)
	assert.Equal(t, "nginxdefault", app.Name)
	tags, err := hubClient.GetTags(app.Id)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*tags))
	assert.Equal(t, "0.0.1", (*tags)[0].Name)
	tag := (*tags)[0]

	err = DownloadTag(tag.Id)
	assert.Nil(t, err)
	apps, err := repo.AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
	MaintainerAndApp := apps[0]
	assert.Equal(t, "sampleuser", MaintainerAndApp.Maintainer)
	assert.Equal(t, "nginxdefault", MaintainerAndApp.Name)
	appId, err := repo.AppRepo.GetAppId("sampleuser", "nginxdefault")
	assert.Nil(t, err)
	repoTags, err := repo.AppRepo.ListTagsOfApp(appId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(repoTags))
	assert.Equal(t, "0.0.1", repoTags[0].Name)

	tagId, err := repo.AppRepo.GetTagId(appId, "0.0.1")
	blob, err := repo.AppRepo.LoadTagBlob(tagId)
	assert.Nil(t, err)
	assert.Equal(t, expectedSampleTagSizeInByte, len(blob))

	// TODO Duplication
	exec.Command("/bin/sh", "-c", "docker network ls | grep -q ocelot-net || docker network create ocelot-net").Run()
	err = StartContainer(appId)
	assert.Nil(t, err)

	err = exec.Command("/bin/sh", "-c", "docker ps | grep -q nginx").Run() // TODO abstract the "nginx"
	assert.Nil(t, err)

	err = StopContainer(appId)
	assert.Nil(t, err)

	err = exec.Command("/bin/sh", "-c", "docker ps -a | grep -q nginx").Run()
	assert.NotNil(t, err)
}

// TODO Make an integration test similar to the test above, but which does a http request to the nginx container.
// TODO Can MaintainerAndApp be merged with UserAndApp?
// TODO New network approach should be added. E.g. Starting stack "gitea" should create a network "gitea-net". This is the only network "gitea" is member of. Ocelot joins the network.
//   -> "docker network connect ocelot-net ocelot-cloud"
//   TODO That network connection does not survive a reboot or a container restart. It must be explicitly reconnected. -> Could be done at start "reconnect to all running containers networks"
// TODO remove the "ocelot-net" from the docker-compose.yml files.
