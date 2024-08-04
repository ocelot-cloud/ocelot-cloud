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

const OriginHeader = "Origin"

func sendJsonResponse(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		Logger.Error("unmarshalling failed: %v", err)
		http.Error(w, "failed to prepare response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

var repo Repository = &SqliteRepository{}

type ValidationJob struct {
	Value   string
	valType ValidationType
}

func readBody[T any](r *http.Request) (*T, error) {
	var result T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read request body: %w", err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	var jobs []ValidationJob
	switch v := any(result).(type) {
	case UserAndApp:
		jobs = []ValidationJob{
			{v.User, User},
			{v.App, App},
		}
	case RegistrationForm:
		jobs = []ValidationJob{
			{v.User, User},
			{v.Password, Password},
			{v.Email, Email},
		}
	case ChangePasswordForm:
		jobs = []ValidationJob{
			{v.OldPassword, Password},
			{v.NewPassword, Password},
		}
	case LoginCredentials:
		jobs = []ValidationJob{
			{v.User, User},
			{v.Password, Password},
			{v.Origin, Origin},
		}
	case TagInfo:
		jobs = []ValidationJob{
			{v.User, User},
			{v.App, App},
			{v.Tag, Tag},
		}
	case AppAndTag:
		jobs = []ValidationJob{
			{v.App, App},
			{v.Tag, Tag},
		}
	}
	if err = validateJobs(jobs); err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func validateJobs(jobs []ValidationJob) error {
	for _, job := range jobs {
		if err := validate(job.Value, job.valType); err != nil {
			return err
		}
	}
	return nil
}

func readBodyAsSingleString(r *http.Request, validationType ValidationType) (string, error) {
	singleString, err := readBody[SingleString](r)
	if err != nil {
		return "", err
	}
	result := singleString.Value

	if err := validate(result, validationType); err != nil {
		return "", err
	}

	return result, nil
}

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
		Name:     cookieName,
		Value:    hex.EncodeToString(bytes),
		Expires:  getTimeIn30Days(),
		SameSite: http.SameSiteLaxMode,
	}, nil
}

func wipeDataHandler(w http.ResponseWriter, r *http.Request) {
	repo.WipeDatabase()
	Logger.Warn("database wipe completed")
	w.WriteHeader(http.StatusOK)
}

func checkAuthentication(w http.ResponseWriter, r *http.Request) (string, error) {
	return doAuthenticationCheck(w, r, true)
}

func doAuthenticationCheck(w http.ResponseWriter, r *http.Request, writeHttpError bool) (string, error) {
	user, httpMsg, status := asdf(w, r)
	if status != 200 {
		if writeHttpError {
			http.Error(w, httpMsg, status)
		}
		return "", fmt.Errorf("")
	}
	return user, nil
}

func asdf(w http.ResponseWriter, r *http.Request) (string, string, int) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		Logger.Info("cookie not set in request: %s", err.Error())
		return "", "cookie not set in request", http.StatusUnauthorized
	}

	if err = validate(cookie.Value, Cookie); err != nil {
		return "", "invalid cookie", http.StatusBadRequest
	}

	if err = validate(r.Header.Get(OriginHeader), Origin); err != nil {
		return "", "invalid origin", http.StatusBadRequest
	}

	user, err := repo.GetUserWithCookie(cookie.Value)
	if err != nil {
		Logger.Warn("error when getting cookie of user: %s", err.Error())
		return "", "cookie not found", http.StatusNotFound
	}

	if !repo.IsOriginCorrect(user, r.Header.Get(OriginHeader)) {
		Logger.Warn("user '%s' used a not matching origin: '%s'", user, r.Header.Get(OriginHeader))
		return "", "origin not matching", http.StatusBadRequest
	}

	if repo.IsCookieExpired(cookie.Value) {
		Logger.Warn("user '%s' used an expired cookie'", user)
		return "", "cookie expired", http.StatusBadRequest
	}

	newExpirationTime := getTimeIn30Days()
	err = repo.SetCookie(user, cookie.Value, newExpirationTime)
	if err != nil {
		Logger.Error("setting new cookie failed: %v", err)
		return "", "setting new cookie failed", http.StatusInternalServerError
	}
	cookie.Expires = newExpirationTime
	// Note: If no path is given, browsers set the default path one level higher than the
	// request path. For example, calling "/a" sets the cookie path to two "/", and calling
	// "/a/b" sets the cookie path to "/a". When updating a cookie, two cookies, the old one
	// and the updated one, with different paths are stored in the browser, causing some
	// requests to fail with "cookie not found".
	cookie.Path = "/"
	http.SetCookie(w, cookie)

	return user, "", 200
}

func handleInvalidRequestMethod(w http.ResponseWriter, r *http.Request, endpoint string) {
	Logger.Warn("invalid request method '%s' on endpoint '%s'", r.Method, endpoint)
	http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
}
