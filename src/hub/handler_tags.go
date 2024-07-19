package main

import (
	"encoding/json"
	"fmt"
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
		Logger.Warn("incoming request for method '%s' to endpoint '%s' which is not allowed", r.Method, tagPath)
		http.Error(w, "method not implemented", http.StatusMethodNotAllowed)
		return
	}
}

func handleDeleteTag(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	tagInfo, err := readBody[AppAndTag](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesTagExist(authenticatedUser, tagInfo.App, tagInfo.Tag) {
		Logger.Info("user '%s' tried to delete tag of app '%s' but tag '%s' does not exist", authenticatedUser, tagInfo.App, tagInfo.Tag)
		http.Error(w, "tag does not exist", http.StatusNotFound)
		return
	}

	err = repo.DeleteTag(authenticatedUser, tagInfo.App, tagInfo.Tag)
	if err != nil {
		Logger.Info("user '%s' tried to delete tag in app '%s' with tag name '%s' but it failed", authenticatedUser, tagInfo.App, tagInfo.Tag)
		http.Error(w, "invalid input", http.StatusInternalServerError)
		return
	}
	Logger.Info("user '%s' deleted in tag in app '%s' with tag name '%s'", authenticatedUser, tagInfo.App, tagInfo.Tag)
	http.Error(w, "tag deleted", http.StatusOK)
}

func handleTagList(w http.ResponseWriter, r *http.Request) {
	userAndApp, err := readBody[UserAndApp](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(userAndApp.User) {
		Logger.Info("someone tried to list tags but user '%s' does not exist", userAndApp.User)
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesAppExist(userAndApp.User, userAndApp.App) {
		Logger.Info("someone tried to list tags but app '%s' does not exist", userAndApp.App)
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	tagsList, err := repo.GetTagList(userAndApp.User, userAndApp.App)
	if err != nil {
		Logger.Error("getting tag list failed for user '%s' and app '%s'", userAndApp.User, userAndApp.App)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, tagsList)
}

const maxPayloadSize = 1024 * 1024 // = 1 MiB
const maxStorageSize = 10 * maxPayloadSize

func handleUpload(w http.ResponseWriter, r *http.Request) {
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	if r.Method != http.MethodPost {
		handleInvalidRequestMethod(w, r, tagPath)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPayloadSize)
	defer r.Body.Close()

	var tagUpload TagUpload
	err = json.NewDecoder(r.Body).Decode(&tagUpload)
	if err != nil {
		if err.Error() == "http: request body too large" {
			Logger.Info("tag upload tag content of user '%s' was too large", user)
			http.Error(w, "tag content too large, the limit is 1MB", http.StatusRequestEntityTooLarge)
			return
		} else {
			Logger.Info("tag upload request body of user '%s' was invalid: %v", user, err)
			http.Error(w, "could not decode request body", http.StatusBadRequest)
			return
		}
	}

	jobs := []ValidationJob{
		{tagUpload.App, App},
		{tagUpload.Tag, Tag},
	}
	if err := validateJobs(jobs); err != nil {
		Logger.Info("tag upload of user '%s' invalid: %v", user, err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	bytesUsed, err := repo.GetUsedSpaceInBytes(user)
	if err != nil {
		Logger.Info("TODO", user, tagUpload.Tag, tagUpload.App)
		http.Error(w, "TODO", http.StatusNotFound)
		return
	} else if bytesUsed+len(tagUpload.Content) > maxStorageSize {
		Logger.Info("user '%s' tried to upload tag '%s', but exceeded max storage size", user, tagUpload.Tag, tagUpload.App)
		asdf := bytesUsed * 100 / maxStorageSize
		msg := fmt.Sprintf("storage limit reached, you can't store more then 10MiB of tag content, currently used storage in bytes: %d/%d (%d percent)", bytesUsed, maxStorageSize, asdf)
		http.Error(w, msg, http.StatusRequestEntityTooLarge)
		return
	}

	if !repo.DoesAppExist(user, tagUpload.App) {
		Logger.Info("user '%s' tried to upload tag '%s', but the app '%s' does not exist", user, tagUpload.Tag, tagUpload.App)
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	if repo.DoesTagExist(user, tagUpload.App, tagUpload.Tag) {
		Logger.Info("user '%s' tried to upload a new tag to app '%s' with tag '%s', but the tag already exists", user, tagUpload.App, tagUpload.Tag)
		http.Error(w, "tag already exists", http.StatusConflict)
		return
	}

	err = repo.CreateTag(user, tagUpload.App, tagUpload.Tag, tagUpload.Content)
	if err != nil {
		Logger.Error("creating tag '%s' for user '%s' failed: %v", tagUpload.App, user, err)
		http.Error(w, "invalid input", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' uploaded a new tag to the app '%s' with the tag name '%s'", user, tagUpload.App, tagUpload.Tag)
	w.WriteHeader(http.StatusOK)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	tagInfo, err := readBody[TagInfo](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(tagInfo.User) {
		Logger.Info("somebody tried to download users '%s' app '%s' with tag '%s', but user does not exist", tagInfo.User, tagInfo.App, tagInfo.Tag)
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesAppExist(tagInfo.User, tagInfo.App) {
		Logger.Info("somebody tried to download users '%s' app '%s' with tag '%s', but app does not exist", tagInfo.User, tagInfo.App, tagInfo.Tag)
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesTagExist(tagInfo.User, tagInfo.App, tagInfo.Tag) {
		Logger.Info("somebody tried to download users '%s' app '%s' with tag '%s', but tag does not exist", tagInfo.User, tagInfo.App, tagInfo.Tag)
		http.Error(w, "tag does not exist", http.StatusNotFound)
		return
	}

	content, err := repo.GetTagContent(tagInfo.User, tagInfo.App, tagInfo.Tag)
	if err != nil {
		Logger.Error("getting tag content failed for user='%s', app='%s' and tag='%s': %v", tagInfo.User, tagInfo.App, tagInfo.Tag, err)
		http.Error(w, "error when accessing tag content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/gzip")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
