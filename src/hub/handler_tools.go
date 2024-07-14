package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TODO Requires Auth
// TODO Only allowed when the target is the user itself. Cant upload stuff to other users.

// TODO Create unit tests
func createAppAndTag(filename string) (*AppAndTag, error) {
	if !strings.HasSuffix(filename, ".tar.gz") {
		return nil, fmt.Errorf("error, filename must end with .tar.gz")
	}
	infos := strings.Split(filename, "_")
	if len(infos) != 2 {
		return nil, fmt.Errorf("error, filenames should have exactly two underscores: %s", filename)
	}
	var info = &AppAndTag{}
	info.App = infos[0]
	info.Tag, _ = strings.CutSuffix(infos[1], ".tar.gz")
	return info, nil
}

func sendJsonResponse(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusInternalServerError)
	}
}

func logAndRespondError(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Error(msg)
	http.Error(w, msg, httpStatus)
}

func logAndRespondDebug(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Debug(msg)
	http.Error(w, msg, httpStatus)
}

type AppAndTag struct {
	App string
	Tag string
}

var fs FileStorage = &FileStorageImpl{}
var repo Repository = &SqliteRepository{}

type SingleString struct {
	Value string `json:"name"`
}

func readBody[T any](r *http.Request) (T, error) {
	var result T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return result, fmt.Errorf("unable to read request body: %w", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("invalid request body: %w", err)
	}

	switch v := any(result).(type) {
	case UserAndApp:
		if !validate(v.User, User) || !validate(v.App, App) {
			return result, fmt.Errorf("invalid input")
		}
	case RegistrationForm:
		if !validate(v.Username, User) || !validate(v.Password, Password) || !validate(v.Email, Email) || !validate(v.Origin, Origin) {
			return result, fmt.Errorf("invalid input")
		}
	case ChangeOriginForm:
		if !validate(v.User, User) || !validate(v.Password, Password) || !validate(v.NewOrigin, Origin) {
			return result, fmt.Errorf("invalid input")
		}
	case ChangePasswordForm:
		if !validate(v.User, User) || !validate(v.OldPassword, Password) || !validate(v.NewPassword, Password) {
			return result, fmt.Errorf("invalid input")
		}
	case LoginCredentials:
		if !validate(v.User, User) || !validate(v.Password, Password) {
			return result, fmt.Errorf("invalid input")
		}
	case AppInfo:
		if !validate(v.User, User) || !validate(v.App, App) {
			return result, fmt.Errorf("invalid input")
		}
	case TagInfo:
		if !validate(v.User, User) || !validate(v.App, App) || !validate(v.Tag, Tag) {
			return result, fmt.Errorf("invalid input")
		}
	}

	return result, nil
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

// TODO Should be used in UserHandler
type RegistrationForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Origin   string `json:"host"`
}

type LoginCredentials struct {
	User     string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("login logic called")
	creds, err := readBody[LoginCredentials](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !repo.IsPasswordCorrect(creds.User, creds.Password) {
		logAndRespondDebug(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	// TODO verify username+password
	// TODO Use safe, randomly generated cookies instead. I think gorilla provides some.
	// TODO Add cookie + expiration time/date to sqlite to survive restarts.
	// TODO Add cookie renewal logic when used in the middleware. Once a day and at boot, delete all expired cookies. A user can have one or multiple active cookies?
	// TODO In the tests, check that cookie has correct length and has different value when requesting a seconds one.
	cookie, err := generateCookie()
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusInternalServerError)
	}

	// TODO And the global profile TEST must be activated.
	if creds.User == "expirationtestuser" {
		cookie.Expires = time.Now().UTC().Add(-1 * time.Second)
	}

	// TODO Must be verified by a test:
	err = repo.SetCookie(creds.User, cookie.Value, cookie.Expires)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusOK)
		return
	}

	http.SetCookie(w, cookie)
	logAndRespondDebug(w, "login successful", http.StatusOK)
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
		Name:    cookieName,
		Value:   hex.EncodeToString(bytes),
		Expires: getTimeIn30Days(),
	}, nil
}

// TODO Add acceptance test checking that this endpoint is not available when using production profile.
func wipeDataHandler(w http.ResponseWriter, r *http.Request) {
	fs.WipeStorage()
	repo.WipeDatabase()
	logAndRespondDebug(w, "wipe completed", http.StatusOK)
}
