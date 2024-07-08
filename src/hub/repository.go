package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// TODO use sqlite to store users, apps (for search), password etc.

var db *sql.DB

func initializeDatabaseWithSource(dataSourceName string) {
	// TODO Add database scheme version?
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		Logger.Fatal("Failed to open database: %v\n", err)
	}

	// TODO add: origin TEXT UNIQUE NOT NULL,
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
    		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_name TEXT UNIQUE NOT NULL,
			hashed_password TEXT NOT NULL,
			cookie TEXT,
			expiration_date TEXT
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create users table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS apps (
			app_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			app_name TEXT,
			UNIQUE(user_id, app_name),
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create apps table: %v\n", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			tag_id INTEGER PRIMARY KEY AUTOINCREMENT,
			app_id INTEGER,
			tag_name TEXT,
			UNIQUE(app_id, tag_id),
			FOREIGN KEY (app_id) REFERENCES apps(app_id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create tags table: %v\n", err)
	}

	Logger.Info("Database initialized")
	// TODO Add initial schemes. With version number table.
}

type Repository interface {
	CreateUser(form *RegistrationForm) error
	DoesUserExist(user string) bool
	DeleteUser(user string) error
	IsPasswordCorrect(user string, password string) bool
	DoesAppExist(user string, app string) bool
	CreateApp(user string, app string) error
	DeleteApp(user string, app string) error
	FindApps(query string) ([]AppInfo, error)
	SetCookie(user string, cookie string, expirationDate time.Time) error
	IsCookieValid(cookie string) bool
	GetUserWithCookie(cookie string) (string, error)
	CreateTag(user string, app string, tag string) error
	DeleteTag(user string, app string, tag string) error
	GetTagList(user string, app string) ([]string, error)
	ChangePassword(user string, newPassword string) error
	ChangeOrigin(user string, newOrigin string) error
}

type SqliteRepository struct{}

func (u *SqliteRepository) IsPasswordCorrect(user string, password string) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT hashed_password FROM users WHERE user_name = ?", user).Scan(&hashedPassword)
	if err != nil {
		Logger.Error("Failed to fetch hashed password: %v\n", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (u *SqliteRepository) DoesUserExist(user string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_name = ?)", user).Scan(&exists)
	if err != nil {
		Logger.Error("Failed to check user existence: %v\n", err)
		return false
	}
	return exists
}

func (u *SqliteRepository) CreateUser(form *RegistrationForm) error {

	hashedPassword, err := hashAndSaltPassword(form.Password)
	if err != nil {
		return Logger.LogAndReturnError("Failed to hash password: %v\n", err)
	}

	// TODO Previously check whether user already exists? Here or in handler?
	_, err = db.Exec("INSERT INTO users (user_name, hashed_password) VALUES (?, ?)", form.Username, hashedPassword)
	if err != nil {
		return Logger.LogAndReturnError("Failed to create user: %v", err)
	}
	return nil
}

func hashAndSaltPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *SqliteRepository) DeleteUser(user string) error {
	if !u.DoesUserExist(user) {
		return Logger.LogAndReturnError("User %s does not exist", user)
	}

	_, err := db.Exec("DELETE FROM users WHERE user_name = ?", user)
	if err != nil {
		return Logger.LogAndReturnError("Failed to delete user: %v", err)
	}

	return nil
}

func (u *SqliteRepository) CreateApp(user string, app string) error {
	if !u.DoesUserExist(user) {
		return Logger.LogAndReturnError("User '%s' does not exist", user)
	} else if u.DoesAppExist(user, app) {
		return Logger.LogAndReturnError("App '%s' already exists for user '%s'", app, user)
	}

	userID, err := getUserId(user)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO apps (user_id, app_name) VALUES (?, ?)`, userID, app)
	if err != nil {
		return Logger.LogAndReturnError("Failed to add app '%s' for user '%s': %v", app, user, err)
	}
	return nil
}

func (u *SqliteRepository) DoesAppExist(user string, app string) bool {
	userID, err := getUserId(user)
	if err != nil {
		return false
	}
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM apps WHERE user_id = ? AND app_name = ?
   		);
	`, userID, app).Scan(&exists)
	if err != nil {
		Logger.Error("Failed to check app existence for user '%s' and app '%s': %v\n", user, app, err)
		return false
	}
	return exists
}

