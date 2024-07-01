package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// TODO use sqlite to store users, apps (for search), password etc.

var databaseFile = "sqlite.db"
var db *sql.DB

func init() {
	initializeDatabase()
}

func initializeDatabase() {
	// TODO Add database scheme version?
	var err error
	db, err = sql.Open("sqlite3", databaseFile)
	if err != nil {
		Logger.Fatal("Failed to open database: %v\n", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
    		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    		user_name TEXT UNIQUE,
    		hashed_password TEXT
		)
	`)
	if err != nil {
		Logger.Fatal("Failed to create users table: %v\n", err)
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

	// TODO Add initial schemes. With version number table.
}

type UserManager interface {
	CreateRepoUser(user string, password string) error
	DoesUserExist(user string) bool
	DeleteRepoUser(user string) error
	IsPasswordCorrect(user string, password string) bool
	DoesAppExist(user string, app string) bool
	AddApp(user string, app string) error
	DeleteApp(user string, app string) error
}

type UserManagerSqlite struct{}

func (u *UserManagerSqlite) IsPasswordCorrect(user string, password string) bool {
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

func (u *UserManagerSqlite) DoesUserExist(user string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_name = ?)", user).Scan(&exists)
	if err != nil {
		Logger.Error("Failed to check user existence: %v\n", err)
		return false
	}
	return exists
}

func (u *UserManagerSqlite) CreateRepoUser(user string, password string) error {
	hashedPassword, err := hashAndSaltPassword(password)
	if err != nil {
		return Logger.LogAndReturnError("Failed to hash password: %v\n", err)
	}

	_, err = db.Exec("INSERT INTO users (user_name, hashed_password) VALUES (?, ?)", user, hashedPassword)
	if err != nil {
		return Logger.LogAndReturnError("Failed to create user: %v\n", err)
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

func (u *UserManagerSqlite) DeleteRepoUser(user string) error {
	if !u.DoesUserExist(user) {
		return Logger.LogAndReturnError("User %s does not exist", user)
	}

	_, err := db.Exec("DELETE FROM users WHERE user_name = ?", user)
	if err != nil {
		return Logger.LogAndReturnError("Failed to delete user: %v", err)
	}

	return nil
}

func (u *UserManagerSqlite) AddApp(user string, app string) error {
	if !u.DoesUserExist(user) {
		return Logger.LogAndReturnError("User '%s' does not exist", user)
	}

	if u.DoesAppExist(user, app) {
		return Logger.LogAndReturnError("App '%s' already exists for user '%s'", app, user)
	}

	_, err := db.Exec(`
		INSERT INTO apps (user_id, app_name)
		VALUES ((SELECT user_id FROM users WHERE user_name = ?), ?)
	`, user, app)
	if err != nil {
		return Logger.LogAndReturnError("Failed to add app '%s' for user '%s': %v", app, user, err)
	}

	return nil
}

func (u *UserManagerSqlite) DoesAppExist(user string, app string) bool {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM apps WHERE user_id = (
				SELECT user_id FROM users WHERE user_name = ?
		  	) AND app_name = ?
   		);
	`, user, app).Scan(&exists)
	if err != nil {
		Logger.Error("Failed to check app existence for user '%s' and app '%s': %v\n", user, app, err)
		return false
	}
	return exists
}

func (u *UserManagerSqlite) DeleteApp(user string, app string) error {
	_, err := db.Exec(`DELETE FROM apps WHERE user_id = (SELECT user_id FROM users WHERE user_name = ?) AND app_name = ?`, user, app)
	if err != nil {
		return Logger.LogAndReturnError("Failed to delete app '%s' of user '%s', error: %v", app, user, err)
	}
	return nil
}
