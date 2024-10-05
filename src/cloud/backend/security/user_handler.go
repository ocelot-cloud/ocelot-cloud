package security

import (
	"context"
	"encoding/json"
	"fmt"
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

func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	err := GetAuthentication(w, r)
	if err != nil {
		return
	}
	Logger.Debug("Cookie was okay.")
	w.WriteHeader(http.StatusOK)
}

func GetAuthentication(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(tools.CookieName)
	if err != nil {
		Logger.Info("cookie error: %v", err)
		http.Error(w, "error with request cookie", http.StatusUnauthorized)
		return fmt.Errorf("")
	}

	auth, err := UserRepo.GetUserViaCookie(cookie.Value)
	if err != nil {
		http.Error(w, "cookie not found", http.StatusUnauthorized)
		return fmt.Errorf("")
	}

	// The context information is not used in an actual subsequent http request, such as when the request is proxied to
	// an application.
	ctx := context.WithValue(r.Context(), "auth", auth)
	r = r.WithContext(ctx)

	return nil
}

// TODO must be authenticated
func getSecretHandler(w http.ResponseWriter, r *http.Request) {
	// TODO read user from context, generate secret and return it
}

// TODO Add cookie expiration checks
