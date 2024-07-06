//go:build acceptance

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
	"time"
)

// TODO There should be a one-liner to get a "hub" instance, that is already logged in, with cookie etc, to start functional testing
// TODO this should be generated from anew for each test.

func TestFileUploadDownload(t *testing.T) {
	defer cleanup()
	// TODO Should be global?
	fs := FileStorageImpl{}
	assert.Nil(t, fs.CreateUser(sampleUser))
	assert.Nil(t, fs.CreateApp(sampleUser, sampleApp))
	filename := "myuser_myapp_v0.1.0.tar.gz"

	fileContent := []byte("hello")
	fileBuffer := bytes.NewBuffer(fileContent)

	uploadURL := rootUrl + uploadPath
	responseCode, err := uploadFile(uploadURL, filename, fileBuffer)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, responseCode)

	downloadURL := rootUrl + downloadPath + filename
	downloadedContent, err := downloadFile(downloadURL)
	assert.Nil(t, err)
	assert.Equal(t, fileContent, downloadedContent)
}

func uploadFile(url, filename string, fileBuffer *bytes.Buffer) (int, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return 0, Logger.LogAndReturnError("Error creating file header: %v\n", err)
	}

	if _, err := io.Copy(part, fileBuffer); err != nil {
		return 0, Logger.LogAndReturnError("Error copying content: %v\n", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, Logger.LogAndReturnError("failed to download file, status code: %d", resp.StatusCode)
	}

	downloadedContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return downloadedContent, nil
}

// TODO High level test: create myuser, create myapp, findApps -> One element {myuser, myapp}

type HubClient struct {
	Username        string
	Password        string
	Origin          string
	Email           string
	SetOriginHeader bool
	App             string
	Cookie          *http.Cookie
}

func getRegistrationForm(hub *HubClient) *RegistrationForm {
	return &RegistrationForm{
		Username: hub.Username,
		Password: hub.Password,
		Origin:   hub.Origin,
		Email:    hub.Email,
	}
}

func getHub() *HubClient {
	return &HubClient{
		Username:        "testuser",
		Password:        "password123",
		Origin:          "http://localhost:8082",
		Email:           "testuser@example.com",
		SetOriginHeader: true,
		App:             sampleApp,
	}
}

// TODO Test if cookie expiration date updates when making a successful request.

func TestCookie(t *testing.T) {
	hub := getHub()
	defer assert.Nil(t, hub.deleteUser())
	form := getRegistrationForm(hub)
	assert.Nil(t, hub.registerUser(form))
	cookie, err := hub.login()
	assert.Nil(t, err)
	assert.NotNil(t, cookie)
	assert.Equal(t, cookieName, cookie.Name)
	assert.True(t, getTimeIn30Days().Add(1*time.Second).After(cookie.Expires))
	assert.True(t, getTimeIn30Days().Add(-1*time.Second).Before(cookie.Expires))
	assert.Equal(t, 64, len(cookie.Value))

	cookie2, err := hub.login()
	assert.Nil(t, err)
	assert.NotNil(t, cookie2)
	assert.NotEqual(t, cookie.Value, cookie2.Value)
}

func TestCreateApp(t *testing.T) {
	hub := getHubAndLogin(t)
	assert.Nil(t, hub.createApp())
	foundApps, err := hub.findApps(sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundApps))
	app := foundApps[0]
	assert.Equal(t, hub.Username, app.Username)
	assert.Equal(t, hub.App, app.AppName)
}

func getHubAndLogin(t *testing.T) *HubClient {
	hub := getHub()
	defer assert.Nil(t, hub.deleteUser())
	form := getRegistrationForm(hub)
	assert.Nil(t, hub.registerUser(form))

	cookie, err := hub.login()
	assert.Nil(t, err)
	hub.Cookie = cookie
	return hub
}

// TODO After running "ci-runner" hub tests, I still have a "testuser" folder in data. Should actually be deleted after the test creating it.

// TODO Can just be done, when I have a protected endpoint
func TestOriginPolicy(t *testing.T) {
	hub := getHub()
	form := getRegistrationForm(hub)
	fakeOrigin := "http://non-existing-subdomain.localhost:8082"
	assert.Nil(t, hub.registerUser(form))

	// TODO
	/*hub.SetOriginHeader = false
	err := hub.deleteUser()
	assert.NotNil(t, err)
	expected := fmt.Sprintf("Security policy does not allow this request without 'Origin' header")
	assert.Equal(t, expected, err.Error())
	*/

	form.Origin = fakeOrigin
	// TODO expected := fmt.Sprintf("Security policy does not allow requests from origin: %s", fakeOrigin)
}

func (h *HubClient) registerUser(form *RegistrationForm) error {
	_, err := h.doRequest(registrationPath, form, "User registered", http.StatusCreated, "POST", Register)
	return err
}

func (h *HubClient) login() (*http.Cookie, error) {
	creds := LoginCredentials{
		Username: h.Username,
		Password: h.Password,
	}

	result, err := h.doRequest(loginPath, creds, "login successful", http.StatusOK, "GET", Login)
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
	return cookies[0], nil // TODO return cookie for assertion
}

func (h *HubClient) deleteUser() error {
	_, err := h.doRequest(userPath, SingleString{h.Username}, "User deleted", http.StatusOK, "DELETE", DeleteUser)
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

	if h.Cookie != nil {
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
	} else if operation == FindApps {
		return respBody, nil
	} else {
		return nil, nil
	}

}

func (h *HubClient) createApp() error {
	_, err := h.doRequest(appPath, SingleString{h.App}, "created app successfully", http.StatusCreated, "POST", CreateApp)
	return err
}

func (h *HubClient) findApps(searchTerm string) ([]App, error) {
	result, err := h.doRequest(appPath, SingleString{searchTerm}, "", http.StatusOK, "GET", FindApps)
	if err != nil {
		return nil, err
	}

	respBody, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("Failed to assert result to []byte")
	}

	var apps []App
	err = json.Unmarshal(respBody, &apps)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v\n", err)
	}

	return apps, nil
}

// TODO assert that no other object should be send in body, should be nil, when IsCredentialsRequired == true
