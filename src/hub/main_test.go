package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestFileUploadDownload(t *testing.T) {
	filename := "myuser_myapp_v1.0.tar.gz"
	usersDir := "users"
	uploadedFilePath := filepath.Join(usersDir, filename)

	if _, err := os.Stat(usersDir); os.IsNotExist(err) {
		t.Fatalf("Users directory does not exist: %v", err)
	}

	fileContent := []byte("hello")
	fileBuffer := bytes.NewBuffer(fileContent)

	uploadURL := "http://localhost:8082/upload"
	responseCode, err := uploadFile(uploadURL, filename, fileBuffer)
	if err != nil {
		t.Fatalf("Failed to upload file: %v", err)
	}
	if responseCode != http.StatusOK {
		t.Fatalf("Failed to upload file. HTTP status code: %d", responseCode)
	}

	if _, err := os.Stat(uploadedFilePath); os.IsNotExist(err) {
		t.Fatalf("Uploaded file does not exist in the users directory")
	}

	downloadURL := "http://localhost:8082/download/" + filename
	downloadedContent, err := downloadFile(downloadURL)
	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	if !bytes.Equal(fileContent, downloadedContent) {
		t.Fatalf("Downloaded content does not match uploaded content")
	}
}

func uploadFile(url, filename string, fileBuffer *bytes.Buffer) (int, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return 0, err
	}

	if _, err := io.Copy(part, fileBuffer); err != nil {
		return 0, err
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
		return nil, fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	downloadedContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return downloadedContent, nil
}
