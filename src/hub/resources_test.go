package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"io"
	"net/http"
	"strings"
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
		sampleOrigin,
	}
)

type HubClient struct {
	User            string
	Password        string
	Origin          string
	Email           string
	App             string
	Cookie          *http.Cookie
	Tag             string
	SetOriginHeader bool
	SetCookieHeader bool
	UploadContent   []byte
}

type Operation int

const (
	FindApps Operation = iota
	DownloadApp
	Register
	ChangeOrigin
	ChangePassword
	Login
	DeleteUser
	CreateApp
	DeleteApp
	UploadTag
	DeleteTag
	GetTags
	WipeData
	Logout
)

func getRegistrationForm(hub *HubClient) *RegistrationForm {
	return &RegistrationForm{
		User:     hub.User,
		Password: hub.Password,
		Origin:   hub.Origin,
		Email:    hub.Email,
	}
}

func getHub() *HubClient {
	hub := &HubClient{
		User:            sampleUser,
		Password:        samplePassword,
		Origin:          rootUrl,
		Email:           sampleEmail,
		App:             sampleApp,
		Tag:             sampleTag,
		SetOriginHeader: true,
		SetCookieHeader: true,
		UploadContent:   []byte(sampleTagFileContent),
	}
	hub.wipeData()
	return hub
}

func (h *HubClient) registerUser() error {
	form := getRegistrationForm(h)
	_, err := h.doRequest(registrationPath, form, "", "POST", Register)
	return err
}

func (h *HubClient) login() error {
	creds := LoginCredentials{
		User:     h.User,
		Password: h.Password,
	}

	result, err := h.doRequest(loginPath, creds, "", "GET", Login)
	if err != nil {
		return err
	}

	resp, ok := result.(*http.Response)
	if !ok {
		return fmt.Errorf("Failed to assert result to *http.Response")
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		return fmt.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	h.Cookie = cookies[0]
	return nil
}

func (h *HubClient) deleteUser() error {
	_, err := h.doRequest(userPath, nil, "", "DELETE", DeleteUser)
	return err
}

func (h *HubClient) doRequest(path string, payload interface{}, expectedMessage string, method string, operation Operation) (interface{}, error) {
	url := rootUrl + path

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal payload: %v", err)
	}
	payloadReader := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	setCookieAndOriginHeaders(req, h)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %v", err)
	}

	respBody, err := assertOkStatusAndExtractBody(resp)
	if err != nil {
		return nil, err
	}

	responseMessage, _ := strings.CutSuffix(string(respBody), "\n")
	if expectedMessage != "" && expectedMessage != responseMessage {
		return nil, fmt.Errorf("Expected response message '%s', got '%s'", expectedMessage, responseMessage)
	}

	if len(resp.Cookies()) == 1 {
		h.Cookie = resp.Cookies()[0]
	}

	if operation == Login {
		return resp, nil
	} else if operation == FindApps || operation == GetTags || operation == DownloadApp {
		return respBody, nil
	} else {
		return nil, nil
	}
}

func setCookieAndOriginHeaders(req *http.Request, h *HubClient) {
	if h.SetOriginHeader {
		req.Header.Set(OriginHeader, h.Origin)
	}
	if h.SetCookieHeader && h.Cookie != nil {
		req.AddCookie(h.Cookie)
	}
}

func assertOkStatusAndExtractBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	var bodyBuffer bytes.Buffer
	teeReader := io.TeeReader(resp.Body, &bodyBuffer)

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(teeReader)
		if err != nil {
			return nil, fmt.Errorf("Expected status code 200, but got %d. Also failed to read response body: %v", resp.StatusCode, err)
		}
		errorMessage := getErrMsg(resp.StatusCode, string(respBody))
		trimmedStr := strings.TrimSuffix(errorMessage, "\n")
		return nil, fmt.Errorf(trimmedStr)
	}

	respBody, err := io.ReadAll(teeReader)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	return respBody, nil
}

func getErrMsg(actualStatusCode int, respBodyMsg string) string {
	var msg string
	if respBodyMsg == "" {
		msg = ""
	} else {
		msg = fmt.Sprintf(" Response body: %s", respBodyMsg)
	}
	return fmt.Sprintf("Expected status code 200, but got %d.%s", actualStatusCode, msg)
}

func (h *HubClient) createApp() error {
	_, err := h.doRequest(appPath, SingleString{h.App}, "", "POST", CreateApp)
	return err
}

func (h *HubClient) findApps(searchTerm string) ([]UserAndApp, error) {
	result, err := h.doRequest(appPath, SingleString{searchTerm}, "", "GET", FindApps)
	if err != nil {
		return nil, err
	}

	apps, err := unpackResponse[[]UserAndApp](result)
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
	_, err := h.doRequest(tagPath, tapUpload, "", "POST", UploadTag)
	return err
}

func (h *HubClient) downloadTag() (string, error) {
	tagInfo := &TagInfo{
		User: h.User,
		App:  h.App,
		Tag:  h.Tag,
	}

	result, err := h.doRequest(downloadPath, tagInfo, "", "GET", DownloadApp)
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
		User: h.User,
		App:  h.App,
	}

	result, err := h.doRequest(tagPath, usernameAndApp, "", "GET", GetTags)
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
	_, err := h.doRequest(tagPath, tagInfo, "", "DELETE", DeleteTag)
	return err
}

func (h *HubClient) deleteApp() error {
	_, err := h.doRequest(appPath, SingleString{h.App}, "", "DELETE", DeleteApp)
	return err
}

func (h *HubClient) ChangePassword(newPassword string) error {
	form := ChangePasswordForm{
		User:        h.User,
		OldPassword: h.Password,
		NewPassword: newPassword,
	}

	_, err := h.doRequest(changePasswordPath, form, "", "POST", ChangePassword)
	return err
}

func (h *HubClient) ChangeOrigin(newOrigin string) error {
	form := ChangeOriginForm{
		User:      h.User,
		Password:  h.Password,
		NewOrigin: newOrigin,
	}

	_, err := h.doRequest(changeOriginPath, form, "", "POST", ChangeOrigin)
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
	_, err := h.doRequest(wipeDataPath, nil, "", "GET", WipeData)
	return err
}

func (h *HubClient) logout() error {
	_, err := h.doRequest(logoutPath, nil, "", "GET", Logout)
	return err
}
