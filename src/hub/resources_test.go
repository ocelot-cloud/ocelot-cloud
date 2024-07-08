package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
)

var (
	sampleUser                    = "myuser"
	sampleApp                     = "myapp"
	sampleTag                     = "v0.0.1"
	singleUserDir                 = usersDir + "/" + sampleUser
	appDir                        = singleUserDir + "/" + sampleApp
	sampleFile                    = appDir + fmt.Sprintf("/%s.tar.gz", sampleTag)
	sampleTagFileContent          = "hello"
	sampleTaggedFileContentBuffer = bytes.NewBuffer([]byte("hello"))
	sampleFileInfo                = &FileInfo{sampleUser, sampleApp, sampleTag, sampleFile}
	sampleMail                    = "testuser@example.com"
	samplePassword                = "mypassword"
	sampleOrigin                  = rootUrl
	sampleForm                    = &RegistrationForm{
		sampleUser,
		samplePassword,
		sampleMail,
		sampleOrigin,
	}
)

func cleanup() {
	err := deleteIfExist(dataDir)
	if err != nil {
		Logger.Error("Cleanup: Could not delete dir: %s", dataDir)
		os.Exit(1)
	}
}

type HubClient struct {
	User            string
	Password        string
	Origin          string
	Email           string
	App             string
	Cookie          *http.Cookie
	Tag             string
	TagFilename     string
	SetOriginHeader bool
	SetCookieHeader bool
}

func (h *HubClient) getRegistrationForm() *RegistrationForm {
	return &RegistrationForm{
		Username: h.User,
		Password: h.Password,
		Origin:   h.Origin,
		Email:    h.Email,
	}
}

func getHub() *HubClient {
	return &HubClient{
		User:            sampleUser,
		Password:        samplePassword,
		Origin:          rootUrl,
		Email:           sampleMail,
		App:             sampleApp,
		Tag:             sampleTag,
		TagFilename:     strings.Join([]string{sampleUser, sampleApp, sampleTag}, "_") + ".tar.gz",
		SetOriginHeader: true,
		SetCookieHeader: true,
	}
}

func (h *HubClient) registerUser(form *RegistrationForm) error {
	_, err := h.doRequest(registrationPath, form, "User registered\n", http.StatusCreated, "POST", Register)
	return err
}

