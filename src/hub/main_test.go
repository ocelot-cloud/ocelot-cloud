package main

import (
	"bytes"
	"github.com/ocelot-cloud/shared"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
)

// TODO Clean up "users" folder before and after test. Also store in "data/users" instead.
func TestFileUploadDownload(t *testing.T) {
	defer cleanup()
	shared.AssertNil(t, CreateUser(sampleUser))
	shared.AssertNil(t, CreateApp(sampleUser, sampleApp))
	filename := "myuser_myapp_v0.1.0.tar.gz"

	fileContent := []byte("hello")
	fileBuffer := bytes.NewBuffer(fileContent)

	uploadURL := rootUrl + uploadPath
	responseCode, err := uploadFile(uploadURL, filename, fileBuffer)
	shared.AssertNil(t, err)
	shared.AssertEqual(t, http.StatusOK, responseCode)

	downloadURL := rootUrl + downloadPath + filename
	downloadedContent, err := downloadFile(downloadURL)
	shared.AssertNil(t, err)
	shared.AssertEqual(t, fileContent, downloadedContent)
}

func uploadFile(url, filename string, fileBuffer *bytes.Buffer) (int, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return 0, logger.LogAndReturnError("Error creating file header: %v\n", err)
	}

	if _, err := io.Copy(part, fileBuffer); err != nil {
		return 0, logger.LogAndReturnError("Error copying content: %v\n", err)
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
		return nil, logger.LogAndReturnError("failed to download file, status code: %d", resp.StatusCode)
	}

	downloadedContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return downloadedContent, nil
}
