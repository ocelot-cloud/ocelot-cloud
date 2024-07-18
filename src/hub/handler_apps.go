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
		Logger.Debug("error reading body: %v", err)
		http.Error(w, "failed reading body", http.StatusBadRequest)
	}
	app := singleString.Value

	if !validate(app, App) {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesAppExist(user, app) {
		logAndRespondDebug(w, "app does not exist", http.StatusNotFound)
		return
	}

	err = repo.DeleteApp(user, app)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' deleted app '%s'", user, app)
	w.WriteHeader(http.StatusOK)
}

func findApps(w http.ResponseWriter, r *http.Request) {
	appSearchTerm, err := readBodyAsSingleString(r, User)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	apps, err := repo.FindApps(appSearchTerm)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Info("conducted app search for '%s'", appSearchTerm)
	sendJsonResponse(w, apps)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	singleString, err := readBody[SingleString](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}
	app := singleString.Value

	if !validate(authenticatedUser, User) || !validate(app, App) {
		logAndRespondDebug(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(authenticatedUser) {
		logAndRespondDebug(w, "user does not exists", http.StatusNotFound)
		return
	}
	if repo.DoesAppExist(authenticatedUser, app) {
		logAndRespondDebug(w, "app already exists", http.StatusConflict)
		return
	}
	err = repo.CreateApp(authenticatedUser, app)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	Logger.Info("user '%s' created app '%s'", authenticatedUser, app)
}
