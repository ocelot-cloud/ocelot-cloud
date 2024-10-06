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
	_, err := GetRequestWithAuthContext(w, r)
	if err != nil {
		return
	}
	Logger.Debug("Cookie was okay.")
	w.WriteHeader(http.StatusOK)
}

func GetRequestWithAuthContext(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	auth, err := GetAuthentication(w, r)
	if err != nil {
		return nil, err
	}

	// The context information is not used in an actual subsequent http request, such as when the request is proxied to
	// an application.
	ctx := context.WithValue(r.Context(), "auth", auth)
	r = r.WithContext(ctx)

	return r, nil
}

func GetAuthentication(w http.ResponseWriter, r *http.Request) (*tools.Authorization, error) {
	cookie, err := r.Cookie(tools.CookieName)
	if err != nil {
		Logger.Info("cookie error: %v", err)
		http.Error(w, "error with request cookie", http.StatusUnauthorized)
		return nil, fmt.Errorf("")
	}

	auth, err := UserRepo.GetUserViaCookie(cookie.Value)
	if err != nil {
		http.Error(w, "cookie not found", http.StatusUnauthorized)
		return nil, fmt.Errorf("")
	}
	return auth, nil
}

// TODO must be authenticated
func SecretHandler(w http.ResponseWriter, r *http.Request) {
	auth, err := tools.GetAuthFromContext(w, r)
	if err != nil {
		return
	}

	secret, err := UserRepo.GenerateSecret(auth.User)
	if err != nil {
		http.Error(w, "secret generation failed", http.StatusInternalServerError)
		return
	}
	Logger.Debug("Secret generated: " + secret) // TODO temp

	utils.SendJsonResponse(w, secret) // TODO I should use that more often instead of w.writeHeader and w.write?
}

// TODO Add cookie expiration checks and renewals.
