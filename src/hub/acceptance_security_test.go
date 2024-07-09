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

	_, err := hub.findApps("notexistingapp")
	assert.Nil(t, err)

	testInputInvalidation(t, hub, "not-existing-app", AppField, FindApps)
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

	testInputInvalidation(t, hub, "invalid-tag", TagField, DownloadApp)
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

	testInputInvalidation(t, hub, "invalid-user", UserField, GetTags)
	testInputInvalidation(t, hub, "invalid-app", AppField, GetTags)
}

func TestRegisterSecurity(t *testing.T) {
	hub := getHub()
	testInputInvalidation(t, hub, "invalid-password-with-letter-ä", PasswordField, Register)
	testInputInvalidation(t, hub, "invalid-username", UserField, Register)
	testInputInvalidation(t, hub, "asd@asd.d", EmailField, Register)
	testInputInvalidation(t, hub, "https:/only-single-slash-invalid-domain.de", OriginField, Register)
}

func TestChangeOriginSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	hub.SetCookieHeader = false
	hub.SetOriginHeader = false
	assert.Nil(t, hub.ChangeOrigin(sampleOrigin))

	/*
		err := hub.ChangeOrigin("http:/single-slash-invalid-domain.com")
		assert.NotNil(t, err)
		assert.Equal(t, "asfg", err.Error())
	*/
}

type FieldType int

const (
	UserField FieldType = iota
	PasswordField
	EmailField
	OriginField
	AppField
	TagField
)

func testInputInvalidation(t *testing.T, hub *HubClient, invalidValue string, fieldType FieldType, operation Operation) {
	originalValue := returnCurrentValueAndSetField(hub, fieldType, invalidValue)

	switch operation {
	case Register:
		err := hub.registerUser()
		assert.NotNil(t, err)
		assert.Equal(t, "Expected status code 201, but got 400. Response body: invalid input\n", err.Error())
	case GetTags:
		_, err := hub.getTags()
		assert.NotNil(t, err)
		assert.Equal(t, "Expected status code 200, but got 400. Response body: invalid input\n", err.Error())
	case DownloadApp:
		_, err := hub.downloadApp()
		assert.NotNil(t, err)
		assert.Equal(t, "Expected status code 200, but got 400. Response body: file name is invalid\n", err.Error())
	case FindApps:
		_, err := hub.findApps(hub.App)
		assert.NotNil(t, err)
		// TODO Resolve duplication
		assert.Equal(t, "Expected status code 200, but got 400. Response body: invalid input\n", err.Error())
	}

	hub.deleteUser()
	returnCurrentValueAndSetField(hub, fieldType, originalValue)
}

func returnCurrentValueAndSetField(hub *HubClient, fieldType FieldType, value string) string {
	var originalValue string
	switch fieldType {
	case PasswordField:
		originalValue = hub.Password
		hub.Password = value
	case UserField:
		originalValue = hub.User
		hub.User = value
	case EmailField:
		originalValue = hub.Email
		hub.Email = value
	case OriginField:
		originalValue = hub.Origin
		hub.Origin = value
	case AppField:
		originalValue = hub.App
		hub.App = value
	case TagField:
		originalValue = hub.Tag
		hub.Tag = value
		hub.TagFilename = getTagFileName(sampleUser, sampleApp, value)
	}
	return originalValue
}
