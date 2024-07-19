package main

import (
	"net/http"
	"time"
)

const (
	testUserWithExpiredCookie          = "expcookietestuser"
	testUserWithOldButNotExpiredCookie = "oldcookietestuser"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	creds, err := readBody[LoginCredentials](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.IsPasswordCorrect(creds.User, creds.Password) {
		Logger.Info("Password of user '%s' was not correct", creds.User)
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	cookie, err := generateCookie()
	if err != nil {
		Logger.Error("cookie generation failed: %v", err)
		http.Error(w, "cookie generation failed", http.StatusInternalServerError)
	}

	if profile == TEST {
		if creds.User == testUserWithExpiredCookie {
			cookie.Expires = time.Now().UTC().Add(-1 * time.Second)
		} else if creds.User == testUserWithOldButNotExpiredCookie {
			cookie.Expires = time.Now().UTC().Add(24 * time.Hour)
		}
	}

	err = repo.SetCookie(creds.User, cookie.Value, cookie.Expires)
	if err != nil {
		Logger.Error("setting cookie failed: %v", err)
		http.Error(w, "setting cookie failed", http.StatusInternalServerError)
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
		handleInvalidRequestMethod(w, r, userPath)
		return
	}
}

func deleteReceivedUser(w http.ResponseWriter, r *http.Request) {
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	if !repo.DoesUserExist(user) {
		Logger.Error("user '%s' wanted to delete his account but seems not to exist although authenticated", user)
		http.Error(w, "user does not exist", http.StatusInternalServerError)
		return
	}

	err = repo.DeleteUser(user)
	if err != nil {
		Logger.Error("user '%s' deletion failed", err)
		http.Error(w, "user deletion failed", http.StatusInternalServerError)
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
		Logger.Error("user '%s' registration failed: %v", form.User, err)
		http.Error(w, "user registration failed", http.StatusInternalServerError)
		return
	}

	Logger.Info("user registered: " + form.User)
	w.WriteHeader(http.StatusOK)
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleInvalidRequestMethod(w, r, userPath)
		return
	}

	form, err := readBody[ChangePasswordForm](r)
	if err != nil {
		Logger.Warn("could not read request body: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(form.User) {
		Logger.Warn("somebody tried to change password but user '%s' does not exist", form.User)
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.IsPasswordCorrect(form.User, form.OldPassword) {
		Logger.Info("incorrect credentials for user '%s' when trying to change password", form.User)
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	err = repo.ChangePassword(form.User, form.NewPassword)
	if err != nil {
		Logger.Error("changing password for user '%s' failed: %v", form.User, err)
		http.Error(w, "error when trying to change password", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' changed his password", form.User)
	w.WriteHeader(http.StatusOK)
}

func changeOriginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleInvalidRequestMethod(w, r, userPath)
		return
	}

	form, err := readBody[ChangeOriginForm](r)
	if err != nil {
		Logger.Warn("invalid input: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(form.User) {
		Logger.Warn("user '%s' tried to change origin but he does not exist", form.User)
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.IsPasswordCorrect(form.User, form.Password) {
		Logger.Info("incorrect credentials for user '%s' when trying to change origin", form.User)
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	err = repo.ChangeOrigin(form.User, form.NewOrigin)
	if err != nil {
		Logger.Error("changing origin for user '%s' failed: %v", form.User, err)
		http.Error(w, "error when trying to change origin", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' changed his origin", form.User)
	w.WriteHeader(http.StatusOK)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleInvalidRequestMethod(w, r, userPath)
		return
	}

	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	err = repo.Logout(user)
	if err != nil {
		Logger.Error("logout of user '%s' failed: %v", user, err)
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' logged out", user)
	w.WriteHeader(http.StatusOK)
}

// TODO I think the validator should return nil/err for improved logging messages.
// TODO During testing the log level should be DEBUG, by default, it should be INFO

// TODO Implement Frontend, do I have to add CORS policy to allow cross domain access?

// TODO When upload is implemented in hub, then I can delete alls the stacks in the cloud. Acceptance tests need to integrate hub and need to implement download of stacks at the beginning?
// TODO In the end, add deploy script which only works on my device, since I have the correct SSH keys and config.
