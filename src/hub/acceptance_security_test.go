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
	downloadedContent, err := hub.downloadTag()
	assert.Nil(t, err)
	assert.Equal(t, sampleTagFileContent, downloadedContent)

	testInputInvalidation(t, hub, "invalid-user", UserField, DownloadTag)
	testInputInvalidation(t, hub, "invalid-app", AppField, DownloadTag)
	testInputInvalidation(t, hub, "invalid-tag", TagField, DownloadTag)
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
}

func TestChangePasswordSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	hub.NewPassword = samplePassword + "x"
	correctlyFormattedButNotMatchingPassword := samplePassword + "xy"
	hub.Password = correctlyFormattedButNotMatchingPassword
	err := hub.changePassword()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "incorrect username or password"), err.Error())
	hub.Password = samplePassword

	testInputInvalidation(t, hub, "invalid-password-ä", PasswordField, ChangePassword)
	// TODO New password field check
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
	testInputInvalidation(t, hub, "https:/only-single-slash-invalid-domain.de", OriginField, Login)

	correctlyFormattedButNotMatchingPassword := samplePassword + "x"
	hub.Password = correctlyFormattedButNotMatchingPassword
	err = hub.login()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(401, "incorrect username or password"), err.Error())
	hub.Password = samplePassword
}

// TestDeleteUserSecurity is not necessary, since there are no further tests to conducted.

func TestCreateAppSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	testInputInvalidation(t, hub, "invalid-app", AppField, CreateApp)
}

func TestDeleteAppSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	testInputInvalidation(t, hub, "invalid-app", AppField, DeleteApp)
}

func TestUploadTagSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	testInputInvalidation(t, hub, "invalid-app", AppField, UploadTag)
	testInputInvalidation(t, hub, "invalid-tag", TagField, UploadTag)
}

func TestDeleteTagSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	testInputInvalidation(t, hub, "invalid-app", AppField, DeleteTag)
	testInputInvalidation(t, hub, "invalid-tag", TagField, DeleteTag)
}

// There is some specific logic for this user in the production code when handling cookie.
func TestCookieExpirationAndRenewal(t *testing.T) {
	hub := getHubAndLogin(t)
	hub.User = testUserWithExpiredCookie
	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())
	assert.True(t, time.Now().UTC().After(hub.Cookie.Expires))
	err := hub.createApp()
	assert.NotNil(t, err)
	assert.Equal(t, getErrMsg(400, "cookie expired"), err.Error())
	hub.User = sampleUser

	hub.User = testUserWithOldButNotExpiredCookie
	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())
	assert.True(t, time.Now().UTC().Before(hub.Cookie.Expires))
	assert.True(t, time.Now().UTC().Add(48*time.Hour).After(hub.Cookie.Expires))
	assert.Nil(t, hub.createApp())
	assert.True(t, time.Now().UTC().AddDate(0, 0, 29).Before(hub.Cookie.Expires))
	assert.True(t, time.Now().UTC().AddDate(0, 0, 31).After(hub.Cookie.Expires))
	hub.User = sampleUser
}

func TestCookieAndHostProtection(t *testing.T) {
	hub := getHub()
	doCookieAndHostPolicyChecks(t, hub, hub.deleteUser)
	doCookieAndHostPolicyChecks(t, hub, hub.createApp)
	doCookieAndHostPolicyChecks(t, hub, hub.deleteApp)
	doCookieAndHostPolicyChecks(t, hub, hub.uploadTag)
	doCookieAndHostPolicyChecks(t, hub, hub.deleteTag)
	doCookieAndHostPolicyChecks(t, hub, hub.checkAuth)
	doCookieAndHostPolicyChecks(t, hub, hub.changePassword)
}

func doCookieAndHostPolicyChecks(t *testing.T, hub *HubClient, operation func() error) {
	defer hub.wipeData()
	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())

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

	hub.User = testUserWithExpiredCookie
	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())
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
	case DownloadTag:
		_, err := hub.downloadTag()
		assertInvalidInputError(t, err)
	case FindApps:
		_, err := hub.findApps(hub.App)
		assertInvalidInputError(t, err)
	case ChangePassword:
		assertInvalidInputError(t, hub.changePassword())
	case Login:
		assertInvalidInputError(t, hub.login())
	case DeleteApp:
		assertInvalidInputError(t, hub.deleteApp())
	case UploadTag:
		assertInvalidInputError(t, hub.uploadTag())
	case DeleteTag:
		assertInvalidInputError(t, hub.deleteTag())
	case CheckAuth:
		assertInvalidInputError(t, hub.checkAuth())
	case CreateApp:
		assertInvalidInputError(t, hub.createApp())
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
	default:
		panic("Unsupported field type")
	}
	return originalValue
}
