package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// TODO Requires Auth
// TODO Only allowed when the target is the user itself. Cant upload stuff to other users.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logAndRespondError(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logAndRespondError(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// TODO Add test
	// TODO Make security test that user and repo are in the name correctly, and that both exist.
	if !strings.HasSuffix(header.Filename, ".tar.gz") {
		logAndRespondError(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	fileInfo, err := createFileInfo(header.Filename)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fileBuffer bytes.Buffer
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		logAndRespondError(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	err = CreateTag(fileInfo, &fileBuffer)
	if err != nil {
		logAndRespondError(w, "Failed to write content to local file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	Logger.Info("File uploaded successfully: %s", header.Filename)
}

// TODO Create unit tests
func createFileInfo(filename string) (*FileInfo, error) {
	infos := strings.Split(filename, "_")
	if len(infos) != 3 {
		return nil, fmt.Errorf("error, filenames should have exactly two underscores: %s", filename)
	}
	var info = &FileInfo{}
	info.User = infos[0]
	info.App = infos[1]
	// TODO consider error here and test it.
	info.FileName = infos[2]
	info.Tag, _ = strings.CutSuffix(infos[2], ".tar.gz")
	return info, nil
}

func logAndRespondError(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Error(msg)
	http.Error(w, msg, httpStatus)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	uploadName := strings.TrimPrefix(r.URL.Path, downloadPath)
	if uploadName == "" {
		logAndRespondError(w, "File name is missing", http.StatusBadRequest)
		return
	}

	fileInfo, err := createFileInfo(uploadName)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	path := fmt.Sprintf("%s/%s/%s/%s", usersDir, fileInfo.User, fileInfo.App, fileInfo.FileName)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		logAndRespondError(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, path)
}

type FileInfo struct {
	User     string
	App      string
	Tag      string
	FileName string
}

// TODO All functions below require auth
func userHandler(w http.ResponseWriter, r *http.Request) {
	// TODO create/delete user
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	// TODO create/delete app, search for app: search
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	// TODO delete tag, getListOfTags(app)
}
