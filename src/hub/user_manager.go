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

	// TODO Handle error
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		hashed_password TEXT
	)`)
	// TODO Add initial schemes. With version number table.
}

type UserManager interface {
	CreateRepoUser(user string, password string) error
	DoesUserExist(user string) bool
}

type UserManagerSqlite struct{}

func (u *UserManagerSqlite) DoesUserExist(user string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user).Scan(&exists)
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

	// Insert user into database
	_, err = db.Exec("INSERT INTO users (username, hashed_password) VALUES (?, ?)", user, hashedPassword)
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
