//go:build acceptance

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/hub"
	"github.com/ocelot-cloud/shared/utils"
	"testing"
	"time"
)

func TestFindAppsSecurity(t *testing.T) {
	hub := GetHubAndLogin(t)

	hub.Parent.SetCookieHeader = false
	hub.Parent.SetOriginHeader = false

	_, err := hub.FindApps("notexistingapp")
	assert.Nil(t, err)

	testInputInvalidation(t, hub, "not-existing-app", AppField, FindApps)
}

func TestDownloadAppSecurity(t *testing.T) {
	client := GetHubAndLogin(t)

	assert.Nil(t, client.CreateApp())
	assert.Nil(t, client.UploadTag())

	client.Parent.SetCookieHeader = false
	client.Parent.SetOriginHeader = false
	downloadedContent, err := client.DownloadTag()
	assert.Nil(t, err)
	assert.Equal(t, hub.SampleTagFileContent, downloadedContent)

	testInputInvalidation(t, client, "invalid-user", UserField, DownloadTag)
	testInputInvalidation(t, client, "invalid-app", AppField, DownloadTag)
	testInputInvalidation(t, client, "invalid-tag", TagField, DownloadTag)
}

func TestGetTagsSecurity(t *testing.T) {
	client := GetHubAndLogin(t)

	assert.Nil(t, client.CreateApp())
	assert.Nil(t, client.UploadTag())

	client.Parent.SetCookieHeader = false
	client.Parent.SetOriginHeader = false
	tags, err := client.GetTags()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, hub.SampleTag, tags[0])

	testInputInvalidation(t, client, "invalid-user", UserField, GetTags)
	testInputInvalidation(t, client, "invalid-app", AppField, GetTags)
}

func TestRegisterSecurity(t *testing.T) {
	hub := hub.GetHub()
	hub.Parent.SetCookieHeader = false
	hub.Parent.SetOriginHeader = false
	testInputInvalidation(t, hub, "invalid-password-with-letter-ä", PasswordField, Register)
	testInputInvalidation(t, hub, "invalid-username", UserField, Register)
	testInputInvalidation(t, hub, "asd@asd.d", EmailField, Register)
}

func TestChangePasswordSecurity(t *testing.T) {
	client := GetHubAndLogin(t)

	client.Parent.NewPassword = hub.SamplePassword + "x"
	correctlyFormattedButNotMatchingPassword := hub.SamplePassword + "xy"
	client.Parent.Password = correctlyFormattedButNotMatchingPassword
	err := client.ChangePassword()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "incorrect username or password"), err.Error())
	client.Parent.Password = hub.SamplePassword

	testInputInvalidation(t, client, "invalid-password-ä", PasswordField, ChangePassword)
	testInputInvalidation(t, client, "invalid-password-ä", NewPasswordField, ChangePassword)
}

func TestLoginSecurity(t *testing.T) {
	client := hub.GetHub()
	err := client.RegisterUser()
	assert.Nil(t, err)

	assert.Nil(t, client.Parent.Cookie)
	assert.Nil(t, client.Login())
	assert.NotNil(t, client.Parent.Cookie)
	client.Parent.Cookie = nil

	testInputInvalidation(t, client, "invalid-user", UserField, Login)
	testInputInvalidation(t, client, "invalid-password-ä", PasswordField, Login)
	testInputInvalidation(t, client, "https:/only-single-slash-invalid-domain.de", OriginField, Login)

	correctlyFormattedButNotMatchingPassword := hub.SamplePassword + "x"
	client.Parent.Password = correctlyFormattedButNotMatchingPassword
	err = client.Login()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "incorrect username or password"), err.Error())
	client.Parent.Password = hub.SamplePassword
}

// TestDeleteUserSecurity is not necessary, since there are no further tests to conducted.

func TestCreateAppSecurity(t *testing.T) {
	hub := GetHubAndLogin(t)
	testInputInvalidation(t, hub, "invalid-app", AppField, CreateApp)
}

func TestDeleteAppSecurity(t *testing.T) {
	hub := GetHubAndLogin(t)
	testInputInvalidation(t, hub, "invalid-app", AppField, DeleteApp)
}

func TestUploadTagSecurity(t *testing.T) {
	hub := GetHubAndLogin(t)

	testInputInvalidation(t, hub, "invalid-app", AppField, UploadTag)
	testInputInvalidation(t, hub, "invalid-tag", TagField, UploadTag)
}

func TestDeleteTagSecurity(t *testing.T) {
	hub := GetHubAndLogin(t)

	testInputInvalidation(t, hub, "invalid-app", AppField, DeleteTag)
	testInputInvalidation(t, hub, "invalid-tag", TagField, DeleteTag)
}

