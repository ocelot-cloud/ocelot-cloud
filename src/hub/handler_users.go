package main

import (
	"net/http"
	"time"
)

const OriginHeader = "Origin"
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

	err = repo.SetCookie(creds.User, cookie.Value, cookie.Expires)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	logAndRespondDebug(w, "login successful", http.StatusOK)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		deleteReceivedUser(w, r)
	} else {
		logAndRespondDebug(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func deleteReceivedUser(w http.ResponseWriter, r *http.Request) {
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	if !repo.DoesUserExist(user) {
		logAndRespondError(w, "user does not exist", http.StatusNotFound)
		return
	}

	err = repo.DeleteUser(user)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Info("Deleted user: %s", user)

	logAndRespondDebug(w, "User deleted", http.StatusOK)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	form, err := readBody[RegistrationForm](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if repo.DoesUserExist(form.Username) {
		logAndRespondError(w, "user already exists", http.StatusConflict)
		return
	}

	err = repo.CreateUser(form)
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
