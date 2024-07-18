//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
	"time"
)

func TestTagDownload(t *testing.T) {
	hub := getHub()

	_, err := hub.downloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(404, "user does not exist"), err.Error())

	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())

	_, err = hub.downloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.createApp())
	_, err = hub.downloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(404, "tag does not exist"), err.Error())

	assert.Nil(t, hub.uploadTag())
	foundTags, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, sampleTag, foundTags[0])

	downloadedContent, err := hub.downloadTag()
	assert.Nil(t, err)
	assert.Equal(t, sampleTagFileContent, downloadedContent)
}

// TODO Test if cookie expiration date updates when making a successful request.

func TestCookie(t *testing.T) {
	hub := getHubAndLogin(t)

	assert.Equal(t, cookieName, hub.Cookie.Name)
	assert.True(t, getTimeIn30Days().Add(1*time.Second).After(hub.Cookie.Expires))
	assert.True(t, getTimeIn30Days().Add(-1*time.Second).Before(hub.Cookie.Expires))
	assert.Equal(t, 64, len(hub.Cookie.Value))

	cookie1 := hub.Cookie
	err := hub.login()
	assert.Nil(t, err)
	cookie2 := hub.Cookie
	assert.NotNil(t, cookie2)
	assert.NotEqual(t, cookie1.Value, cookie2.Value)
}

func TestCreateApp(t *testing.T) {
	hub := getHubAndLogin(t)

	err := hub.deleteApp()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.createApp())
	foundApps, err := hub.findApps(sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundApps))
	foundApp := foundApps[0]
	assert.Equal(t, hub.User, foundApp.User)
	assert.Equal(t, hub.App, foundApp.App)

	err = hub.createApp()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(409, "app already exists"), err.Error())

	assert.Nil(t, hub.deleteApp())
	foundApps, err = hub.findApps(sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundApps))
}

func TestUploadTag(t *testing.T) {
	hub := getHubAndLogin(t)

	err := hub.uploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag())

	err = hub.uploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(409, "tag already exists"), err.Error())

	tags, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, hub.Tag, tags[0])

	assert.Nil(t, hub.deleteTag())
	tags, err = hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tags))

	err = hub.deleteTag()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(404, "tag does not exist"), err.Error())
}

func TestChangePassword(t *testing.T) {
	hub := getHubAndLogin(t)

	newPassword := hub.Password + "x"

	assert.Nil(t, hub.ChangePassword(newPassword))
	err := hub.login()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "incorrect username or password"), err.Error())

	hub.Password = newPassword
	err = hub.login()
	assert.Nil(t, err)
	assert.NotNil(t, hub.Cookie)
}

func TestChangeOrigin(t *testing.T) {
	hub := getHubAndLogin(t)

	newOrigin := "http://wrong-origin.de"

	assert.Nil(t, hub.ChangeOrigin(newOrigin))
	err := hub.createApp()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(400, "origin not matching"), err.Error())

	hub.Origin = newOrigin
	err = hub.createApp()
	assert.Nil(t, err)
}

func TestRegistration(t *testing.T) {
	hub := getHub()
	assert.Nil(t, hub.registerUser())
	err := hub.registerUser()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(409, "user already exists"), err.Error())
}

// TODO test case: user does not exist? -> getTagList?
