//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/hub"
	"github.com/ocelot-cloud/shared/utils"
	"strconv"
	"testing"
	"time"
)

func TestTagDownload(t *testing.T) {
	client := hub.GetHub()

	_, err := client.DownloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "user does not exist"), err.Error())

	assert.Nil(t, client.RegisterUser())
	assert.Nil(t, client.Login())

	_, err = client.DownloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, client.CreateApp())
	_, err = client.DownloadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "tag does not exist"), err.Error())

	assert.Nil(t, client.UploadTag())
	foundTags, err := client.GetTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, hub.SampleTag, foundTags[0])

	downloadedContent, err := client.DownloadTag()
	assert.Nil(t, err)
	assert.Equal(t, hub.SampleTagFileContent, downloadedContent)
}

func TestCookie(t *testing.T) {
	hub := GetHubAndLogin(t)

	assert.Equal(t, cookieName, hub.Parent.Cookie.Name)
	assert.True(t, utils.GetTimeIn30Days().Add(1*time.Second).After(hub.Parent.Cookie.Expires))
	assert.True(t, utils.GetTimeIn30Days().Add(-1*time.Second).Before(hub.Parent.Cookie.Expires))
	assert.Equal(t, 64, len(hub.Parent.Cookie.Value))

	cookie1 := hub.Parent.Cookie
	err := hub.Login()
	assert.Nil(t, err)
	cookie2 := hub.Parent.Cookie
	assert.NotNil(t, cookie2)
	assert.NotEqual(t, cookie1.Value, cookie2.Value)
}

func TestCreateApp(t *testing.T) {
	client := GetHubAndLogin(t)

	err := client.DeleteApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, client.CreateApp())
	foundApps, err := client.FindApps(hub.SampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundApps))
	foundApp := foundApps[0]
	assert.Equal(t, client.Parent.User, foundApp.User)
	assert.Equal(t, client.App, foundApp.App)

	err = client.CreateApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(409, "app already exists"), err.Error())

	assert.Nil(t, client.DeleteApp())
	foundApps, err = client.FindApps(hub.SampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundApps))
}

func TestUploadTag(t *testing.T) {
	hub := GetHubAndLogin(t)

	err := hub.UploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.CreateApp())
	assert.Nil(t, hub.UploadTag())

	err = hub.UploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(409, "tag already exists"), err.Error())

	tags, err := hub.GetTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, hub.Tag, tags[0])

	assert.Nil(t, hub.DeleteTag())
	tags, err = hub.GetTags()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tags))

	err = hub.DeleteTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "tag does not exist"), err.Error())
}

func TestLogin(t *testing.T) {
	hub := hub.GetHub()
	err := hub.Login()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "user does not exist"), err.Error())
}

func TestChangePassword(t *testing.T) {
	hub := GetHubAndLogin(t)

	hub.Parent.NewPassword = hub.Parent.Password + "x"

	assert.Nil(t, hub.ChangePassword())
	err := hub.Login()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "incorrect username or password"), err.Error())

	hub.Parent.Password = hub.Parent.NewPassword
	hub.Parent.Cookie = nil
	err = hub.Login()
	assert.Nil(t, err)
	assert.NotNil(t, hub.Parent.Cookie)
}

func TestRegistration(t *testing.T) {
	hub := hub.GetHub()
	assert.Nil(t, hub.RegisterUser())
	err := hub.RegisterUser()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(409, "user already exists"), err.Error())
}

func TestGetTagsUnhappyPath(t *testing.T) {
	hub := hub.GetHub()

	_, err := hub.GetTags()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "user does not exist"), err.Error())

	assert.Nil(t, hub.RegisterUser())
	assert.Nil(t, hub.Login())
	_, err = hub.GetTags()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(404, "app does not exist"), err.Error())

	assert.Nil(t, hub.CreateApp())
	tagList, err := hub.GetTags()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tagList))
}

func TestLimitsForUploadsAndTagStorage(t *testing.T) {
	client := GetHubAndLogin(t)
	assert.Nil(t, client.CreateApp())

	const oneMebibyteInBytes = 1024 * 1024
	client.UploadContent = make([]byte, oneMebibyteInBytes+1)
	err := client.UploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(413, "tag content too large, the limit is 1MB"), err.Error())

	client.UploadContent = make([]byte, 1024*750) // Must be a little smaller than 1024*1024 since conversion to json makes it bigger.
	// Upload tags until we are just a little bit below the 10MiB storage limit.
	for i := 0; i < 13; i++ {
		client.Tag = hub.SampleTag + "." + strconv.Itoa(i)
		assert.Nil(t, client.UploadTag())
	}
	// Tag whose upload exceeds the 10MiB storage limit.
	client.Tag = hub.SampleTag + ".x"
	err = client.UploadTag()
	assert.NotNil(t, err)
	expectedMsg := "storage limit reached, you can't store more then 10MiB of tag content, currently used storage in bytes: 9984000/10485760 (95 percent)"
	assert.Equal(t, utils.GetErrMsg(413, expectedMsg), err.Error())
}

func TestLogout(t *testing.T) {
	hub := GetHubAndLogin(t)
	assert.Nil(t, hub.CreateApp())
	assert.Nil(t, hub.Logout())
	err := hub.CreateApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "cookie not found"), err.Error())
}

func TestGetAppList(t *testing.T) {
	client := GetHubAndLogin(t)
	apps, err := client.GetApps()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(apps))
	assert.Nil(t, client.CreateApp())
	apps, err = client.GetApps()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, hub.SampleApp, apps[0])
}
