package main

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
)

const maxPayloadSize = 1024 * 1024 // = 1 MiB
const maxStorageSize = 10 * maxPayloadSize

func tagUploadHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	r.Body = http.MaxBytesReader(w, r.Body, maxPayloadSize)
	defer r.Body.Close()

	var tagUpload TagUpload
	err := json.NewDecoder(r.Body).Decode(&tagUpload)
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
		{tagUpload.App, AppType},
		{tagUpload.Tag, TagType},
	}
	if err := validateJobs(jobs); err != nil {
		Logger.Info("tag upload of user '%s' invalid: %v", user, err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	bytesUsed, err := repo.GetUsedSpaceInBytes(user)
	if err != nil {
		Logger.Info("user '%s' tried to upload tag '%s' to app '%s', but getting current storage size failed", user, tagUpload.Tag, tagUpload.App)
		http.Error(w, "reading currently used storage failed", http.StatusInternalServerError)
		return
	} else if bytesUsed+len(tagUpload.Content) > maxStorageSize {
		Logger.Info("user '%s' tried to upload tag '%s' to app '%s', but exceeded max storage size", user, tagUpload.Tag, tagUpload.App)
		asdf := bytesUsed * 100 / maxStorageSize
		msg := fmt.Sprintf("storage limit reached, you can't store more then 10MiB of tag content, currently used storage in bytes: %d/%d (%d percent)", bytesUsed, maxStorageSize, asdf)
		http.Error(w, msg, http.StatusRequestEntityTooLarge)
		return
	}

	appId, err := repo.GetAppId(user, tagUpload.App)
	if err != nil {
		// TODO
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesAppExist(appId) {
		Logger.Info("user '%s' tried to upload tag '%s', but the app '%s' does not exist", user, tagUpload.Tag, tagUpload.App)
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	_, err = repo.GetTagId(appId, tagUpload.Tag)
	if err == nil {
		// TODO
		http.Error(w, "tag already exists", http.StatusConflict)
		return
	}

	err = repo.CreateTag(appId, tagUpload.Tag, tagUpload.Content)
	if err != nil {
		Logger.Error("creating tag '%s' for user '%s' failed: %v", tagUpload.App, user, err)
		http.Error(w, "invalid input", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' uploaded a new tag to the app '%s' with the tag name '%s'", user, tagUpload.App, tagUpload.Tag)
	w.WriteHeader(http.StatusOK)
}

func tagDeleteHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	tagInfo, err := readBody[AppAndTag](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	appId, err := repo.GetAppId(user, tagInfo.App)
	if err != nil {
		// TODO
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	tagId, err := repo.GetTagId(appId, tagInfo.Tag)
	if err != nil {
		// TODO
		http.Error(w, "tag does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesTagExist(tagId) {
		Logger.Info("user '%s' tried to delete tag of app '%s' but tag '%s' does not exist", user, tagInfo.App, tagInfo.Tag)
		http.Error(w, "tag does not exist", http.StatusNotFound)
		return
	}

	err = repo.DeleteTag(tagId)
	if err != nil {
		Logger.Info("user '%s' tried to delete tag in app '%s' with tag name '%s' but it failed", user, tagInfo.App, tagInfo.Tag)
		http.Error(w, "invalid input", http.StatusInternalServerError)
		return
	}
	Logger.Info("user '%s' deleted in tag in app '%s' with tag name '%s'", user, tagInfo.App, tagInfo.Tag)
	http.Error(w, "tag deleted", http.StatusOK)
}

func getTagsHandler(w http.ResponseWriter, r *http.Request) {
	appId, err := readBodyAsSingleInteger(r)
	if err != nil {
		Logger.Info("TODO: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// TODO check existence and ownership?

	if !repo.DoesAppExist(appId) {
		Logger.Info("someone tried to list tags but app with ID '%d' does not exist", appId)
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	tagsList, err := repo.GetTagList(appId)
	if err != nil {
		Logger.Error("getting tag list failed for app with ID '%d'", appId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJsonResponse(w, tagsList)
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

	appId, err := repo.GetAppId(tagInfo.User, tagInfo.App)
	if err != nil {
		// TODO
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesAppExist(appId) {
		Logger.Info("somebody tried to download users '%s' app '%s' with tag '%s', but app does not exist", tagInfo.User, tagInfo.App, tagInfo.Tag)
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	tagId, err := repo.GetTagId(appId, tagInfo.Tag)
	if err != nil {
		// TODO
		http.Error(w, "tag does not exist", http.StatusNotFound)
		return
	}

	if !repo.DoesTagExist(tagId) {
		Logger.Info("somebody tried to download users '%s' app '%s' with tag '%s', but tag does not exist", tagInfo.User, tagInfo.App, tagInfo.Tag)
		http.Error(w, "tag does not exist", http.StatusNotFound)
		return
	}

	content, err := repo.GetTagContent(tagId)
	if err != nil {
		Logger.Error("getting tag content failed for user='%s', app='%s' and tag='%s': %v", tagInfo.User, tagInfo.App, tagInfo.Tag, err)
		http.Error(w, "error when accessing tag content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/gzip")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
