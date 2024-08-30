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

	if !repo.DoesUserExist(creds.User) {
		Logger.Info("user '%s' does not exist", creds.User)
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.IsPasswordCorrect(creds.User, creds.Password) {
		Logger.Info("Password of user '%s' was not correct", creds.User)
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	err = repo.SetOrigin(creds.User, creds.Origin)
	if err != nil {
		Logger.Error("setting origin failed: %v", err)
		http.Error(w, "setting origin failed", http.StatusInternalServerError)
		return
	}

	cookie, err := generateCookie()
	if err != nil {
		Logger.Error("cookie generation failed: %v", err)
		http.Error(w, "cookie generation failed", http.StatusInternalServerError)
		return
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

func authCheckHandler(w http.ResponseWriter, r *http.Request) {
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}
	sendJsonResponse(w, SingleString{user})
}

func userDeleteHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := checkAuthentication(w, r)
	if err != nil {
		return
	}

	form, err := readBody[ChangePasswordForm](r)
	if err != nil {
		Logger.Warn("could not read request body: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if !repo.DoesUserExist(user) {
		Logger.Warn("somebody tried to change password but user '%s' does not exist", user)
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	if !repo.IsPasswordCorrect(user, form.OldPassword) {
		Logger.Info("incorrect credentials for user '%s' when trying to change password", user)
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	err = repo.ChangePassword(user, form.NewPassword)
	if err != nil {
		Logger.Error("changing password for user '%s' failed: %v", user, err)
		http.Error(w, "error when trying to change password", http.StatusInternalServerError)
		return
	}

	Logger.Info("user '%s' changed his password", user)
	w.WriteHeader(http.StatusOK)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
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
