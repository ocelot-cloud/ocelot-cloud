//go:build acceptance

package main

import (
	"bytes"
	"github.com/ocelot-cloud/shared/assert"
	"testing"
	"time"
)

func TestFileUploadDownload(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()

	assert.Nil(t, hub.createApp())
	fileContent := []byte("hello")
	fileBuffer := bytes.NewBuffer(fileContent)

	assert.Nil(t, hub.uploadFile(fileBuffer))

	/* TODO Implement
	asdf, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(asdf))
	// TODO verify fields of fileInfo
	*/

	downloadedContent, err := hub.downloadFile()
	assert.Nil(t, err)
	assert.Equal(t, fileContent, downloadedContent)
}

// TODO Test if cookie expiration date updates when making a successful request.

func TestCookie(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()

	assert.Equal(t, cookieName, hub.Cookie.Name)
	assert.True(t, getTimeIn30Days().Add(1*time.Second).After(hub.Cookie.Expires))
	assert.True(t, getTimeIn30Days().Add(-1*time.Second).Before(hub.Cookie.Expires))
	assert.Equal(t, 64, len(hub.Cookie.Value))

	cookie1 := hub.Cookie
	_, err := hub.login()
	assert.Nil(t, err)
	cookie2 := hub.Cookie
	assert.NotNil(t, cookie2)
	assert.NotEqual(t, cookie1.Value, cookie2.Value)
}

func TestCreateApp(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()
	assert.Nil(t, hub.createApp())
	foundApps, err := hub.findApps(sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundApps))
	foundApp := foundApps[0]
	assert.Equal(t, hub.Username, foundApp.Username)
	assert.Equal(t, hub.App, foundApp.AppName)
}

// TODO Can just be done, when I have a protected endpoint
func TestOriginPolicy(t *testing.T) {
	/*hub := getHub()
	form := hub.getRegistrationForm()
	fakeOrigin := "http://non-existing-subdomain.localhost:8082"
	assert.Nil(t, hub.registerUser(form))

	hub.SetOriginHeader = false
	err := hub.deleteUser()
	assert.NotNil(t, err)
	expected := fmt.Sprintf("Security policy does not allow this request without 'Origin' header")
	assert.Equal(t, expected, err.Error())


	form.Origin = fakeOrigin
	*/
	// TODO expected := fmt.Sprintf("Security policy does not allow requests from origin: %s", fakeOrigin)
}

// TODO assert that no other object should be send in body, should be nil, when IsCredentialsRequired == true