func (h *HubClient) login() (*http.Cookie, error) {
	creds := LoginCredentials{
		Username: h.User,
		Password: h.Password,
	}

	result, err := h.doRequest(loginPath, creds, "login successful\n", http.StatusOK, "GET", Login)
	if err != nil {
		return nil, err
	}

	resp, ok := result.(*http.Response)
	if !ok {
		return nil, fmt.Errorf("Failed to assert result to *http.Response")
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		return nil, fmt.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	h.Cookie = cookies[0]
	return cookies[0], nil
}

func (h *HubClient) deleteUser() error {
	_, err := h.doRequest(userPath, SingleString{h.User}, "User deleted", http.StatusOK, "DELETE", DeleteUser)
	return err
}

func (h *HubClient) doRequest(path string, payload interface{}, expectedMessage string, expectedStatusCode int, method string, operation Operation) (interface{}, error) {
	url := rootUrl + path

	/* TODO
	policy := securityPolicies.getPolicyFor(operation)
	if policy.IsCredentialsRequired && payload != nil {
		return nil, fmt.Errorf("Security policy uses credentials in json body, so you can't define an addition apyload.")
	}*/

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
	if h.SetOriginHeader {
		req.Header.Set("Origin", h.Origin)
	}
	if h.SetCookieHeader && h.Cookie != nil {
		req.AddCookie(h.Cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var bodyBuffer bytes.Buffer
	teeReader := io.TeeReader(resp.Body, &bodyBuffer)

	if resp.StatusCode != expectedStatusCode {
		// Attempt to read the response body for additional error information
		respBody, err := io.ReadAll(teeReader)
		if err != nil {
			return nil, fmt.Errorf("Expected status code %d, but got %d. Also failed to read response body: %v", expectedStatusCode, resp.StatusCode, err)
		}
		return nil, fmt.Errorf("Expected status code %d, but got %d. Response body: %s", expectedStatusCode, resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	if expectedMessage != "" && string(respBody) != expectedMessage {
		return nil, fmt.Errorf("Expected response message '%s', got '%s'", expectedMessage, string(respBody))
	}

	if operation == Login {
		return resp, nil
	} else if operation == FindApps || operation == GetTags {
		return respBody, nil
	} else {
		return nil, nil
	}
}

func (h *HubClient) createApp() error {
	_, err := h.doRequest(appPath, SingleString{h.App}, "app created\n", http.StatusCreated, "POST", CreateApp)
	return err
}

func (h *HubClient) findApps(searchTerm string) ([]AppInfo, error) {
	result, err := h.doRequest(appPath, SingleString{searchTerm}, "", http.StatusOK, "GET", FindApps)
	if err != nil {
		return nil, err
	}

	respBody, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("Failed to assert result to []byte")
	}

	var apps []AppInfo
	err = json.Unmarshal(respBody, &apps)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v\n", err)
	}

	return apps, nil
}

func (h *HubClient) uploadTag(content string) error {
	fileContent := []byte(content)
	fileBuffer := bytes.NewBuffer(fileContent)

	url := rootUrl + tagPath
	filename := h.TagFilename
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return Logger.LogAndReturnError("Error creating file header: %v\n", err)
	}

	if _, err := io.Copy(part, fileBuffer); err != nil {
		return Logger.LogAndReturnError("Error copying content: %v\n", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}
	return nil
}

func (h *HubClient) downloadApp() (string, error) {
	downloadURL := rootUrl + downloadPath + h.TagFilename
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", Logger.LogAndReturnError("failed to download file, status code: %d", resp.StatusCode)
	}

	downloadedContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(downloadedContent), nil
}

// TODO Resolve duplication
func (h *HubClient) getTags() ([]string, error) {
	usernameAndApp := &UsernameAndApp{
		Username: h.User,
		App:      h.App,
	}

	result, err := h.doRequest(tagPath, usernameAndApp, "", http.StatusOK, "GET", GetTags)
	if err != nil {
		return nil, err
	}

	respBody, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("Failed to assert result to []byte")
	}

	var tags []string
	err = json.Unmarshal(respBody, &tags)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v\n", err)
	}

	return tags, nil
}

func (h *HubClient) deleteTag() error {
	tagInfo := &TagInfo{
		User: h.User,
		App:  h.App,
		Tag:  h.Tag,
	}
	// TODO check expected message
	_, err := h.doRequest(tagPath, tagInfo, "", http.StatusOK, "DELETE", DeleteTag)
	return err
}

func (h *HubClient) deleteApp() error {
	appInfo := &AppInfo{
		User: h.User,
		App:  h.App,
	}
	_, err := h.doRequest(appPath, appInfo, "app deleted\n", http.StatusOK, "DELETE", DeleteApp)
	return err
}

func (h *HubClient) ChangePassword(newPassword string) error {
	form := ChangePasswordForm{
		User:        h.User,
		OldPassword: h.Password,
		NewPassword: newPassword,
	}

	_, err := h.doRequest(changePasswordPath, form, "password changed\n", http.StatusOK, "POST", ChangePassword)
	return err
}

func (h *HubClient) ChangeOrigin(newOrigin string) interface{} {
	form := ChangeOriginForm{
		User:      h.User,
		Password:  h.Password,
		NewOrigin: newOrigin,
	}

	_, err := h.doRequest(changeOriginPath, form, "origin changed\n", http.StatusOK, "POST", ChangeOrigin)
	return err
}

func getHubAndLogin(t *testing.T) *HubClient {
	hub := getHub()
	form := hub.getRegistrationForm()
	assert.Nil(t, hub.registerUser(form))

	cookie, err := hub.login()
	assert.Nil(t, err)
	hub.Cookie = cookie
	return hub
}
