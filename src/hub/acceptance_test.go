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

var hub = Hub{
	Username: "testuser",
	Password: "password123",
	Host:     "http://localhost:8082",
	Email:    "testuser@example.com",
}

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

type Hub struct {
	Username string
	Password string
	Host     string
	Email    string
}

func getRegistrationForm() *RegistrationForm {
	return &RegistrationForm{
		Username: hub.Username,
		Password: hub.Password,
		Host:     hub.Host,
		Email:    hub.Email,
	}
}

func TestCreateUser(t *testing.T) {
	defer assert.Nil(t, hub.deleteUser())
	form := getRegistrationForm()
	assert.Nil(t, hub.createUser(form))
	cookie, err := hub.login()
	assert.Nil(t, err)
	assert.NotNil(t, cookie)
	assert.Equal(t, "auth", cookie.Name)
	assert.True(t, getTimeIn30Days().Add(1*time.Second).After(cookie.Expires))
	assert.True(t, getTimeIn30Days().Add(-1*time.Second).Before(cookie.Expires))
	assert.Equal(t, 64, len(cookie.Value))

	cookie2, err := hub.login()
	assert.Nil(t, err)
	assert.NotNil(t, cookie2)
	assert.NotEqual(t, cookie.Value, cookie2.Value)
}

// TODO Can just be done, when I have a protected endpoint
func TestOriginPolicy(t *testing.T) {
	form := getRegistrationForm()
	form.Host = "http://non-existing-subdomain.localhost:8082"
	hub.createUser(form)
	//assert.NotNil(t, err)
	//assert.Equal(t, "Origin policy ", err.Error())
}

func (h *Hub) createUser(form *RegistrationForm) error {
	_, err := h.doRequest("/users", form, "User created", http.StatusCreated, "POST")
	return err
}

func (h *Hub) login() (*http.Cookie, error) {
	creds := LoginCredentials{
		Username: "testuser",
		Password: "password123",
	}
	resp, err := h.doRequest("/login", creds, "login successful", http.StatusOK, "GET")
	if err != nil {
		return nil, err
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		return nil, fmt.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	return cookies[0], nil // TODO return cookie for assertion
}

func (h *Hub) deleteUser() error {
	_, err := h.doRequest("/users", User{"testuser2"}, "User deleted", http.StatusOK, "DELETE")
	return err
}

func (h *Hub) doRequest(path string, payload interface{}, expectedMessage string, expectedStatusCode int, method string) (*http.Response, error) {
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		return nil, fmt.Errorf("Expected status code %d, but got %d", expectedStatusCode, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) != expectedMessage {
		return nil, fmt.Errorf("Expected response body %s, got %s", expectedMessage, string(respBody))
	}
	return resp, nil
}
