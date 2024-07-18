package main

import (
	"net/http"
	"time"
)

const expirationTestUser = "expirationtestuser"

func loginHandler(w http.ResponseWriter, r *http.Request) {
	creds, err := readBody[LoginCredentials](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.IsPasswordCorrect(creds.User, creds.Password) {
		// TODO Log
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	// TODO Add cookie renewal logic when used in the checkAuthentication. Once a day and at boot, delete all expired cookies. A user can have one or multiple active cookies?
	cookie, err := generateCookie()
	if err != nil {
		// TODO Log
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if profile == TEST && creds.User == expirationTestUser {
		cookie.Expires = time.Now().UTC().Add(-1 * time.Second)
	}

	err = repo.SetCookie(creds.User, cookie.Value, cookie.Expires)
	if err != nil {
		// TODO Log
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	Logger.Info("user '%s' logged in successfully", creds.User)
	w.WriteHeader(http.StatusOK)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		deleteReceivedUser(w, r)
	} else {
		// TODO Log
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func deleteReceivedUser(w http.ResponseWriter, r *http.Request) {
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	if !repo.DoesUserExist(user) {
		// TODO Log
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	err = repo.DeleteUser(user)
	if err != nil {
		// TODO Log
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Info("deleted user: %s", user)
	w.WriteHeader(http.StatusOK)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	form, err := readBody[RegistrationForm](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if repo.DoesUserExist(form.User) {
		Logger.Info("user '%s' tried to register but he already exists", form.User)
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	err = repo.CreateUser(form)
	if err != nil {
		// TODO Log
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Info("user registered: " + form.User)
	w.WriteHeader(http.StatusOK)
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// TODO Log
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	form, err := readBody[ChangePasswordForm](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(form.User) {
		// TODO Log
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.IsPasswordCorrect(form.User, form.OldPassword) {
		// TODO Log
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	err = repo.ChangePassword(form.User, form.NewPassword)
	if err != nil {
		// TODO Log
		http.Error(w, "error when trying to change password", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' changed his password", form.User)
	w.WriteHeader(http.StatusOK)
}

func changeOriginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// TODO Log
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	form, err := readBody[ChangeOriginForm](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(form.User) {
		// TODO Log
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.IsPasswordCorrect(form.User, form.Password) {
		// TODO Log
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	err = repo.ChangeOrigin(form.User, form.NewOrigin)
	if err != nil {
		// TODO Log
		http.Error(w, "error when trying to change origin", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' changed his origin", form.User)
	w.WriteHeader(http.StatusOK)
}
