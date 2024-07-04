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

type Hub struct{}

func TestCreateUser(t *testing.T) {
	hub := Hub{}
	form := &RegistrationForm{
		Username: "testuser",
		Password: "password123",
		Host:     "http://localhost:8082",
		Email:    "testuser@example.com",
	}
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

// TODO Add a cleanup function after that.
// TODO initialize hub with form, maybe in a setup function?
func TestDeleteUser(t *testing.T) {
	hub := Hub{}
	form := &RegistrationForm{
		Username: "testuser2",
		Password: "password123",
		Host:     "http://localhost:8082",
		Email:    "testuser@example.com",
	}
	assert.Nil(t, hub.createUser(form))
	assert.Nil(t, hub.deleteUser())
}

// TODO Can just be done, when I have a protected endpoint
func TestOriginPolicy(t *testing.T) {
	hub := Hub{}
	form := &RegistrationForm{
		Username: "testuser3",
		Password: "password123",
		Host:     "http://non-existing-domain:8082",
		Email:    "testuser@example.com",
	}
	hub.createUser(form)
}

func (h *Hub) createUser(form *RegistrationForm) error {
	_, err := h.doPostRequest("/users", form, "User created", http.StatusCreated)
	return err
}

func (h *Hub) login() (*http.Cookie, error) {
	url := rootUrl + "/login"
	user := LoginCredentials{
		Username: "testuser",
		Password: "password123",
	}
	payloadBytes, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal credentials: %v", err)
	}
	payload := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", url, payload)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	expectedResponse := "login successful"
	if string(respBody) != expectedResponse {
		return nil, fmt.Errorf("Expected response body %s, got %s", expectedResponse, string(respBody))
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		return nil, fmt.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	return cookies[0], nil // TODO return cookie for assertion
}

func (h *Hub) deleteUser() error {
	url := rootUrl + "/users"
	user := User{
		Name: "testuser2",
	}
	payloadBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %v", err)
	}
	payload := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("DELETE", url, payload)
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %v", err)
	}

	expectedResponse := "User deleted"
	if string(respBody) != expectedResponse {
		return fmt.Errorf("Expected response body %s, got %s", expectedResponse, string(respBody))
	}
	return nil
}

func (h *Hub) doPostRequest(path string, payload interface{}, expectedMessage string, expectedStatusCode int) (*http.Response, error) {
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
