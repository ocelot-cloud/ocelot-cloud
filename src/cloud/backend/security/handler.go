package security

import (
	"encoding/json"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
	"ocelot/backend/tools"
)

// TODO Can be abstracted in shared module?
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		Logger.Info("decoding credentials failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !userRepo.IsPasswordCorrect(creds.Username, creds.Password) {
		Logger.Info("password of user '%s' not matching", creds.Username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie, err := utils.GenerateCookie()
	if err != nil {
		Logger.Error("generating cookie failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie.Name = tools.CookieName
	cookie.SameSite = http.SameSiteLaxMode // TODO Necessary at all? should maybe only be enabled for TEST profile, write tests for it?
	http.SetCookie(w, cookie)

	Logger.Debug("login successful")
	w.WriteHeader(http.StatusOK)
}

// TODO Duplication with handleBackendApiRequest
func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO store generated cookie in a repo and check if their value is correct.
	_, err := r.Cookie(tools.CookieName)
	if err != nil {
		Logger.Debug("Cookie error: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Logger.Debug("Cookie was okay.")
		w.WriteHeader(http.StatusOK)
	}
}
