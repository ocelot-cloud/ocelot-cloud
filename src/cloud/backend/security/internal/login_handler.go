package internal

import (
	"encoding/json"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TODO Insecure
var users = map[string]string{
	"admin": "password",
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("login logic called")
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		Logger.Debug("password not matching")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO Use safe, randomly generated cookies instead. I think gorilla provides some.
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "valid",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}