func (u *SqliteRepository) DeleteApp(user string, app string) error {
	userID, err := getUserId(user)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM apps WHERE user_id = ? AND app_name = ?`, userID, app)
	if err != nil {
		return Logger.LogAndReturnError("Failed to delete app '%s' of user '%s', error: %v", app, user, err)
	}
	return nil
}

type AppInfo struct {
	User string
	App  string
}

func (u *SqliteRepository) FindApps(query string) ([]AppInfo, error) {
	var apps []AppInfo

	rows, err := db.Query(`
		SELECT u.user_name, a.app_name 
		FROM users u 
		JOIN apps a ON u.user_id = a.user_id
		WHERE u.user_name LIKE ? OR a.app_name LIKE ?
		LIMIT 100
	`, "%"+query+"%", "%"+query+"%")

	if err != nil {
		return nil, Logger.LogAndReturnError("Failed to fetch apps: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var app AppInfo
		err := rows.Scan(&app.User, &app.App)
		if err != nil {
			Logger.Error("Error scanning app row: %v\n", err)
			continue
		}
		apps = append(apps, app)
	}
	err = rows.Err()
	if err != nil {
		return nil, Logger.LogAndReturnError("Error iterating over app rows: %v\n", err)
	}
	return apps, nil
}

func (u *SqliteRepository) SetCookie(user string, cookie string, expirationDate time.Time) error {
	_, err := db.Exec("UPDATE users SET cookie = ?, expiration_date = ? WHERE user_name = ?", cookie, expirationDate.Format(time.RFC3339), user)
	if err != nil {
		return Logger.LogAndReturnError("Failed to set cookie: %v", err)
	}
	return nil
}

func (u *SqliteRepository) IsCookieValid(cookie string) bool {
	var expirationDateStr string
	err := db.QueryRow("SELECT expiration_date FROM users WHERE cookie = ?", cookie).Scan(&expirationDateStr)
	if err != nil {
		Logger.Error("Failed to fetch expiration date: %v", err)
		return true
	} else if expirationDateStr == "" {
		return true
	}

	expirationDate, err := time.Parse(time.RFC3339, expirationDateStr)
	if err != nil {
		Logger.Error("Failed to parse expiration date: %v\n", err)
		return true
	}

	return time.Now().UTC().After(expirationDate)
}

func (u *SqliteRepository) GetUserWithCookie(cookie string) (string, error) {
	if cookie == "" {
		return "", Logger.LogAndReturnError("Can't search for empty string cookies")
	}

	var user string
	err := db.QueryRow("SELECT user_name FROM users WHERE cookie = ?", cookie).Scan(&user)
	if err != nil {
		return "", Logger.LogAndReturnError("Failed to fetch user with cookie because: %v", err)
	}

	return user, nil
}

// TODO Avoid duplication of "getIdOf(user/app) logic."
func (u *SqliteRepository) CreateTag(user string, app string, tag string) error {
	appID, err := getAppIdFromUsername(user, app)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO tags (app_id, tag_name) VALUES (?, ?)", appID, tag)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

func (u *SqliteRepository) DeleteTag(user string, app string, tag string) error {
	appID, err := getAppIdFromUsername(user, app)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM tags WHERE app_id = ? AND tag_name = ?", appID, tag)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func getAppIdFromUsername(user string, app string) (int, error) {
	userID, err := getUserId(user)
	if err != nil {
		return 0, err
	}

	appID, err := getAppId(userID, app)
	if err != nil {
		return 0, err
	}
	return appID, nil
}

func (u *SqliteRepository) GetTagList(user string, app string) ([]string, error) {
	appID, err := getAppIdFromUsername(user, app)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT tag_name FROM tags WHERE app_id = ?", appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tags, nil
}

// TODO
func (u *SqliteRepository) ChangePassword(user string, newPassword string) error {
	return nil
}

// TODO
func (u *SqliteRepository) ChangeOrigin(user string, newOrigin string) error {
	return nil
}

func getAppId(userID int, app string) (int, error) {
	var appID int
	err := db.QueryRow("SELECT app_id FROM apps WHERE user_id = ? AND app_name = ?", userID, app).Scan(&appID)
	if err != nil {
		return 0, fmt.Errorf("app not found: %w", err)
	}
	return appID, nil
}

func getUserId(user string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE user_name = ?", user).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}
	return userID, nil
}
