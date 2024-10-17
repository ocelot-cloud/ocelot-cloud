package main

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/utils"
	"io"
	"net/http"
)

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
	case App:
		jobs = []ValidationJob{
			{v.Maintainer, User},
			{v.App, AppType},
		}
	case RegistrationForm:
		jobs = []ValidationJob{
			{v.User, User},
			{v.Password, Password},
			{v.Email, Email},
		}
	case utils.ChangePasswordForm:
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
			{v.App, AppType},
			{v.Tag, Tag},
		}
	case AppAndTag:
		jobs = []ValidationJob{
			{v.App, AppType},
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
	singleString, err := readBody[utils.SingleString](r)
	if err != nil {
		return "", err
	}
	result := singleString.Value

	if err = validate(result, validationType); err != nil {
		return "", err
	}

	return result, nil
}

func wipeDataHandler(w http.ResponseWriter, r *http.Request) {
	repo.WipeDatabase()
	Logger.Warn("database wipe completed")
	w.WriteHeader(http.StatusOK)
}

func checkAuthentication(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		Logger.Info("cookie not set in request: %s", err.Error())
		http.Error(w, "cookie not set in request", http.StatusUnauthorized)
		return "", fmt.Errorf("")
	}

	if err = validate(cookie.Value, Cookie); err != nil {
		http.Error(w, "invalid cookie", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	if err = validate(r.Header.Get(utils.OriginHeader), Origin); err != nil {
		http.Error(w, "invalid origin", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	user, err := repo.GetUserWithCookie(cookie.Value)
	if err != nil {
		Logger.Warn("error when getting cookie of user: %s", err.Error())
		http.Error(w, "cookie not found", http.StatusUnauthorized)
		return "", fmt.Errorf("")
	}

	if !repo.IsOriginCorrect(user, r.Header.Get(utils.OriginHeader)) {
		Logger.Warn("user '%s' used a not matching origin: '%s'", user, r.Header.Get(utils.OriginHeader))
		http.Error(w, "origin not matching", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	if repo.IsCookieExpired(cookie.Value) {
		Logger.Warn("user '%s' used an expired cookie'", user)
		http.Error(w, "cookie expired", http.StatusBadRequest)
		return "", fmt.Errorf("")
	}

	newExpirationTime := utils.GetTimeIn30Days()
	err = repo.HashAndSaveCookie(user, cookie.Value, newExpirationTime)
	if err != nil {
		Logger.Error("setting new cookie failed: %v", err)
		http.Error(w, "setting new cookie failed", http.StatusInternalServerError)
		return "", fmt.Errorf("")
	}
	cookie.Expires = newExpirationTime
	// Note: If no path is given, browsers set the default path one level higher than the
	// request path. For example, calling "/a" sets the cookie path to two "/", and calling
	// "/a/b" sets the cookie path to "/a". When updating a cookie, two cookies, the old one
	// and the updated one, with different paths are stored in the browser, causing some
	// requests to fail with "cookie not found".
	cookie.Path = "/"
	http.SetCookie(w, cookie)

	return user, nil
}
