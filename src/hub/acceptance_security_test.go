//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
	"time"
)

func TestFindAppsSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	hub.SetCookieHeader = false
	hub.SetOriginHeader = false

	_, err := hub.findApps("notexistingapp")
	assert.Nil(t, err)

	testInputInvalidation(t, hub, "not-existing-app", AppField, FindApps)
}

func TestDownloadAppSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag())

	hub.SetCookieHeader = false
	hub.SetOriginHeader = false
	downloadedContent, err := hub.downloadApp()
	assert.Nil(t, err)
	assert.Equal(t, sampleTagFileContent, downloadedContent)

	testInputInvalidation(t, hub, "invalid-tag", TagField, DownloadApp)
}

func TestGetTagsSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag())

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
	hub.SetCookieHeader = false
	hub.SetOriginHeader = false
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

	correctlyFormattedButNotMatchingPassword := samplePassword + "x"
	hub.Password = correctlyFormattedButNotMatchingPassword
	err := hub.ChangeOrigin(sampleOrigin)
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "incorrect username or password"), err.Error())
	hub.Password = samplePassword

	testInputInvalidation(t, hub, "invalid-origin", OriginField, ChangeOrigin)
	testInputInvalidation(t, hub, "invalid-user", UserField, ChangeOrigin)
	testInputInvalidation(t, hub, "invalid-password-ä", PasswordField, ChangeOrigin)
}

func TestChangePasswordSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	hub.SetCookieHeader = false
	hub.SetOriginHeader = false

	assert.Nil(t, hub.ChangePassword(samplePassword))

	correctlyFormattedButNotMatchingPassword := samplePassword + "x"
	hub.Password = correctlyFormattedButNotMatchingPassword
	err := hub.ChangePassword(samplePassword)
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "incorrect username or password"), err.Error())
	hub.Password = samplePassword

	testInputInvalidation(t, hub, "invalid-user", UserField, ChangePassword)
	testInputInvalidation(t, hub, "invalid-password-ä", PasswordField, ChangePassword)

	oldPassword := hub.Password
	hub.Password = "invalid-old-password-ä"
	err = hub.ChangePassword("new-valid-password")
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(400, "invalid input"), err.Error())
	hub.Password = oldPassword
}

func TestLoginSecurity(t *testing.T) {
	hub := getHub()
	err := hub.registerUser()
	assert.Nil(t, err)

	assert.Nil(t, hub.Cookie)
	err = hub.login()
	assert.Nil(t, err)
	assert.NotNil(t, hub.Cookie)
	hub.Cookie = nil

	testInputInvalidation(t, hub, "invalid-user", UserField, Login)
	testInputInvalidation(t, hub, "invalid-password-ä", PasswordField, Login)

	correctlyFormattedButNotMatchingPassword := samplePassword + "x"
	hub.Password = correctlyFormattedButNotMatchingPassword
	err = hub.login()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "incorrect username or password"), err.Error())
	hub.Password = samplePassword
}

// TODO Not finished yet.
func TestDeleteUserSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	// TODO
	print(hub)
}

func TestCreateAppSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	// TODO
	print(hub)
}

func TestDeleteAppSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	hub.User += "x"
	err := hub.deleteApp()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "deletion of apps not belonging to you is not allowed"), err.Error())
	hub.User = sampleUser

	testInputInvalidation(t, hub, "invalid-user", UserField, DeleteApp)
	testInputInvalidation(t, hub, "invalid-app", AppField, DeleteApp)
}

func TestUploadTagSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	hub.User += "x"
	// TODO I think it is easier to get rid of that hub.field, and build it only when necessary.
	hub.TagFilename = getTagFileName(hub.User, hub.App, hub.Tag)
	err := hub.uploadTag()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "upload of tags not belonging to you is not allowed"), err.Error())
	hub.User = sampleUser
	hub.TagFilename = getTagFileName(sampleUser, sampleApp, sampleTag)

	testInputInvalidation(t, hub, "invalid-user", UserField, UploadTag)
	testInputInvalidation(t, hub, "invalid-app", AppField, UploadTag)
	testInputInvalidation(t, hub, "invalid-tag", TagField, UploadTag)
}

func TestDeleteTagSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	// TODO
	print(hub)
}

// TODO test: update expiration date when calling middleware
func TestCookieAndHostProtection(t *testing.T) {
	hub := getHubAndLogin(t)
	// There is some specific logic for this user in the production code when handling cookie.
	hub.User = "expirationtestuser" // TODO Abstract duplication
	assert.Nil(t, hub.registerUser())
	hub.User = sampleUser

	// TODO It would be cool, if I could abstract that even more like in the security policy collection.
	// TODO authorization checks missing for these functions: authenticated user can only apply this to entities he owns
	// TODO input validation missing
	doCookieAndHostPolicyChecks(t, hub, hub.deleteUser)
	doCookieAndHostPolicyChecks(t, hub, hub.createApp)
	doCookieAndHostPolicyChecks(t, hub, hub.deleteApp)
	doCookieAndHostPolicyChecks(t, hub, hub.uploadTag)
	doCookieAndHostPolicyChecks(t, hub, hub.deleteTag)
}

// TODO Re-check if those tests are also needed/applied by the other hub client functions above.
func doCookieAndHostPolicyChecks(t *testing.T, hub *HubClient, operation func() error) {
	hub.SetCookieHeader = false
	hub.SetOriginHeader = false

	err := operation()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "cookie not set in request"), err.Error())

	hub.SetCookieHeader = true
	hub.Cookie.Value = "some-invalid-cookie-value"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(400, "invalid cookie"), err.Error())

	err = hub.login()
	assert.Nil(t, err)
	hub.Origin = "http:/single-slash-invalid-origin"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(400, "invalid origin"), err.Error())

	hub.SetOriginHeader = true
	hub.Origin = "http://valid-but-incorrect-origin.com"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(400, "origin not matching"), err.Error())
	hub.Origin = sampleOrigin

	validButNonExistentCookie := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	hub.Cookie.Value = validButNonExistentCookie
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(404, "cookie not found"), err.Error())
	assert.Nil(t, hub.login())

	hub.User = "expirationtestuser"
	err = hub.login()
	assert.Nil(t, err)
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(400, "cookie expired"), err.Error())
	assert.True(t, time.Now().UTC().After(hub.Cookie.Expires))
	hub.User = sampleUser
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
		assertInvalidInputError(t, hub.registerUser())
	case GetTags:
		_, err := hub.getTags()
		assertInvalidInputError(t, err)
	case DownloadApp:
		_, err := hub.downloadApp()
		assertInvalidInputError(t, err)
	case FindApps:
		_, err := hub.findApps(hub.App)
		assertInvalidInputError(t, err)
	case ChangeOrigin:
		assertInvalidInputError(t, hub.ChangeOrigin(hub.Origin))
	case ChangePassword:
		assertInvalidInputError(t, hub.ChangePassword(hub.Password))
	case Login:
		assertInvalidInputError(t, hub.login())
	case DeleteApp:
		assertInvalidInputError(t, hub.deleteApp())
	case UploadTag:
		assertInvalidInputError(t, hub.uploadTag())
	default:
		panic("Unsupported operation")
	}

	returnCurrentValueAndSetField(hub, fieldType, originalValue)
}

func assertInvalidInputError(t *testing.T, err error) {
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(400, "invalid input"), err.Error())
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
		hub.TagFilename = getTagFileName(hub.User, hub.App, value)
	case EmailField:
		originalValue = hub.Email
		hub.Email = value
	case OriginField:
		originalValue = hub.Origin
		hub.Origin = value
	case AppField:
		originalValue = hub.App
		hub.App = value
		hub.TagFilename = getTagFileName(hub.User, hub.App, value)
	case TagField:
		originalValue = hub.Tag
		hub.Tag = value
		hub.TagFilename = getTagFileName(hub.User, hub.App, value)
	default:
		panic("Unsupported field type")
	}
	return originalValue
}
