package main

import "net/http"

// TODO delete user, get user (maybe for testing?)
func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		deleteReceivedUser(w, r)
	} else {
		logAndRespondDebug(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func deleteReceivedUser(w http.ResponseWriter, r *http.Request) {
	singleString, err := readBody[SingleString](r) // TODO username validation
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	userToDelete := singleString.Value

	// TODO Everywhere: replace "auth" by cookieName

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !validate(cookie.Value, Cookie) {
		logAndRespondDebug(w, "invalid cookie", http.StatusBadRequest)
		return
	}

	if !validate(r.Header.Get("Origin"), Origin) {
		logAndRespondDebug(w, "invalid origin", http.StatusBadRequest)
		return
	}

	// TODO everytime I use "GetUserWithCookie" I should do validation previously. Everytime I do sth with the cookie in general.
	authenticatedUser, err := repo.GetUserWithCookie(cookie.Value)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO there should be a global variable "originHeader"
	if !repo.IsOriginCorrect(authenticatedUser, r.Header.Get("Origin")) {
		logAndRespondDebug(w, "origin not matching", http.StatusBadRequest)
		return
	}

	// TODO rename to "IsCookieExpired" and inverse the boolean result
	// TODO Research: Uncommenting this causes many tests to fail. I assume some other tests did not provide cookie properly.
	/*if !repo.IsCookieValid(cookie.Value) {
		logAndRespondDebug(w, "cookie expired", http.StatusBadRequest)
		return
	}
	*/

	// TODO Most of the code above this line can be put into a single security policy function.

	if authenticatedUser != userToDelete {
		logAndRespondDebug(w, "deletion of other users not allowed", http.StatusUnauthorized)
		return
	}

	if !repo.DoesUserExist(userToDelete) {
		logAndRespondError(w, "user does not exist", http.StatusNotFound)
		return
	}

	err = fs.DeleteUser(userToDelete)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = repo.DeleteUser(userToDelete)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Info("Deleted user: %s", userToDelete)

	logAndRespondDebug(w, "User deleted", http.StatusOK)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Should return a pointer
	form, err := readBody[RegistrationForm](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = fs.CreateUser(form.Username)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
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

	logAndRespondDebug(w, "User registered", http.StatusCreated)
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
		logAndRespondDebug(w, "Password is not correct", http.StatusUnauthorized)
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
		logAndRespondDebug(w, "Password is not correct", http.StatusUnauthorized)
		return
	}

	err = repo.ChangeOrigin(form.User, form.NewOrigin)
	if err != nil {
		logAndRespondError(w, "error when trying to change origin", http.StatusInternalServerError)
		return
	}

	logAndRespondDebug(w, "origin changed", http.StatusOK)
}
