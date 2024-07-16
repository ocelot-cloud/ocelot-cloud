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
		logAndRespondError(w, "method not implemented", http.StatusMethodNotAllowed)
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
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}
	app := singleString.Value

	if !validate(app, App) {
		logAndRespondDebug(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesAppExist(user, app) {
		logAndRespondDebug(w, "app does not exist", http.StatusNotFound)
		return
	}

	// TODO Replace all "authenticatedUser" simply with "user"

	err = repo.DeleteApp(user, app)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logAndRespondDebug(w, "app deleted", http.StatusOK)
}

func findApps(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("finding apps")
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

	logAndRespondDebug(w, "app created", http.StatusOK)
}
