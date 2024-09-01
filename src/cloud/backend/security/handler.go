package security

import (
	"encoding/json"
	"github.com/ocelot-cloud/shared/utils"
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
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

	cookie, err := utils.GenerateCookie()
	if err != nil {
		// TODO
	}
	cookie.SameSite = http.SameSiteLaxMode // TODO Necessary at all? should maybe only be enabled for TEST profile, write tests for it?
	http.SetCookie(w, cookie)

	Logger.Debug("cookie was set")
	w.WriteHeader(http.StatusOK)
}

// TODO Duplication with handleBackendApiRequest
func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	// TODO store generated cookie in a repo and check if their value is correct.
	println("cookie: %s", cookie.Value) // TODO to be removed
	if err != nil {
		Logger.Trace("Cookie error.")
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Logger.Trace("Cookie was okay.")
		w.WriteHeader(http.StatusOK)
	}
}
