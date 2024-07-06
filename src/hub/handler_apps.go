package main

import (
	"encoding/json"
	"net/http"
)

func appHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		findApps(w, r)
	} else if r.Method == http.MethodPost {
		createApp(w, r)
	} else {
		logAndRespondError(w, "method not implemented", http.StatusMethodNotAllowed)
		return
	}
}

func findApps(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("finding apps")
	singleString, err := readBody[SingleString](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	searchTerm := singleString.Value

	apps, err := repo.FindApps(searchTerm)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(apps)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie == nil || cookie.Value == "" {
		logAndRespondDebug(w, "Cookie not contained in request", http.StatusBadRequest)
		return
	}
	user, err := repo.GetUserWithCookie(cookie.Value)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusInternalServerError)
		return
	}

	singleString, err := readBody[SingleString](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}
	app := singleString.Value

	if !repo.DoesUserExist(user) {
		logAndRespondDebug(w, "user does not exists", http.StatusNotFound)
		return
	}
	if repo.DoesAppExist(user, app) {
		logAndRespondDebug(w, "app already exists", http.StatusConflict)
		return
	}
	err = fs.CreateApp(user, app)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = repo.CreateApp(user, app)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logAndRespondDebug(w, "app created", http.StatusCreated)
}
