package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ocelot-cloud/shared/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var db *sql.DB

func initializeDatabaseWithSource(dataSourceName string) {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		Logger.Fatal("Failed to open database: %v\n", err)
	}

	// Prevents concurrency problems. My guess is that the sqlite client does not handle concurrency correctly.
	// Approach works, but may be too slow. Another client or DB may be needed in the future.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	EnsureSchemaVersionTable()

	// TODO Store only hashed cookies. Should also be UNIQUE
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
    		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_name TEXT UNIQUE NOT NULL,
			hashed_password TEXT NOT NULL UNIQUE,
			origin TEXT,
			cookie TEXT,
			expiration_date TEXT,
		    used_space BIGINT NOT NULL
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
			app_id INTEGER NOT NULL,
			tag_name TEXT NOT NULL,
			data BLOB NOT NULL,
			UNIQUE(app_id, tag_id),
			FOREIGN KEY (app_id) REFERENCES apps(app_id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create tags table: %v\n", err)
	}

	Logger.Info("Database initialized")
}

type Repository interface {
	CreateUser(form *RegistrationForm) error
	DoesUserExist(user string) bool
	DeleteUser(user string) error
	IsPasswordCorrect(user string, password string) bool
	DoesAppExist(user string, app string) bool
	CreateApp(user string, app string) error
	DeleteApp(user string, app string) error
	FindApps(query string) ([]UserAndApp, error)
	SetCookie(user string, cookie string, expirationDate time.Time) error
	IsCookieExpired(cookie string) bool
	GetUserWithCookie(cookie string) (string, error)
	CreateTag(user string, app string, tag string, data []byte) error
	DeleteTag(user string, app string, tag string) error
	GetTagList(user string, app string) ([]string, error)
	ChangePassword(user string, newPassword string) error
	SetOrigin(user string, newOrigin string) error
	IsOriginCorrect(user string, origin string) bool
	DoesTagExist(user string, app string, tag string) bool
	GetTagContent(user string, app string, tag string) ([]byte, error)
	GetUsedSpaceInBytes(user string) (int, error)
	Logout(user string) error
	GetAppList(user string) ([]string, error)
	WipeDatabase()
}

type SqliteRepository struct{}

func (u *SqliteRepository) GetTagContent(user string, app string, tag string) ([]byte, error) {
	appID, err := getAppIdFromUsername(user, app)
	if err != nil {
		return nil, err
	}

	var data []byte
	err = db.QueryRow("SELECT data FROM tags WHERE app_id = ? and tag_name = ?", appID, tag).Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

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
	hashedPassword, err := utils.SaltAndHash(form.Password)
	if err != nil {
		return logAndReturnError("Failed to hash password: %v\n", err)
	}

	_, err = db.Exec("INSERT INTO users (user_name, hashed_password, used_space) VALUES (?, ?, ?)", form.User, hashedPassword, 0)
	if err != nil {
		return logAndReturnError("Failed to create user: %v", err)
	}
	return nil
}

func (u *SqliteRepository) DeleteUser(user string) error {
	if !u.DoesUserExist(user) {
		return logAndReturnError("User %s does not exist", user)
	}

	_, err := db.Exec("DELETE FROM users WHERE user_name = ?", user)
	if err != nil {
		return logAndReturnError("Failed to delete user: %v", err)
	}

	return nil
}

