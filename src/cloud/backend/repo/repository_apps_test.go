package repo

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/tools"
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

	tagInfo := tools.TagInfo{sampleMaintainer, sampleApp, sampleTag}
	assert.False(t, AppRepo.DoesTagExist(tagInfo))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	apps, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.NotNil(t, apps)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, sampleMaintainer, (apps)[0].Maintainer)
	assert.Equal(t, sampleApp, (apps)[0].Name)
	assert.True(t, AppRepo.DoesTagExist(tagInfo))

	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	tags, err := AppRepo.ListTagsOfApp(appId)
	assert.Nil(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, sampleTag, tags[0].Name)

	tagId, err := AppRepo.GetTagId(appId, sampleTag)
	assert.Nil(t, err)
	blob, err := AppRepo.LoadTagBlob(tagId)
	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, sampleBlob, blob)

	// TODO Deleting the only tag left should also delete the app.
	//   assert.False(t, AppRepo.DoesTagExist(tagInfo))
	// TODO Test creating a second app with different tag, as this should not cause collisions if handled correctly.
}

func TestDeleteApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, AppRepo.DeleteApp(sampleMaintainer, sampleApp))

	apps, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Nil(t, apps)

	_, err = AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.NotNil(t, err)
}

func TestCreatingTwoTagsInApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	sampleTag2 := "2.0"
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag2, sampleBlob))

	app, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(app))
	assert.Equal(t, sampleMaintainer, app[0].Maintainer)
	assert.Equal(t, sampleApp, app[0].Name)

	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)

	tags, err := AppRepo.ListTagsOfApp(appId)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tags))

	assert.True(t, contain(tags, sampleTag))
	assert.True(t, contain(tags, sampleTag2))
}

func contain(tags []Tag, expectedTagName string) bool {
	for _, actualTag := range tags {
		if actualTag.Name == expectedTagName {
			return true
		}
	}
	return false
}

func TestDeleteTag(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, AppRepo.DeleteTag(sampleMaintainer, sampleApp, sampleTag))

	maintainersAndApps, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(maintainersAndApps))

	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	tags, err := AppRepo.ListTagsOfApp(appId)
	assert.Nil(t, err)
	assert.Nil(t, tags)

	_, err = AppRepo.GetTagId(appId, sampleTag)
	assert.NotNil(t, err)
}

// TODO check if expiration is working
// TODO can't set a cookie without user
// TODO all inconsistencies should be handled in this layer -> user does not exist, user already existing etc.
// TODO error: user already exists
// TODO SetCookie, DeleteCookie, IsCookieValid
// TODO the DB interface appears to grow quite large when all all use cases are implemented. Check if could be split up.
// TODO Test deletion cascading, e.g. deleting user should also delete his group memberships etc.
