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

	_, err = hub.findApps("not-existing-app")
	assert.NotNil(t, err)
	// TODO Resolve duplication
	assert.Equal(t, "Expected status code 200, but got 400. Response body: invalid input\n", err.Error())
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

	hub.Tag = "invalid-tag"
	hub.TagFilename = getTagFileName(sampleUser, sampleApp, hub.Tag)
	downloadedContent, err = hub.downloadApp()
	assert.NotNil(t, err)
	assert.Equal(t, "Expected status code 200, but got 400. Response body: file name is invalid\n", err.Error())
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

	hub.deleteUser()
	hub.User = "invalid-user"
	hub.registerUser()
	hub.createApp()
	_, err = hub.getTags()
	assert.NotNil(t, err)
	assert.Equal(t, "Expected status code 200, but got 400. Response body: invalid input\n", err.Error())
	hub.deleteUser()
	hub.User = sampleUser

	hub.App = "invalid-app"
	hub.registerUser()
	hub.createApp()
	_, err = hub.getTags()
	assert.NotNil(t, err)
	assert.Equal(t, "Expected status code 200, but got 400. Response body: invalid input\n", err.Error())

}

func TestRegisterSecurity(t *testing.T) {
	hub := getHub()
	defer hub.deleteUser()

	testInputInvalidation(t, hub, "invalid-password-with-letter-ä", samplePassword, PasswordField, Register)
	testInputInvalidation(t, hub, "invalid-username", sampleUser, UserField, Register)
	testInputInvalidation(t, hub, "asd@asd.d", sampleEmail, EmailField, Register)
	testInputInvalidation(t, hub, "https:/only-single-slash-invalid-domain.de", sampleOrigin, OriginField, Register)
}

type FieldType int

const (
	UserField FieldType = iota
	PasswordField
	EmailField
	OriginField
)

func testInputInvalidation(t *testing.T, hub *HubClient, invalidValue string, validValue string, fieldType FieldType, operation Operation) {
	setField(hub, fieldType, invalidValue)

	switch operation {
	case Register:
		err := hub.registerUser()
		assert.NotNil(t, err)
		assert.Equal(t, "Expected status code 201, but got 400. Response body: invalid input\n", err.Error())
	}

	hub.deleteUser()
	setField(hub, fieldType, validValue)
}

func setField(hub *HubClient, fieldType FieldType, value string) {
	switch fieldType {
	case PasswordField:
		hub.Password = value
	case UserField:
		hub.User = value
	case EmailField:
		hub.Email = value
	case OriginField:
		hub.Origin = value
	}
}
