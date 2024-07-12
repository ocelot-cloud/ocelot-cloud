package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type TagUpload struct {
	App     string `json:"app"`
	Tag     string `json:"tag"`
	Content []byte `json:"content"`
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

	var tagUpload TagUpload
	err = json.NewDecoder(r.Body).Decode(&tagUpload)
	if err != nil {
		logAndRespondError(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	if !validate(tagUpload.App, App) || !validate(tagUpload.Tag, Tag) {
		logAndRespondDebug(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesAppExist(authenticatedUser, tagUpload.App) {
		logAndRespondDebug(w, "app does not exist", http.StatusNotFound)
		return
	}

	if repo.DoesTagExist(authenticatedUser, tagUpload.App, tagUpload.Tag) {
		logAndRespondDebug(w, "tag already exists", http.StatusConflict)
		return
	}

	fileBuffer := bytes.NewBuffer(tagUpload.Content)

	fileInfo := &FileInfo{authenticatedUser, tagUpload.App, tagUpload.Tag}
	err = fs.CreateTag(fileInfo, fileBuffer)
	if err != nil {
		logAndRespondError(w, "Failed to write content to local file", http.StatusInternalServerError)
		return
	}

	err = repo.CreateTag(authenticatedUser, tagUpload.App, tagUpload.Tag)
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

	// TODO I think this should be handled via json, not a upload path.
	fileInfo, err := createFileDownloadInfo(uploadName)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO Validate and re-enable test

	path := fmt.Sprintf("%s/%s/%s/%s", usersDir, fileInfo.User, fileInfo.App, fileInfo.Tag+".tar.gz")
	if _, err = os.Stat(path); os.IsNotExist(err) {
		logAndRespondError(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, path)
}

// TODO Should be removed when using jsons instead of paths.
func getDownloadFileName(user string, app string, tag string) string {
	return user + "_" + app + "_" + tag + ".tar.gz"
}

func getUploadFileName(app string, tag string) string {
	return app + "_" + tag + ".tar.gz"
}
