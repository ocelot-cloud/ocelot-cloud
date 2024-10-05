package security

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

var (
	sampleMaintainer = "maintainer"
	sampleApp        = "app"
	sampleTag        = "1.0"
	sampleBlob       = []byte("hello")
)

// TODO Trigger already existing and not existing errors.
func TestAppLifecycle(t *testing.T) {
	defer dbRepo.WipeDatabase()

	maintainersAndApps, err := repo.ListApps()
	assert.Nil(t, err)
	assert.Nil(t, maintainersAndApps)

	tags, err := repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.NotNil(t, err)
	assert.Nil(t, tags)

	blob, err := repo.LoadTagBlob(sampleMaintainer, sampleApp, sampleTag)
	assert.NotNil(t, err)
	assert.Nil(t, blob)

	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	maintainersAndApps, err = repo.ListApps()
	assert.Nil(t, err)
	assert.NotNil(t, maintainersAndApps)
	assert.Equal(t, 1, len(maintainersAndApps))
	assert.Equal(t, sampleMaintainer, maintainersAndApps[0].Maintainer)
	assert.Equal(t, sampleApp, maintainersAndApps[0].App)

	tags, err = repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, sampleTag, tags[0])

	blob, err = repo.LoadTagBlob(sampleMaintainer, sampleApp, sampleTag)
	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, sampleBlob, blob)

	// TODO Deleting the only tag left should also delete the app.
	// TODO Test creating a second app with different tag, as this should not cause collisions if handled correctly.
}

func TestDeleteApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, repo.DeleteApp(sampleMaintainer, sampleApp))

	apps, err := repo.ListApps()
	assert.Nil(t, err)
	assert.Nil(t, apps)

	tags, err := repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.NotNil(t, err)
	assert.Nil(t, tags)
}

func TestCreatingTwoTagsInApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	sampleTag2 := "2.0"
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag2, sampleBlob))

	app, err := repo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(app))
	assert.Equal(t, sampleMaintainer, app[0].Maintainer)
	assert.Equal(t, sampleApp, app[0].App)

	tags, err := repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tags))

	assert.True(t, contain(tags, sampleTag))
	assert.True(t, contain(tags, sampleTag2))
}

func contain(tags []string, expectedTag string) bool {
	for _, actualTag := range tags {
		if actualTag == expectedTag {
			return true
		}
	}
	return false
}

func TestDeleteTag(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, repo.DeleteTag(sampleMaintainer, sampleApp, sampleTag))

	maintainersAndApps, err := repo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(maintainersAndApps))

	tags, err := repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.Nil(t, tags)

	blob, err := repo.LoadTagBlob(sampleMaintainer, sampleApp, sampleTag)
	assert.NotNil(t, err)
	assert.Nil(t, blob)
}

// TODO check if expiration is working
// TODO can't set a cookie without user
// TODO all inconsistencies should be handled in this layer -> user does not exist, user already existing etc.
// TODO error: user already exists
// TODO SetCookie, DeleteCookie, IsCookieValid
// TODO the DB interface appears to grow quite large when all all use cases are implemented. Check if could be split up.
// TODO Test deletion cascading, e.g. deleting user should also delete his group memberships etc.
