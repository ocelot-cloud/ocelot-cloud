//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestFindAppsSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	defer hub.deleteUser()
	hub.SetCookieHeader = false
	hub.SetOriginHeader = false

	apps, err := hub.findApps("not-existing-app")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(apps))

}

func TestDownloadAppSecurity(t *testing.T) {
	// TODO duplication
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
	// TODO duplication
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