func TestCookieExpirationAndRenewal(t *testing.T) {
	client := GetHubAndLogin(t)
	// There is some specific logic for this user in the production code when handling cookie.
	client.Parent.User = testUserWithExpiredCookie
	assert.Nil(t, client.RegisterUser())
	assert.Nil(t, client.Login())
	assert.True(t, time.Now().UTC().After(client.Parent.Cookie.Expires))
	err := client.CreateApp()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "cookie expired"), err.Error())
	client.Parent.User = hub.SampleUser

	// There is some specific logic for this user in the production code when handling cookie.
	client.Parent.User = testUserWithOldButNotExpiredCookie
	assert.Nil(t, client.RegisterUser())
	assert.Nil(t, client.Login())
	assert.True(t, time.Now().UTC().Before(client.Parent.Cookie.Expires))
	assert.True(t, time.Now().UTC().Add(48*time.Hour).After(client.Parent.Cookie.Expires))
	assert.Nil(t, client.CreateApp())
	assert.True(t, time.Now().UTC().AddDate(0, 0, 29).Before(client.Parent.Cookie.Expires))
	assert.True(t, time.Now().UTC().AddDate(0, 0, 31).After(client.Parent.Cookie.Expires))
	client.Parent.User = hub.SampleUser
}

func TestCookieAndHostProtection(t *testing.T) {
	hub := hub.GetHub()
	tests := []func() error{
		hub.DeleteUser,
		hub.CreateApp,
		hub.DeleteApp,
		hub.UploadTag,
		hub.DeleteTag,
		hub.ChangePassword,
		hub.CheckAuth,
	}
	for _, test := range tests {
		doCookieAndHostPolicyChecks(t, hub, test)
	}
}

func doCookieAndHostPolicyChecks(t *testing.T, client *hub.HubClient, operation func() error) {
	defer client.WipeData()
	assert.Nil(t, client.RegisterUser())
	assert.Nil(t, client.Login())

	client.Parent.SetCookieHeader = false
	client.Parent.SetOriginHeader = false

	err := operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "cookie not set in request"), err.Error())

	client.Parent.SetCookieHeader = true
	client.Parent.Cookie.Value = "some-invalid-cookie-value"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "invalid cookie"), err.Error())

	err = client.Login()
	assert.Nil(t, err)
	client.Parent.Origin = "http:/single-slash-invalid-origin"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "invalid origin"), err.Error())

	client.Parent.SetOriginHeader = true
	client.Parent.Origin = "http://valid-but-incorrect-origin.com"
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "origin not matching"), err.Error())

	client.Parent.Origin = hub.SampleOrigin
	validButNonExistentCookie := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	client.Parent.Cookie.Value = validButNonExistentCookie
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(401, "cookie not found"), err.Error())

	assert.Nil(t, client.Login())

	client.Parent.User = testUserWithExpiredCookie
	assert.Nil(t, client.RegisterUser())
	assert.Nil(t, client.Login())
	err = operation()
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "cookie expired"), err.Error())
	assert.True(t, time.Now().UTC().After(client.Parent.Cookie.Expires))
	client.Parent.User = hub.SampleUser
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

func testInputInvalidation(t *testing.T, hub *hub.HubClient, invalidValue string, fieldType FieldType, operation Operation) {
	originalValue := returnCurrentValueAndSetField(hub, fieldType, invalidValue)

	switch operation {
	case Register:
		assertInvalidInputError(t, hub.RegisterUser())
	case GetTags:
		_, err := hub.GetTags()
		assertInvalidInputError(t, err)
	case DownloadTag:
		_, err := hub.DownloadTag()
		assertInvalidInputError(t, err)
	case FindApps:
		_, err := hub.FindApps(hub.App)
		assertInvalidInputError(t, err)
	case ChangePassword:
		assertInvalidInputError(t, hub.ChangePassword())
	case Login:
		assertInvalidInputError(t, hub.Login())
	case DeleteApp:
		assertInvalidInputError(t, hub.DeleteApp())
	case UploadTag:
		assertInvalidInputError(t, hub.UploadTag())
	case DeleteTag:
		assertInvalidInputError(t, hub.DeleteTag())
	case CheckAuth:
		assertInvalidInputError(t, hub.CheckAuth())
	case CreateApp:
		assertInvalidInputError(t, hub.CreateApp())
	default:
		panic("Unsupported operation")
	}

	returnCurrentValueAndSetField(hub, fieldType, originalValue)
}

func assertInvalidInputError(t *testing.T, err error) {
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(400, "invalid input"), err.Error())
}

func returnCurrentValueAndSetField(hub *hub.HubClient, fieldType FieldType, value string) string {
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
