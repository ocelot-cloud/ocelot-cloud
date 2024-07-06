//go:build acceptance

package main

import (
	"bytes"
	"github.com/ocelot-cloud/shared/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestFileUploadDownload(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()
	assert.Nil(t, hub.createApp())
	uploadFilename := strings.Join([]string{hub.Username, hub.App, hub.Tag}, "_") + ".tar.gz"

	fileContent := []byte("hello")
	fileBuffer := bytes.NewBuffer(fileContent)

	responseCode, err := hub.uploadFile(rootUrl+tagPath, uploadFilename, fileBuffer)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, responseCode)

	downloadURL := rootUrl + downloadPath + uploadFilename
	downloadedContent, err := hub.downloadFile(downloadURL)
	assert.Nil(t, err)
	assert.Equal(t, fileContent, downloadedContent)
}

// TODO Test if cookie expiration date updates when making a successful request.

func TestCookie(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()

	assert.NotNil(t, hub.Cookie)
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
	app := foundApps[0]
	assert.Equal(t, hub.Username, app.Username)
	assert.Equal(t, hub.App, app.AppName)
}

func getHubAndLogin(t *testing.T) *HubClient {
	hub := getHub()
	form := hub.getRegistrationForm()
	assert.Nil(t, hub.registerUser(form))

	cookie, err := hub.login()
	assert.Nil(t, err)
	hub.Cookie = cookie
	return hub
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
