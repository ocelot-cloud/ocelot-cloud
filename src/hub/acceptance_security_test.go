//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

// TODO Input validation missing for the first three operations.

func TestFindAppsSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()
	hub.SetCookieHeader = false
	hub.SetOriginHeader = false

	_, err := hub.findApps("notexistingapp")
	assert.Nil(t, err)

	_, err = hub.findApps("not-existing-app")
	assert.NotNil(t, err)
	// TODO Resolve duplication
	assert.Equal(t, "Expected status code 200, but got 400. Response body: invalid app search term\n", err.Error())
}

func TestDownloadAppSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()
	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag(sampleTagFileContent))

	hub.SetCookieHeader = false
	hub.SetOriginHeader = false
	downloadedContent, err := hub.downloadApp()
	assert.Nil(t, err)
	assert.Equal(t, sampleTagFileContent, downloadedContent)
}

func TestGetTagsSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()
	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag(sampleTagFileContent))

	hub.SetCookieHeader = false
	hub.SetOriginHeader = false
	tags, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, sampleTag, tags[0])
}

func TestRegisterSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()
	// TODO cases: wrong password, input validation
}
