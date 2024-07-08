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
	appInfo, err := readBody[AppInfo](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
	}

	err = fs.DeleteApp(appInfo.User, appInfo.App)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = repo.DeleteApp(appInfo.User, appInfo.App)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logAndRespondDebug(w, "app deleted", http.StatusOK)
}

func findApps(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("finding apps")
	singleString, err := readBody[SingleString](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	appSearchTerm := singleString.Value
	if !isValidName(appSearchTerm) {
		logAndRespondDebug(w, "invalid app search term", http.StatusBadRequest)
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

	// TODO Should be used by security policy logic, maybe?
	if !repo.IsOriginCorrect(user, r.Header.Get("Origin")) {
		logAndRespondDebug(w, "wrong origin", http.StatusUnauthorized)
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
