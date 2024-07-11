//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
	"time"
)

func TestFileUploadDownload(t *testing.T) {
	hub := getHubAndLogin(t)

	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag())

	foundTags, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, sampleTag, foundTags[0])

	downloadedContent, err := hub.downloadApp()
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

	assert.Nil(t, hub.createApp())
	foundApps, err := hub.findApps(sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundApps))
	foundApp := foundApps[0]
	assert.Equal(t, hub.User, foundApp.User)
	assert.Equal(t, hub.App, foundApp.App)

	assert.Nil(t, hub.deleteApp())
	foundApps, err = hub.findApps(sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundApps))
}

func TestCreateTags(t *testing.T) {
	hub := getHubAndLogin(t)

	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag())
	tags, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, hub.Tag, tags[0])

	assert.Nil(t, hub.deleteTag())
	tags, err = hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tags))
}

func TestChangePassword(t *testing.T) {
	hub := getHubAndLogin(t)

	newPassword := "new-password"

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

// TODO assert that no other object should be send in body, should be nil, when IsCredentialsRequired == true
