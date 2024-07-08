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
	singleString, err := readBody[SingleString](r)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := singleString.Value

	if !repo.DoesUserExist(user) {
		logAndRespondError(w, "user does not exist", http.StatusNotFound)
		return
	}

	err = fs.DeleteUser(user)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
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
