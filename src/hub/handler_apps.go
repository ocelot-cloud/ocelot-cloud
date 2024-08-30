package main

import (
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
)

func appHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	app, err := readBodyAsSingleString(r, App)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(authenticatedUser) {
		Logger.Info("user '%s' tried to create app '%s' but it does not exist", authenticatedUser, app)
		http.Error(w, "user does not exists", http.StatusNotFound)
		return
	}
	if repo.DoesAppExist(authenticatedUser, app) {
		Logger.Info("user '%s' tried to create app '%s' but it already exists", authenticatedUser, app)
		http.Error(w, "app already exists", http.StatusConflict)
		return
	}

	err = repo.CreateApp(authenticatedUser, app)
	if err != nil {
		Logger.Error("user '%s' tried to create app '%s' but it failed: %v", authenticatedUser, app, err)
		http.Error(w, "app creation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	Logger.Info("user '%s' created app '%s'", authenticatedUser, app)
}

func appDeleteHandler(w http.ResponseWriter, r *http.Request) {
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	app, err := readBodyAsSingleString(r, App)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesAppExist(user, app) {
		Logger.Info("user '%s' tried to delete app '%s' but it does not exist", user, app)
		http.Error(w, "app does not exist", http.StatusNotFound)
		return
	}

	err = repo.DeleteApp(user, app)
	if err != nil {
		Logger.Error("user '%s' tried to delete app '%s' but it failed", user, app)
		http.Error(w, "app deletion failed", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' deleted app '%s'", user, app)
	w.WriteHeader(http.StatusOK)
}

func appGetListHandler(w http.ResponseWriter, r *http.Request) {
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	list, err := repo.GetAppList(user)
	if err != nil {
		Logger.Warn("error getting app list: %v", err)
		http.Error(w, "error getting app list", http.StatusInternalServerError)
	}

	Logger.Info("got apps of user '%s'", user)
	utils.SendJsonResponse(w, list)
}

func searchAppsHandler(w http.ResponseWriter, r *http.Request) {
	appSearchTerm, err := readBodyAsSingleString(r, User)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	apps, err := repo.FindApps(appSearchTerm)
	if err != nil {
		Logger.Warn("error finding apps: %v", err)
		http.Error(w, "error finding apps", http.StatusInternalServerError)
		return
	}

	Logger.Info("conducted app search with search term '%s'", appSearchTerm)
	utils.SendJsonResponse(w, apps)
}
