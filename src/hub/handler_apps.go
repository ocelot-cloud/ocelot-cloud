package main

import (
	"net/http"
)

func appHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		findApps(w, r)
	} else if r.Method == http.MethodPost {
		createApp(w, r)
	} else if r.Method == http.MethodDelete {
		handleDeleteApp(w, r)
	} else {
		logAndRespondWarn(w, "method not implemented", http.StatusMethodNotAllowed)
		return
	}
}

func handleDeleteApp(w http.ResponseWriter, r *http.Request) {
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	singleString, err := readBody[SingleString](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
	}
	app := singleString.Value

	if !validate(app, App) {
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

func findApps(w http.ResponseWriter, r *http.Request) {
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
	sendJsonResponse(w, apps)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	singleString, err := readBody[SingleString](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	app := singleString.Value

	if !validate(authenticatedUser, User) || !validate(app, App) {
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
