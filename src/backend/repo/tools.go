package repo

import (
	"fmt"
	"net/http"
	"ocelot/backend/tools"
)

var Logger = tools.Logger

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
