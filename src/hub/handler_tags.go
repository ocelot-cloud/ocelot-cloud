package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func tagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleUpload(w, r)
	} else if r.Method == http.MethodGet {
		handleTagList(w, r)
	} else if r.Method == http.MethodDelete {
		handleDeleteTag(w, r)
	} else {
		logAndRespondDebug(w, "method not implemented", http.StatusMethodNotAllowed)
	}
}

// TODO My impression is there might be some duplication with other data structure. To be checked for abstration.
type TagInfo struct {
	User string `json:"user"`
	App  string `json:"app"`
	Tag  string `json:"tag"`
}

func handleDeleteTag(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := middleware(w, r)
	if err != nil {
		return
	}

	tagInfo, err := readBody[TagInfo](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validate(tagInfo.User, User) || !validate(tagInfo.App, App) || !validate(tagInfo.Tag, Tag) {
		logAndRespondDebug(w, "invalid input", http.StatusBadRequest)
		return
	}

	if authenticatedUser != tagInfo.User {
		logAndRespondDebug(w, "deleting tags not belonging to you is not allowed", http.StatusUnauthorized)
		return
	}

	fs.DeleteTag(tagInfo.User, tagInfo.App, tagInfo.Tag) // TODO make it return an error.
	err = repo.DeleteTag(tagInfo.User, tagInfo.App, tagInfo.Tag)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logAndRespondDebug(w, "tag deleted", http.StatusOK)
}

type UserAndApp struct {
	User string `json:"username"`
	App  string `json:"app"`
}

// TODO Implement
func handleTagList(w http.ResponseWriter, r *http.Request) {
	usernameAndApp, err := readBody[UserAndApp](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}

	tagsList, err := repo.GetTagList(usernameAndApp.User, usernameAndApp.App)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, tagsList)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := middleware(w, r)
	if err != nil {
		return
	}

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

	if !validate(fileInfo.User, User) || !validate(fileInfo.App, App) || !validate(fileInfo.Tag, Tag) {
		logAndRespondDebug(w, "invalid input", http.StatusBadRequest)
		return
	}

	if authenticatedUser != fileInfo.User {
		logAndRespondDebug(w, "upload of tags not belonging to you is not allowed", http.StatusUnauthorized)
		return
	}

	var fileBuffer bytes.Buffer
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		logAndRespondError(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	err = fs.CreateTag(fileInfo, &fileBuffer)
	if err != nil {
		logAndRespondError(w, "Failed to write content to local file", http.StatusInternalServerError)
		return
	}

	err = repo.CreateTag(fileInfo.User, fileInfo.App, fileInfo.Tag)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logAndRespondDebug(w, "file uploaded successfully", http.StatusOK)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	uploadName := strings.TrimPrefix(r.URL.Path, downloadPath)
	if uploadName == "" {
		logAndRespondError(w, "File name is missing", http.StatusBadRequest)
		return
	}

	if !validate(uploadName, TagFile) {
		logAndRespondError(w, "invalid input", http.StatusBadRequest)
		return
	}

	fileInfo, err := createFileInfo(uploadName)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("%s/%s/%s/%s", usersDir, fileInfo.User, fileInfo.App, fileInfo.Tag+".tar.gz")
	if _, err = os.Stat(path); os.IsNotExist(err) {
		logAndRespondError(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, path)
}

func getTagFileName(user string, app string, tag string) string {
	return strings.Join([]string{user, app, tag}, "_") + ".tar.gz"
}
