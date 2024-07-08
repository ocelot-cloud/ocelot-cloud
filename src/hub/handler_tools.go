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
func createFileInfo(filename string) (*FileInfo, error) {
	infos := strings.Split(filename, "_")
	if len(infos) != 3 {
		return nil, fmt.Errorf("error, filenames should have exactly two underscores: %s", filename)
	}
	var info = &FileInfo{}
	info.User = infos[0]
	info.App = infos[1]
	// TODO consider error here and test it.
	info.FileName = infos[2]
	info.Tag, _ = strings.CutSuffix(infos[2], ".tar.gz")
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

type FileInfo struct {
	User     string
	App      string
	Tag      string
	FileName string
}

var fs FileStorage = &FileStorageImpl{}
var repo Repository = &SqliteRepository{}

// TODO All functions below require auth
// TODO There must be a "login" handler. When credentials are correct, set a cookie header. -> usually browsers then send that cookie for all subsequent but I have to do that manually

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

	return result, nil
}

// TODO Add security: auth, origin policy and according security tests
// TODO auth: for required actions, some are public like findApps, reuse code from backend?
// TODO origin policy: user creation requires "host" parameter, so all security relevant actions must have this "host" as "Origin" header
// TODO changing the "host" parameter is not origin-protected, but requires password again.
// TODO Add input validation: usernames only lower letters, lengths etc. And test trying to break this with according error messages.
// TODO Restrict maximum space used by user to 10MB
// TODO logging: 1) make sure folder "data" exists. If so, store logs in "data/logs.txt"
// TODO store sqlite.db in "data" folder
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
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("login logic called")
	creds, err := readBody[LoginCredentials](r)
	if err != nil {
		logAndRespondDebug(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !repo.IsPasswordCorrect(creds.Username, creds.Password) {
		logAndRespondDebug(w, "wrong password", http.StatusUnauthorized)
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

	// TODO Must be verified by a test:
	err = repo.SetCookie(creds.Username, cookie.Value, cookie.Expires)
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
