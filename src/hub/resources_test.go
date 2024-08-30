package main

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"testing"
)

var (
	sampleUser           = "myuser"
	sampleApp            = "myapp"
	sampleTag            = "v0.0.1"
	sampleTagFileContent = "hello"
	sampleEmail          = "testuser@example.com"
	samplePassword       = "mypassword"
	sampleOrigin         = rootUrl
	sampleForm           = &RegistrationForm{
		sampleUser,
		samplePassword,
		sampleEmail,
	}
)

type HubClient struct {
	Parent        utils.ComponentClient
	Email         string
	App           string
	Tag           string
	UploadContent []byte
}

type Operation int

const (
	FindApps Operation = iota
	DownloadTag
	Register
	ChangePassword
	Login
	CreateApp
	DeleteApp
	UploadTag
	DeleteTag
	GetTags
	CheckAuth
)

func getRegistrationForm(hub *HubClient) *RegistrationForm {
	return &RegistrationForm{
		User:     hub.Parent.User,
		Password: hub.Parent.Password,
		Email:    hub.Email,
	}
}

func getHub() *HubClient {
	hub := &HubClient{
		Parent: utils.ComponentClient{
			User:            sampleUser,
			Password:        samplePassword,
			Origin:          rootUrl,
			SetOriginHeader: true,
			SetCookieHeader: true,
			RootUrl:         rootUrl,
		},

		Email:         sampleEmail,
		App:           sampleApp,
		Tag:           sampleTag,
		UploadContent: []byte(sampleTagFileContent),
	}
	hub.wipeData()
	return hub
}

func (h *HubClient) registerUser() error {
	form := getRegistrationForm(h)
	_, err := h.Parent.DoRequest(registrationPath, form, "")
	return err
}

func (h *HubClient) login() error {
	creds := LoginCredentials{
		User:     h.Parent.User,
		Password: h.Parent.Password,
		Origin:   h.Parent.Origin,
	}

	resp, err := h.Parent.DoRequestWithFullResponse(loginPath, creds, "")
	if err != nil {
		return err
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		return fmt.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	h.Parent.Cookie = cookies[0]
	return nil
}

func (h *HubClient) deleteUser() error {
	_, err := h.Parent.DoRequest(deleteUserPath, nil, "")
	return err
}

func (h *HubClient) createApp() error {
	_, err := h.Parent.DoRequest(appCreationPath, utils.SingleString{h.App}, "")
	return err
}

func (h *HubClient) findApps(searchTerm string) ([]UserAndApp, error) {
	result, err := h.Parent.DoRequest(searchAppsPath, utils.SingleString{searchTerm}, "")
	if err != nil {
		return nil, err
	}

	apps, err := unpackResponse[[]UserAndApp](result)
	if err != nil {
		return nil, err
	}

	return *apps, nil
}

func (h *HubClient) GetApps() ([]string, error) {
	result, err := h.Parent.DoRequest(appGetListPath, nil, "")
	if err != nil {
		return nil, err
	}

	apps, err := unpackResponse[[]string](result)
	if err != nil {
		return nil, err
	}

	return *apps, nil
}

func (h *HubClient) uploadTag() error {
	tapUpload := &TagUpload{
		App:     h.App,
		Tag:     h.Tag,
		Content: h.UploadContent,
	}
	_, err := h.Parent.DoRequest(tagUploadPath, tapUpload, "")
	return err
}

func (h *HubClient) downloadTag() (string, error) {
	tagInfo := &TagInfo{
		User: h.Parent.User,
		App:  h.App,
		Tag:  h.Tag,
	}

	result, err := h.Parent.DoRequest(downloadPath, tagInfo, "")
	if err != nil {
		return "", err
	}

	downloadedContent, ok := result.([]byte)
	if !ok {
		return "", fmt.Errorf("Failed to assert result to []byte")
	}

	return string(downloadedContent), nil
}

func (h *HubClient) getTags() ([]string, error) {
	usernameAndApp := &UserAndApp{
		User: h.Parent.User,
		App:  h.App,
	}

	result, err := h.Parent.DoRequest(getTagsPath, usernameAndApp, "")
	if err != nil {
		return nil, err
	}

	tags, err := unpackResponse[[]string](result)
	if err != nil {
		return nil, err
	}

	return *tags, nil
}

func unpackResponse[T any](object interface{}) (*T, error) {
	respBody, ok := object.([]byte)
	if !ok {
		return nil, fmt.Errorf("Failed to assert result to []byte")
	}

	var result T
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v", err)
	}
	return &result, nil
}

func (h *HubClient) deleteTag() error {
	tagInfo := &AppAndTag{
		App: h.App,
		Tag: h.Tag,
	}
	_, err := h.Parent.DoRequest(tagDeletePath, tagInfo, "")
	return err
}

func (h *HubClient) deleteApp() error {
	_, err := h.Parent.DoRequest(appDeletePath, utils.SingleString{h.App}, "")
	return err
}

func (h *HubClient) changePassword() error {
	form := utils.ChangePasswordForm{
		OldPassword: h.Parent.Password,
		NewPassword: h.Parent.NewPassword,
	}

	_, err := h.Parent.DoRequest(changePasswordPath, form, "")
	return err
}

func getHubAndLogin(t *testing.T) *HubClient {
	hub := getHub()
	assert.Nil(t, hub.registerUser())
	err := hub.login()
	assert.Nil(t, err)
	return hub
}

func (h *HubClient) wipeData() error {
	_, err := h.Parent.DoRequest(wipeDataPath, nil, "")
	return err
}

func (h *HubClient) logout() error {
	_, err := h.Parent.DoRequest(logoutPath, nil, "")
	return err
}

func (h *HubClient) checkAuth() error {
	_, err := h.Parent.DoRequest(authCheckPath, nil, "")
	return err
}
