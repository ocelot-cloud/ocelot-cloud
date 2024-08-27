package security

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ocelot-cloud/shared"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var databaseFile = shared.DataDir + "/sqlite.db"

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

	// TODO EnsureSchemaVersionTable()

	initializeTables()
	Logger.Info("Database initialized")
}

// TODO Initial design for the data model, during implementation/testing you should check if correct.
func initializeTables() {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		user_name VARCHAR(255) NOT NULL UNIQUE,
		hashed_password VARCHAR(255) NOT NULL,
		hashed_cookie VARCHAR(255),
		is_admin BOOLEAN NOT NULL
	);
`)
	if err != nil {
		Logger.Fatal("Failed to create users table: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS apps (
		app_id SERIAL PRIMARY KEY,
		maintainer_name VARCHAR(255) NOT NULL,
		app_name VARCHAR(255) NOT NULL,
		UNIQUE (maintainer_name, app_name)
	);
`)
	if err != nil {
		Logger.Fatal("Failed to create apps table: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tags (
		app_id INT REFERENCES apps(app_id) ON DELETE CASCADE,
		tag_id SERIAL PRIMARY KEY,
		tag_name VARCHAR(255) NOT NULL,
		content_blob BYTEA NOT NULL,
		UNIQUE (app_id, tag_name)
	);
`)
	if err != nil {
		Logger.Fatal("Failed to create tags table: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS groups (
		group_id SERIAL PRIMARY KEY,
		group_name VARCHAR(255) NOT NULL UNIQUE 
	);
`)
	if err != nil {
		Logger.Fatal("Failed to create groups table: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user_to_group (
		user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
		group_id INT REFERENCES groups(group_id) ON DELETE CASCADE,
		PRIMARY KEY (user_id, group_id)
	);
`)
	if err != nil {
		Logger.Fatal("Failed to create user_to_group table: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS app_access (
		group_id INT REFERENCES groups(group_id) ON DELETE CASCADE,
		app_id INT REFERENCES apps(app_id) ON DELETE CASCADE,
		PRIMARY KEY (group_id, app_id)
	);
`)
	if err != nil {
		Logger.Fatal("Failed to create app_access table: %v", err)
	}
}

type Authorization struct {
	user    string
	isAdmin bool
}

// TODO To be implemented
type Repository interface {
	// User Handling
	CreateUser(user string, password string, isAdmin bool) error
	WipeDatabase()
	IsPasswordCorrect(user string, password string) bool
	DeleteUser(user string) error
	/*
		Logout(user string) error
		ChangePassword(user string, newPassword string) error

		// Auth
		DoesUserExist(user string) bool
		GetUserWithCookie(cookie string) (Authorization, error)


		// TODO Matrial from hub which might be an inspiration. If not used, please delete.
		DoesAppExist(user string, app string) bool
		CreateApp(user string, app string) error
		DeleteApp(user string, app string) error
		FindApps(query string) ([]string, error)
		SetCookie(user string, cookie string, expirationDate time.Time) error
		IsCookieExpired(cookie string) bool

		CreateTag(user string, app string, tag string, data []byte) error
		DeleteTag(user string, app string, tag string) error
		GetTagList(user string, app string) ([]string, error)

		SetOrigin(user string, newOrigin string) error
		IsOriginCorrect(user string, origin string) bool
		DoesTagExist(user string, app string, tag string) bool
		GetTagContent(user string, app string, tag string) ([]byte, error)
		GetUsedSpaceInBytes(user string) (int, error)

		GetAppList(user string) ([]string, error)
	*/
}

type MyRepository struct{}

func (r *MyRepository) CreateUser(user string, password string, isAdmin bool) error {
	hashedPassword, err := hashAndSaltPassword(password)
	if err != nil {
		Logger.Warn("Failed to hash password: %v", err)
		return fmt.Errorf("failed to hash password")
	}
	_, err = db.Exec("INSERT INTO users (user_name, hashed_password, is_admin) VALUES (?, ?, ?)", user, hashedPassword, isAdmin)
	if err != nil {
		Logger.Warn("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user")
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

func (r *MyRepository) WipeDatabase() {

}

func (r *MyRepository) IsPasswordCorrect(user string, password string) bool {
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

func (r *MyRepository) DeleteUser(user string) error {
	_, err := db.Exec("DELETE FROM users WHERE user_name = ?", user)
	if err != nil {
		Logger.Warn("Failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user")
	}
	return nil
}
