package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func sendJsonResponse(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		Logger.Error("unmarshalling failed: %v", err)
		http.Error(w, "TODO", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		// TODO
		logAndRespondDebug(w, err.Error(), http.StatusInternalServerError)
	}
}

func logAndRespondWarn(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Warn(msg)
	http.Error(w, msg, httpStatus)
}

func logAndRespondInfo(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Info(msg)
	http.Error(w, msg, httpStatus)
}

func logAndRespondError(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Error(msg)
	http.Error(w, msg, httpStatus)
}

func logAndRespondDebug(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Debug(msg)
	http.Error(w, msg, httpStatus)
}

var repo Repository = &SqliteRepository{}

func readBody[T any](r *http.Request) (*T, error) {
	var result T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read request body: %w", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	switch v := any(result).(type) {
	case UserAndApp:
		if !validate(v.User, User) || !validate(v.App, App) {
			return nil, fmt.Errorf("invalid input")
		}
	case RegistrationForm:
		if !validate(v.User, User) || !validate(v.Password, Password) || !validate(v.Email, Email) || !validate(v.Origin, Origin) {
			return nil, fmt.Errorf("invalid input")
		}
	case ChangeOriginForm:
		if !validate(v.User, User) || !validate(v.Password, Password) || !validate(v.NewOrigin, Origin) {
			return nil, fmt.Errorf("invalid input")
		}
	case ChangePasswordForm:
		if !validate(v.User, User) || !validate(v.OldPassword, Password) || !validate(v.NewPassword, Password) {
			return nil, fmt.Errorf("invalid input")
		}
	case LoginCredentials:
		if !validate(v.User, User) || !validate(v.Password, Password) {
			return nil, fmt.Errorf("invalid input")
		}
	case TagInfo:
		if !validate(v.User, User) || !validate(v.App, App) || !validate(v.Tag, Tag) {
			return nil, fmt.Errorf("invalid input")
		}
	}

	return &result, nil
}

func readBodyAsSingleString(r *http.Request, validationType ValidationType) (string, error) {
	singleString, err := readBody[SingleString](r)
	if err != nil {
		return "", err
	}
	result := singleString.Value

	if !validate(result, validationType) {
		return "", fmt.Errorf("invalid input")
	}

	return result, nil
}

// TODO Restrict maximum space used by user to 10MB
// TODO Introduce ENV variable "DISABLE_EMAIL_VERIFICATION", default is false.
//  Disable for development. If enabled, I think I should throw an error if it
//  did not got email stuff. If it got them, It will run a test to check whether
//  it works. If so, start normally. If not, it exits immediately.
//  If it exits early, make a short explanation on why it does that. "printEmailExplanation"?

func getTimeIn30Days() time.Time {
	return time.Now().UTC().AddDate(0, 0, 30)
}

func generateCookie() (*http.Cookie, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		Logger.Error("Failed to generate cookie: %v", err)
		return nil, err
	}
	return &http.Cookie{
		Name:    cookieName,
		Value:   hex.EncodeToString(bytes),
		Expires: getTimeIn30Days(),
	}, nil
}

func wipeDataHandler(w http.ResponseWriter, r *http.Request) {
	repo.WipeDatabase()
	Logger.Warn("database wipe completed")
	w.WriteHeader(http.StatusOK)
}

func checkAuthentication(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		Logger.Debug("cookie not set in request: %s", err.Error())
		logAndRespondDebug(w, "cookie not set in request", http.StatusUnauthorized)
		return "", fmt.Errorf("")
	}

	if !validate(cookie.Value, Cookie) {
		logAndRespondDebug(w, "invalid cookie", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	if !validate(r.Header.Get(OriginHeader), Origin) {
		logAndRespondDebug(w, "invalid origin", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	authenticatedUser, err := repo.GetUserWithCookie(cookie.Value)
	if err != nil {
		Logger.Debug("error when getting cookie of user: %s", err.Error())
		http.Error(w, "cookie not found", http.StatusNotFound)
		return "", fmt.Errorf("")
	}

	if !repo.IsOriginCorrect(authenticatedUser, r.Header.Get(OriginHeader)) {
		logAndRespondDebug(w, "origin not matching", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	if repo.IsCookieExpired(cookie.Value) {
		logAndRespondDebug(w, "cookie expired", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	newExpirationTime := getTimeIn30Days() // TODO Also set exp time request.
	err = repo.SetCookie(authenticatedUser, cookie.Value, newExpirationTime)
	if err != nil {
		logAndRespondDebug(w, "updating cookie failed", http.StatusInternalServerError)
		return "", err
	}
	cookie.Expires = newExpirationTime
	http.SetCookie(w, cookie)

	return authenticatedUser, nil
}
