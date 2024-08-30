package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/secutils"
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
	}
)

type ComponentClient struct {
	User            string
	Password        string
	NewPassword     string
	Origin          string
	Cookie          *http.Cookie
	SetOriginHeader bool
	SetCookieHeader bool
}

type HubClient struct {
	Parent        ComponentClient
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
		Parent: ComponentClient{
			User:            sampleUser,
			Password:        samplePassword,
			Origin:          rootUrl,
			SetOriginHeader: true,
			SetCookieHeader: true,
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
	_, err := h.Parent.doRequest(registrationPath, form, "")
	return err
}

func (h *HubClient) login() error {
	creds := LoginCredentials{
		User:     h.Parent.User,
		Password: h.Parent.Password,
		Origin:   h.Parent.Origin,
	}

	resp, err := h.Parent.doRequestWithFullResponse(loginPath, creds, "")
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
	_, err := h.Parent.doRequest(deleteUserPath, nil, "")
	return err
}

func (c *ComponentClient) doRequest(path string, payload interface{}, expectedMessage string) (interface{}, error) {
	resp, err := c.doRequestWithFullResponse(path, payload, expectedMessage)
	if err != nil {
		return nil, err
	}

	respBody, err := assertOkStatusAndExtractBody(resp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *ComponentClient) doRequestWithFullResponse(path string, payload interface{}, expectedMessage string) (*http.Response, error) {
	url := rootUrl + path

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal payload: %v", err)
	}
	payloadReader := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", url, payloadReader)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	setCookieAndOriginHeaders(req, c)

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
		c.Cookie = resp.Cookies()[0]
	}

	// Response body can only be read once. When reading it after this function, an error occurs. So a copy is created.
	newResp := &http.Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       io.NopCloser(bytes.NewBuffer(respBody)),
	}
	return newResp, nil
}

func setCookieAndOriginHeaders(req *http.Request, c *ComponentClient) {
	if c.SetOriginHeader {
		req.Header.Set(secutils.OriginHeader, c.Origin)
	}
	if c.SetCookieHeader && c.Cookie != nil {
		req.AddCookie(c.Cookie)
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
	_, err := h.Parent.doRequest(appCreationPath, secutils.SingleString{h.App}, "")
	return err
}

func (h *HubClient) findApps(searchTerm string) ([]UserAndApp, error) {
	result, err := h.Parent.doRequest(searchAppsPath, secutils.SingleString{searchTerm}, "")
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
	result, err := h.Parent.doRequest(appGetListPath, nil, "")
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
	_, err := h.Parent.doRequest(tagUploadPath, tapUpload, "")
	return err
}

func (h *HubClient) downloadTag() (string, error) {
	tagInfo := &TagInfo{
		User: h.Parent.User,
		App:  h.App,
		Tag:  h.Tag,
	}

	result, err := h.Parent.doRequest(downloadPath, tagInfo, "")
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

	result, err := h.Parent.doRequest(getTagsPath, usernameAndApp, "")
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
	_, err := h.Parent.doRequest(tagDeletePath, tagInfo, "")
	return err
}

func (h *HubClient) deleteApp() error {
	_, err := h.Parent.doRequest(appDeletePath, secutils.SingleString{h.App}, "")
	return err
}

func (h *HubClient) changePassword() error {
	form := secutils.ChangePasswordForm{
		OldPassword: h.Parent.Password,
		NewPassword: h.Parent.NewPassword,
	}

	_, err := h.Parent.doRequest(changePasswordPath, form, "")
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
	_, err := h.Parent.doRequest(wipeDataPath, nil, "")
	return err
}

func (h *HubClient) logout() error {
	_, err := h.Parent.doRequest(logoutPath, nil, "")
	return err
}

func (h *HubClient) checkAuth() error {
	_, err := h.Parent.doRequest(authCheckPath, nil, "")
	return err
}
