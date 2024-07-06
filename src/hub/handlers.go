package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func logAndRespondError(w http.ResponseWriter, msg string, httpStatus int) {
	Logger.Error(msg)
	http.Error(w, msg, httpStatus)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	uploadName := strings.TrimPrefix(r.URL.Path, downloadPath)
	if uploadName == "" {
		logAndRespondError(w, "File name is missing", http.StatusBadRequest)
		return
	}

	fileInfo, err := createFileInfo(uploadName)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	path := fmt.Sprintf("%s/%s/%s/%s", usersDir, fileInfo.User, fileInfo.App, fileInfo.FileName)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		logAndRespondError(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, path)
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

// TODO delete user, get user (maybe for testing?)
func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		deleteReceivedUser(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	user, err := readBody[RegistrationForm](r)
	if err != nil {
		// TODO
	}

	err = fs.CreateUser(user.Username)
	if err != nil {
		Logger.Error("error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not create user on filesystem"))
		return
	}

	// TODO Handle error
	repo.CreateUser(user.Username, user.Password)
	Logger.Info("Created user: %s", user.Username)

	// TODO Handle error
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered"))
}

type SingleString struct {
	Value string `json:"name"`
}

func deleteReceivedUser(w http.ResponseWriter, r *http.Request) {
	singleString, err := readBody[SingleString](r)
	if err != nil {
		// TODO
	}
	user := singleString.Value

	// TODO Misses some functions like: Does(User/App/Tag)Exist?
	fs.DeleteUser(user) // TODO Shouldn't that return a potential error?
	repo.DeleteUser(user)

	Logger.Info("Deleted user: %s", user)

	w.WriteHeader(http.StatusOK)
	// TODO Handle error
	w.Write([]byte("User deleted"))
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		findApps(w, r)
	} else if r.Method == http.MethodPost {
		createApp(w, r)
	} else {
		// TODO
	}
	// TODO create/delete app, search for app: search
}

// TODO Much duplication of the handler logic. Should be abstracted.

func findApps(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("finding apps")
	singleString, err := readBody[SingleString](r)
	if err != nil {
		// TODO
	}
	searchTerm := singleString.Value

	apps, err := repo.FindApps(searchTerm)
	if err != nil {
		Logger.Error("Finding apps failed: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred when trying to find apps"))
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(apps)
	if err != nil {
		Logger.Error("Failed to marshal apps: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred when trying to marshal apps"))
		return
	}
	w.Write(jsonData)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie == nil || cookie.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cookie not contained in request"))
		return
	}
	user, err := repo.GetUserWithCookie(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Cookie not found"))
		return
	}

	singleString, err := readBody[SingleString](r)
	if err != nil {
		// TODO
	}
	app := singleString.Value

	if !repo.DoesUserExist(user) {
		// TODO
	}
	if repo.DoesAppExist(user, app) {
		// TODO
	}
	err = fs.CreateApp(user, app)
	if err != nil {
		// TODO
		return
	}
	err = repo.CreateApp(user, app)
	if err != nil {
		// TODO
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("created app successfully"))
	if err != nil {
		// TODO
	}
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

func tagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleUpload(w, r)
	}

	// TODO delete tag, getListOfTags(app)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logAndRespondError(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logAndRespondError(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// TODO Add test
	// TODO Make security test that user and repo are in the name correctly, and that both exist.
	if !strings.HasSuffix(header.Filename, ".tar.gz") {
		logAndRespondError(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	fileInfo, err := createFileInfo(header.Filename)
	if err != nil {
		logAndRespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fileBuffer bytes.Buffer
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		logAndRespondError(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	// TODO Should be global?
	fs := FileStorageImpl{}
	err = fs.CreateTag(fileInfo, &fileBuffer)
	if err != nil {
		logAndRespondError(w, "Failed to write content to local file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	Logger.Info("File uploaded successfully: %s", header.Filename)
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
		// TODO
	}

	// TODO verify username+password
	// TODO Use safe, randomly generated cookies instead. I think gorilla provides some.
	// TODO Add cookie + expiration time/date to sqlite to survive restarts.
	// TODO Add cookie renewal logic when used in the middleware. Once a day and at boot, delete all expired cookies. A user can have one or multiple active cookies?
	// TODO In the tests, check that cookie has correct length and has different value when requesting a seconds one.
	cookie, err := generateCookie()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// TODO Must be verified by a test:
	err = repo.SetCookie(creds.Username, cookie.Value, cookie.Expires)
	if err != nil {
		// TODO
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("login successful"))
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
