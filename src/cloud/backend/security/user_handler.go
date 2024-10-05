package security

import (
	"encoding/json"
	"github.com/ocelot-cloud/shared/utils"
	"net/http"
	"ocelot/backend/tools"
	"time"
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

	if !UserRepo.IsPasswordCorrect(creds.Username, creds.Password) {
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

	err = UserRepo.HashAndSaveCookie(creds.Username, cookie.Value, time.Now())
	if err != nil {
		http.Error(w, "saving cookie failed", http.StatusInternalServerError)
		return
	}

	Logger.Debug("login successful")
	w.WriteHeader(http.StatusOK)
}

// TODO Duplication with handleBackendApiRequest?
func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(tools.CookieName)
	if err != nil {
		Logger.Debug("Cookie error: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = UserRepo.GetUserViaCookie(cookie.Value)
	if err != nil {
		http.Error(w, "cookie not found", http.StatusUnauthorized)
		return
	}

	Logger.Debug("Cookie was okay.")
	w.WriteHeader(http.StatusOK)
}

// TODO must be authenticated
func getSecretHandler(w http.ResponseWriter, r *http.Request) {
	// TODO read user from context, generate secret and return it
}

// TODO Add cookie expiration checks
