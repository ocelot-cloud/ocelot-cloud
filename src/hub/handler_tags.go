package main

import (
	"encoding/json"
	"net/http"
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

type AppAndTag struct {
	App string `json:"app"`
	Tag string `json:"tag"`
}

func handleDeleteTag(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	tagInfo, err := readBody[AppAndTag](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validate(tagInfo.App, App) || !validate(tagInfo.Tag, Tag) {
		logAndRespondDebug(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesTagExist(authenticatedUser, tagInfo.App, tagInfo.Tag) {
		logAndRespondDebug(w, "tag does not exist", http.StatusNotFound)
		return
	}

	err = repo.DeleteTag(authenticatedUser, tagInfo.App, tagInfo.Tag)
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

func handleTagList(w http.ResponseWriter, r *http.Request) {
	usernameAndApp, err := readBody[UserAndApp](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(usernameAndApp.User) {
		logAndRespondDebug(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesAppExist(usernameAndApp.User, usernameAndApp.App) {
		logAndRespondDebug(w, "app does not exist", http.StatusNotFound)
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
	authenticatedUser, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	if r.Method != http.MethodPost {
		logAndRespondError(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// TODO I think this should be done when the other formalities like app/tag extraction/validation are done.
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

	// TODO Should take fileInfo structure as arg
	err = repo.CreateTag(authenticatedUser, tagUpload.App, tagUpload.Tag, tagUpload.Content)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logAndRespondDebug(w, "file uploaded successfully", http.StatusOK)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileInfo, err := readBody[TagInfo](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(fileInfo.User) {
		logAndRespondDebug(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesAppExist(fileInfo.User, fileInfo.App) {
		logAndRespondDebug(w, "app does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesTagExist(fileInfo.User, fileInfo.App, fileInfo.Tag) {
		logAndRespondDebug(w, "tag does not exist", http.StatusNotFound)
		return
	}

	content, err := repo.GetTagContent(fileInfo.User, fileInfo.App, fileInfo.Tag)
	if err != nil {
		logAndRespondError(w, "error when accessing tag content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/gzip")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
