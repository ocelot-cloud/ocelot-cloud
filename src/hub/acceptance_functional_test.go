//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"strconv"
	"testing"
	"time"
)

func TestTagDownload(t *testing.T) {
	hub := getHub()

	_, err := hub.downloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "user does not exist"), err.Error())

	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())

	_, err = hub.downloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.createApp())
	_, err = hub.downloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "tag does not exist"), err.Error())

	assert.Nil(t, hub.uploadTag())
	foundTags, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, sampleTag, foundTags[0].Name)

	downloadedContent, err := hub.downloadTag()
	assert.Nil(t, err)
	assert.Equal(t, sampleTagFileContent, downloadedContent)
}

func TestCookie(t *testing.T) {
	hub := getHubAndLogin(t)

	assert.Equal(t, cookieName, hub.Parent.Cookie.Name)
	assert.True(t, utils.GetTimeIn30Days().Add(1*time.Second).After(hub.Parent.Cookie.Expires))
	assert.True(t, utils.GetTimeIn30Days().Add(-1*time.Second).Before(hub.Parent.Cookie.Expires))
	assert.Equal(t, 64, len(hub.Parent.Cookie.Value))

	cookie1 := hub.Parent.Cookie
	err := hub.login()
	assert.Nil(t, err)
	cookie2 := hub.Parent.Cookie
	assert.NotNil(t, cookie2)
	assert.NotEqual(t, cookie1.Value, cookie2.Value)
}

func TestCreateApp(t *testing.T) {
	hub := getHubAndLogin(t)

	err := hub.deleteApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.createApp())
	foundApps, err := hub.findApps(sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundApps))
	foundApp := foundApps[0]
	assert.Equal(t, hub.Parent.User, foundApp.Maintainer)
	assert.Equal(t, hub.App, foundApp.App)

	err = hub.createApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(409, "app already exists"), err.Error())

	assert.Nil(t, hub.deleteApp())
	foundApps, err = hub.findApps(sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundApps))
}

func TestUploadTag(t *testing.T) {
	hub := getHubAndLogin(t)

	err := hub.uploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag())

	err = hub.uploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(409, "tag already exists"), err.Error())

	tags, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, hub.Tag, tags[0].Name)

	assert.Nil(t, hub.deleteTag())
	tags, err = hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tags))

	err = hub.deleteTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "tag does not exist"), err.Error())
}

func TestLogin(t *testing.T) {
	hub := getHub()
	err := hub.login()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "user does not exist"), err.Error())
}

func TestChangePassword(t *testing.T) {
	hub := getHubAndLogin(t)

	hub.Parent.NewPassword = hub.Parent.Password + "x"

	assert.Nil(t, hub.changePassword())
	err := hub.login()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "incorrect username or password"), err.Error())

	hub.Parent.Password = hub.Parent.NewPassword
	hub.Parent.Cookie = nil
	err = hub.login()
	assert.Nil(t, err)
	assert.NotNil(t, hub.Parent.Cookie)
}

func TestRegistration(t *testing.T) {
	hub := getHub()
	assert.Nil(t, hub.registerUser())
	err := hub.registerUser()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(409, "user already exists"), err.Error())
}

func TestGetTagsUnhappyPath(t *testing.T) {
	hub := getHub()

	_, err := hub.getTags()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "user does not exist"), err.Error())

	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())
	_, err = hub.getTags()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.createApp())
	tagList, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tagList))
}

func TestLimitsForUploadsAndTagStorage(t *testing.T) {
	hub := getHubAndLogin(t)
	assert.Nil(t, hub.createApp())

	const oneMebibyteInBytes = 1024 * 1024
	hub.UploadContent = make([]byte, oneMebibyteInBytes+1)
	err := hub.uploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(413, "tag content too large, the limit is 1MB"), err.Error())

	hub.UploadContent = make([]byte, 1024*750) // Must be a little smaller than 1024*1024 since conversion to json makes it bigger.
	// Upload tags until we are just a little bit below the 10MiB storage limit.
	for i := 0; i < 13; i++ {
		hub.Tag = sampleTag + "." + strconv.Itoa(i)
		assert.Nil(t, hub.uploadTag())
	}
	// Tag whose upload exceeds the 10MiB storage limit.
	hub.Tag = sampleTag + ".x"
	err = hub.uploadTag()
	assert.NotNil(t, err)
	expectedMsg := "storage limit reached, you can't store more then 10MiB of tag content, currently used storage in bytes: 9984000/10485760 (95 percent)"
	assert.Equal(t, utils.GetErrMsg(413, expectedMsg), err.Error())
}

func TestLogout(t *testing.T) {
	hub := getHubAndLogin(t)
	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.logout())
	err := hub.createApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "cookie not found"), err.Error())
}

func TestGetAppList(t *testing.T) {
	hub := getHubAndLogin(t)
	apps, err := hub.GetApps()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(apps))
	assert.Nil(t, hub.createApp())
	apps, err = hub.GetApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, sampleApp, apps[0])
}
