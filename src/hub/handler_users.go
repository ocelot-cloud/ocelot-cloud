package main

import (
	"fmt"
	"net/http"
	"time"
)

const expirationTestUser = "expirationtestuser"

type RegistrationForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Origin   string `json:"host"`
}

type LoginCredentials struct {
	User     string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("login logic called")
	creds, err := readBody[LoginCredentials](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !repo.IsPasswordCorrect(creds.User, creds.Password) {
		logAndRespondDebug(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	// TODO Add cookie renewal logic when used in the checkAuthentication. Once a day and at boot, delete all expired cookies. A user can have one or multiple active cookies?
	cookie, err := generateCookie()
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
	}

	if profile == TEST && creds.User == expirationTestUser {
		cookie.Expires = time.Now().UTC().Add(-1 * time.Second)
	}

	// TODO Must be verified by a test:
	err = repo.SetCookie(creds.User, cookie.Value, cookie.Expires)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusOK)
		return
	}

	http.SetCookie(w, cookie)
	logAndRespondDebug(w, "login successful", http.StatusOK)
}

// TODO delete user, get user (maybe for testing?)
func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		deleteReceivedUser(w, r)
	} else {
		logAndRespondDebug(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

// TODO Everywhere: replace "auth" by cookieName

// TODO put it in handler_tools
func checkAuthentication(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		Logger.Debug("cookie not set in request: %s", err.Error())
		logAndRespondDebug(w, "cookie not set in request", http.StatusUnauthorized)
		return "", fmt.Errorf("")
	}

	if !validate(cookie.Value, Cookie) {
		logAndRespondDebug(w, "invalid cookie", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	if !validate(r.Header.Get("Origin"), Origin) {
		logAndRespondDebug(w, "invalid origin", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	// TODO everytime I use "GetUserWithCookie" I should do validation previously. Everytime I do sth with the cookie in general.
	authenticatedUser, err := repo.GetUserWithCookie(cookie.Value)
	if err != nil {
		Logger.Debug("error when getting cookie of user: %s", err.Error())
		http.Error(w, "cookie not found", http.StatusNotFound)
		return "", fmt.Errorf("")
	}

	// TODO there should be a global variable "originHeader"
	if !repo.IsOriginCorrect(authenticatedUser, r.Header.Get("Origin")) {
		logAndRespondDebug(w, "origin not matching", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	if repo.IsCookieExpired(cookie.Value) {
		logAndRespondDebug(w, "cookie expired", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	newExpirationTime := getTimeIn30Days() // TODO Also set exp time request.
	err = repo.SetCookie(authenticatedUser, cookie.Value, newExpirationTime)
	if err != nil {
		logAndRespondDebug(w, "updating cookie failed", http.StatusInternalServerError)
		return "", err
	}
	cookie.Expires = newExpirationTime
	http.SetCookie(w, cookie)

	return authenticatedUser, nil
}

func deleteReceivedUser(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	if !repo.DoesUserExist(authenticatedUser) {
		logAndRespondError(w, "user does not exist", http.StatusNotFound)
		return
	}

	err = repo.DeleteUser(authenticatedUser)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Info("Deleted user: %s", authenticatedUser)

	logAndRespondDebug(w, "User deleted", http.StatusOK)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Should return a pointer
	form, err := readBody[RegistrationForm](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if repo.DoesUserExist(form.Username) {
		logAndRespondError(w, "user already exists", http.StatusConflict)
		return
	}

	err = repo.CreateUser(&form)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logAndRespondDebug(w, "User registered", http.StatusOK)
}

type ChangePasswordForm struct {
	User        string `json:"user"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logAndRespondError(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// TODO Should return a pointer
	form, err := readBody[ChangePasswordForm](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(form.User) {
		logAndRespondDebug(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.IsPasswordCorrect(form.User, form.OldPassword) {
		logAndRespondDebug(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	err = repo.ChangePassword(form.User, form.NewPassword)
	if err != nil {
		logAndRespondError(w, "error when trying to change password", http.StatusInternalServerError)
		return
	}

	logAndRespondDebug(w, "password changed", http.StatusOK)
}

type ChangeOriginForm struct {
	User      string `json:"user"`
	Password  string `json:"password"`
	NewOrigin string `json:"new_origin"`
}

func changeOriginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logAndRespondError(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// TODO Should return a pointer
	form, err := readBody[ChangeOriginForm](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(form.User) {
		logAndRespondDebug(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.IsPasswordCorrect(form.User, form.Password) {
		logAndRespondDebug(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	err = repo.ChangeOrigin(form.User, form.NewOrigin)
	if err != nil {
		logAndRespondError(w, "error when trying to change origin", http.StatusInternalServerError)
		return
	}

	logAndRespondDebug(w, "origin changed", http.StatusOK)
}
