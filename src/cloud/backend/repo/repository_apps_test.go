package repo

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

	appId, tagId := createAppAndTag(t)
	apps, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.NotNil(t, apps)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, sampleMaintainer, (apps)[0].Maintainer)
	assert.Equal(t, sampleApp, (apps)[0].Name)

	tags, err := AppRepo.ListTagsOfApp(appId)
	assert.Nil(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, sampleTag, tags[0].Name)

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
	appId, _ := createAppAndTag(t)
	assert.Nil(t, AppRepo.DeleteApp(appId))

	apps, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Nil(t, apps)

	_, err = AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.NotNil(t, err)
}

func TestCreatingTwoTagsInApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	sampleTag2 := "2.0"
	appId, _ := createAppAndTag(t)
	assert.Nil(t, AppRepo.CreateTag(appId, sampleTag2, sampleBlob))

	app, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(app))
	assert.Equal(t, sampleMaintainer, app[0].Maintainer)
	assert.Equal(t, sampleApp, app[0].Name)

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
	appId, tagId := createAppAndTag(t)
	assert.Nil(t, AppRepo.DeleteTag(tagId))

	maintainersAndApps, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(maintainersAndApps))

	tags, err := AppRepo.ListTagsOfApp(appId)
	assert.Nil(t, err)
	assert.Nil(t, tags)

	_, err = AppRepo.GetTagId(appId, sampleTag)
	assert.NotNil(t, err)
}

func TestActiveTag(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, AppRepo.CreateApp(sampleMaintainer, sampleApp))
	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)

	app, err := AppRepo.GetApp(appId)
	assert.Nil(t, err)
	assert.Equal(t, "", app.ActiveTagName)
	assert.Equal(t, -1, app.ActiveTagId)
	apps, err := AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, "", apps[0].ActiveTagName)
	assert.Equal(t, -1, apps[0].ActiveTagId)

	assert.Nil(t, AppRepo.CreateTag(appId, sampleTag, sampleBlob))

	app, err = AppRepo.GetApp(appId)
	assert.Nil(t, err)
	assert.Equal(t, sampleTag, app.ActiveTagName)
	assert.True(t, app.ActiveTagId >= 0)
	/* TODO
	apps, err = AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, sampleTag, apps[0].ActiveTagName)
	assert.True(t, apps[0].ActiveTagId >= 0)
	*/

	tagId, err := AppRepo.GetTagId(appId, sampleTag)
	assert.Nil(t, err)
	assert.Nil(t, AppRepo.DeleteTag(tagId))

	app, err = AppRepo.GetApp(appId)
	assert.Nil(t, err)
	assert.Equal(t, "", app.ActiveTagName)
	assert.Equal(t, -1, app.ActiveTagId)
	apps, err = AppRepo.ListApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, "", apps[0].ActiveTagName)
	assert.Equal(t, -1, apps[0].ActiveTagId)
}

// TODO When a request proxy does not work, the user should be redirected to the home page, in order to use the button to visit the app.
// TODO check if expiration of cookies and secret is working?
// TODO can't set a cookie without user
// TODO all inconsistencies should be handled in this layer -> user does not exist, user already existing etc.
// TODO error: user already exists
// TODO SetCookie, DeleteCookie, IsCookieValid
// TODO the DB interface appears to grow quite large when all all use cases are implemented. Check if could be split up.
// TODO Test deletion cascading, e.g. deleting user should also delete his group memberships etc.
