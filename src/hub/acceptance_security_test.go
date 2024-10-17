//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"testing"
	"time"
)

func TestFindAppsSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	hub.Parent.SetCookieHeader = false
	hub.Parent.SetOriginHeader = false

	_, err := hub.findApps("notexistingapp")
	assert.Nil(t, err)

	testInputInvalidation(t, hub, "not-existing-app", AppField, FindApps)
}

func TestDownloadAppSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag())

	hub.Parent.SetCookieHeader = false
	hub.Parent.SetOriginHeader = false
	downloadedContent, err := hub.downloadTag()
	assert.Nil(t, err)
	assert.Equal(t, sampleTagFileContent, downloadedContent)
}

func TestGetTagsSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	assert.Nil(t, hub.createApp())
	assert.Nil(t, hub.uploadTag())

	hub.Parent.SetCookieHeader = false
	hub.Parent.SetOriginHeader = false
	tags, err := hub.getTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, sampleTag, tags[0].Name)
}

func TestRegisterSecurity(t *testing.T) {
	hub := getHub()
	hub.Parent.SetCookieHeader = false
	hub.Parent.SetOriginHeader = false
	testInputInvalidation(t, hub, "invalid-password-with-letter-ä", PasswordField, Register)
	testInputInvalidation(t, hub, "invalid-username", UserField, Register)
	testInputInvalidation(t, hub, "asd@asd.d", EmailField, Register)
}

func TestChangePasswordSecurity(t *testing.T) {
	hub := getHubAndLogin(t)

	hub.Parent.NewPassword = samplePassword + "x"
	correctlyFormattedButNotMatchingPassword := samplePassword + "xy"
	hub.Parent.Password = correctlyFormattedButNotMatchingPassword
	err := hub.changePassword()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "incorrect username or password"), err.Error())
	hub.Parent.Password = samplePassword

	testInputInvalidation(t, hub, "invalid-password-ä", PasswordField, ChangePassword)
	testInputInvalidation(t, hub, "invalid-password-ä", NewPasswordField, ChangePassword)
}

func TestLoginSecurity(t *testing.T) {
	hub := getHub()
	err := hub.registerUser()
	assert.Nil(t, err)

	assert.Nil(t, hub.Parent.Cookie)
	assert.Nil(t, hub.login())
	assert.NotNil(t, hub.Parent.Cookie)
	hub.Parent.Cookie = nil

	testInputInvalidation(t, hub, "invalid-user", UserField, Login)
	testInputInvalidation(t, hub, "invalid-password-ä", PasswordField, Login)
	testInputInvalidation(t, hub, "https:/only-single-slash-invalid-domain.de", OriginField, Login)

	correctlyFormattedButNotMatchingPassword := samplePassword + "x"
	hub.Parent.Password = correctlyFormattedButNotMatchingPassword
	err = hub.login()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "incorrect username or password"), err.Error())
	hub.Parent.Password = samplePassword
}

// TestDeleteUserSecurity is not necessary, since there are no further tests to conducted.

func TestCreateAppSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	testInputInvalidation(t, hub, "invalid-app", AppField, CreateApp)
}

func TestUploadTagSecurity(t *testing.T) {
	hub := getHubAndLogin(t)
	testInputInvalidation(t, hub, "invalid-tag", TagField, UploadTag)
}

func TestCookieExpirationAndRenewal(t *testing.T) {
	hub := getHubAndLogin(t)
	// There is some specific logic for this user in the production code when handling cookie.
	hub.Parent.User = testUserWithExpiredCookie
	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())
	assert.True(t, time.Now().UTC().After(hub.Parent.Cookie.Expires))
	err := hub.createApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "cookie expired"), err.Error())
	hub.Parent.User = sampleUser

	// There is some specific logic for this user in the production code when handling cookie.
	hub.Parent.User = testUserWithOldButNotExpiredCookie
	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())
	assert.True(t, time.Now().UTC().Before(hub.Parent.Cookie.Expires))
	assert.True(t, time.Now().UTC().Add(48*time.Hour).After(hub.Parent.Cookie.Expires))
	assert.Nil(t, hub.createApp())
	assert.True(t, time.Now().UTC().AddDate(0, 0, 29).Before(hub.Parent.Cookie.Expires))
	assert.True(t, time.Now().UTC().AddDate(0, 0, 31).After(hub.Parent.Cookie.Expires))
	hub.Parent.User = sampleUser
}

func TestCookieAndHostProtection(t *testing.T) {
	hub := getHub()
	tests := []func() error{
		hub.deleteUser,
		hub.createApp,
		hub.deleteApp,
		hub.uploadTag,
		hub.deleteTag,
		hub.changePassword,
		hub.checkAuth,
	}
	for _, test := range tests {
		doCookieAndHostPolicyChecks(t, hub, test)
	}
}

func doCookieAndHostPolicyChecks(t *testing.T, hub *HubClient, operation func() error) {
	defer hub.wipeData()
	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())

	hub.Parent.SetCookieHeader = false
	hub.Parent.SetOriginHeader = false

	err := operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "cookie not set in request"), err.Error())

	hub.Parent.SetCookieHeader = true
	hub.Parent.Cookie.Value = "some-invalid-cookie-value"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "invalid cookie"), err.Error())

	err = hub.login()
	assert.Nil(t, err)
	hub.Parent.Origin = "http:/single-slash-invalid-origin"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "invalid origin"), err.Error())

	hub.Parent.SetOriginHeader = true
	hub.Parent.Origin = "http://valid-but-incorrect-origin.com"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "origin not matching"), err.Error())

	hub.Parent.Origin = sampleOrigin
	validButNonExistentCookie := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	hub.Parent.Cookie.Value = validButNonExistentCookie
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "cookie not found"), err.Error())

	assert.Nil(t, hub.login())

	hub.Parent.User = testUserWithExpiredCookie
	assert.Nil(t, hub.registerUser())
	assert.Nil(t, hub.login())
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "cookie expired"), err.Error())
	assert.True(t, time.Now().UTC().After(hub.Parent.Cookie.Expires))
	hub.Parent.User = sampleUser
}

type FieldType int

const (
	UserField FieldType = iota
	PasswordField
	NewPasswordField
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
	assert.Equal(t, utils.GetErrMsg(400, "invalid input"), err.Error())
}

func returnCurrentValueAndSetField(hub *HubClient, fieldType FieldType, value string) string {
	var originalValue string
	switch fieldType {
	case PasswordField:
		originalValue = hub.Parent.Password
		hub.Parent.Password = value
	case NewPasswordField:
		originalValue = hub.Parent.NewPassword
		hub.Parent.NewPassword = value
	case UserField:
		originalValue = hub.Parent.User
		hub.Parent.User = value
	case EmailField:
		originalValue = hub.Email
		hub.Email = value
	case OriginField:
		originalValue = hub.Parent.Origin
		hub.Parent.Origin = value
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
