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

func TestCreateUserViaHttp(t *testing.T) {
	err := createUser()
	assert.Nil(t, err)
}

func createUser() error {
	url := rootUrl + "/users"
	user := User{
		Username: "testuser",
		Password: "password123",
		Host:     "localhost",
	}
	payloadBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to marshal user: %v", err)
	}
	payload := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", url, payload)
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

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %v", err)
	}

	expectedResponse := "User created"
	if string(respBody) != expectedResponse {
		return fmt.Errorf("Expected response body %s, got %s", expectedResponse, string(respBody))
	}
	return nil
}