func (u *SqliteRepository) CreateApp(user string, app string) error {
	if !u.DoesUserExist(user) {
		return logAndReturnError("User '%s' does not exist", user)
	} else if u.DoesAppExist(user, app) {
		return logAndReturnError("App '%s' already exists for user '%s'", app, user)
	}

	userID, err := getUserId(user)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO apps (user_id, app_name) VALUES (?, ?)`, userID, app)
	if err != nil {
		return logAndReturnError("Failed to add app '%s' for user '%s': %v", app, user, err)
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

	appID, err := getAppId(userID, app)
	if err != nil {
		return err
	}

	totalDataSize, err := u.sumBlobSizes(appID)
	if err != nil {
		return err
	}

	_, err = db.Exec(`DELETE FROM apps WHERE user_id = ? AND app_name = ?`, userID, app)
	if err != nil {
		return logAndReturnError("Failed to delete app '%s' of user '%s', error: %v", app, user, err)
	}

	_, err = db.Exec("UPDATE users SET used_space = used_space - ? WHERE user_name = ?", totalDataSize, user)
	if err != nil {
		return fmt.Errorf("failed to update user space: %w", err)
	}

	return nil
}

func (u *SqliteRepository) sumBlobSizes(appID int) (int64, error) {
	var totalSize sql.NullInt64
	err := db.QueryRow("SELECT SUM(LENGTH(data)) FROM tags WHERE app_id = ?", appID).Scan(&totalSize)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total BLOB size: %w", err)
	}

	if !totalSize.Valid {
		return 0, nil
	}

	return totalSize.Int64, nil
}

func (u *SqliteRepository) FindApps(query string) ([]UserAndApp, error) {
	var apps []UserAndApp

	rows, err := db.Query(`
		SELECT u.user_name, a.app_name 
		FROM users u 
		JOIN apps a ON u.user_id = a.user_id
		WHERE u.user_name LIKE ? OR a.app_name LIKE ?
		LIMIT 100
	`, "%"+query+"%", "%"+query+"%")

	if err != nil {
		return nil, logAndReturnError("Failed to fetch apps: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var app UserAndApp
		err := rows.Scan(&app.User, &app.App)
		if err != nil {
			Logger.Error("Error scanning app row: %v\n", err)
			continue
		}
		apps = append(apps, app)
	}
	err = rows.Err()
	if err != nil {
		return nil, logAndReturnError("Error iterating over app rows: %v\n", err)
	}
	return apps, nil
}

func (u *SqliteRepository) SetCookie(user string, cookie string, expirationDate time.Time) error {
	_, err := db.Exec("UPDATE users SET cookie = ?, expiration_date = ? WHERE user_name = ?", cookie, expirationDate.Format(time.RFC3339), user)
	if err != nil {
		return logAndReturnError("Failed to set cookie: %v", err)
	}
	return nil
}

func (u *SqliteRepository) IsCookieExpired(cookie string) bool {
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
		return "", logAndReturnError("Can't search for empty string cookies")
	}

	var user string
	err := db.QueryRow("SELECT user_name FROM users WHERE cookie = ?", cookie).Scan(&user)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", logAndReturnInfo("cookie not found")
		} else {
			return "", logAndReturnError("Failed to fetch user with cookie: %v", err)
		}
	}

	return user, nil
}

func (u *SqliteRepository) CreateTag(user string, app string, tag string, data []byte) error {
	appID, err := getAppIdFromUsername(user, app)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO tags (app_id, tag_name, data) VALUES (?, ?, ?)", appID, tag, data)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	dataSize := len(data)
	_, err = db.Exec("UPDATE users SET used_space = used_space + ? WHERE user_name = ?", dataSize, user)
	if err != nil {
		return fmt.Errorf("failed to update user space: %w", err)
	}

	return nil
}

func (u *SqliteRepository) DeleteTag(user string, app string, tag string) error {
	appID, err := getAppIdFromUsername(user, app)
	if err != nil {
		return err
	}

	dataSize, err := getBlobSize(appID, tag)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM tags WHERE app_id = ? AND tag_name = ?", appID, tag)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	_, err = db.Exec("UPDATE users SET used_space = used_space - ? WHERE user_name = ?", dataSize, user)
	if err != nil {
		return fmt.Errorf("failed to update user space: %w", err)
	}

	return nil
}

func getBlobSize(appID int, tag string) (int64, error) {
	var dataSize int64
	err := db.QueryRow("SELECT LENGTH(data) FROM tags WHERE app_id = ? AND tag_name = ?", appID, tag).Scan(&dataSize)
	if err != nil {
		return 0, fmt.Errorf("failed to get BLOB size: %w", err)
	}
	return dataSize, nil
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

func (u *SqliteRepository) GetAppList(user string) ([]string, error) {
	userID, err := getUserId(user)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT app_name FROM apps WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get apps: %w", err)
	}
	defer rows.Close()

	var apps []string
	for rows.Next() {
		var app string
		if err = rows.Scan(&app); err != nil {
			return nil, fmt.Errorf("failed to scan app: %w", err)
		}
		apps = append(apps, app)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return apps, nil
}

func (u *SqliteRepository) ChangePassword(user string, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return logAndReturnError("failed to hash password: %w", err)
	}

	_, err = db.Exec("UPDATE users SET hashed_password = ? WHERE user_name = ?", hashedPassword, user)
	if err != nil {
		return logAndReturnError("failed to update password: %w", err)
	}

	return nil
}

func (u *SqliteRepository) SetOrigin(user string, newOrigin string) error {
	_, err := db.Exec("UPDATE users SET origin = ? WHERE user_name = ?", newOrigin, user)
	if err != nil {
		return logAndReturnError("failed to update origin: %w", err)
	}
	return nil
}

func (u *SqliteRepository) IsOriginCorrect(user string, origin string) bool {
	var repoOrigin sql.NullString
	err := db.QueryRow("SELECT origin FROM users WHERE user_name = ?", user).Scan(&repoOrigin)
	if err != nil {
		Logger.Error("Failed to fetch origin: %v\n", err)
		return false
	}

	if repoOrigin.Valid {
		return repoOrigin.String == origin
	} else {
		return false
	}
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

func (u *SqliteRepository) WipeDatabase() {
	users := getUsers()
	for _, v := range users {
		if v != "sample" {
			u.DeleteUser(v)
		}
	}
}

func getUsers() []string {
	rows, _ := db.Query("SELECT user_name FROM users")
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var userName string
		rows.Scan(&userName)
		usernames = append(usernames, userName)
	}
	return usernames
}

func (u *SqliteRepository) DoesTagExist(user string, app string, tag string) bool {
	appID, err := getAppIdFromUsername(user, app)
	if err != nil {
		return false
	}

	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM tags WHERE app_id = ? AND tag_name = ?
		);
	`, appID, tag).Scan(&exists)
	if err != nil {
		Logger.Debug("error checking if tag exists")
		return false
	}
	return exists
}

func (u *SqliteRepository) GetUsedSpaceInBytes(user string) (int, error) {
	var usedSpace int
	err := db.QueryRow(`SELECT used_space FROM users WHERE user_name = ?`, user).Scan(&usedSpace)
	if err != nil {
		return 0, logAndReturnError("failed to get current space: %w", err)
	}
	return usedSpace, nil
}

func (u *SqliteRepository) Logout(user string) error {
	_, err := db.Exec("UPDATE users SET cookie = ?, expiration_date = ? WHERE user_name = ?", nil, nil, user)
	if err != nil {
		return logAndReturnError("Failed to logout: %v", err)
	}
	return nil
}

// TODO I think I should get rid of these two functions. I dont want low level errors to be transported to the top, or even displayed to users.
func logAndReturnError(message string, args ...interface{}) error {
	Logger.Error(message, args...)
	return fmt.Errorf(message, args...)
}

func logAndReturnInfo(message string, args ...interface{}) error {
	Logger.Info(message, args...)
	return fmt.Errorf(message, args...)
}
