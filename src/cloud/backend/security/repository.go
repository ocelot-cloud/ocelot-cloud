package security

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ocelot-cloud/shared"
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

	// TODO EnsureSchemaVersionTable()

	initializeTables()
	Logger.Info("Database initialized")
}

// TODO Initial design for the data model, during implementation/testing you should check if correct.
func initializeTables() {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		user_name VARCHAR(255) NOT NULL,
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
		app_name VARCHAR(255) NOT NULL
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
		content_blob BYTEA NOT NULL
	);
`)
	if err != nil {
		Logger.Fatal("Failed to create tags table: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS groups (
		group_id SERIAL PRIMARY KEY,
		group_name VARCHAR(255) NOT NULL
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
	// Login/Logout/Register
	CreateUser(user string, password string) error
	IsPasswordCorrect(user string, password string) bool
	Logout(user string) error
	ChangePassword(user string, newPassword string) error
	DeleteUser(user string) error

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
	WipeDatabase()
}

var databaseFile = shared.DataDir + "/sqlite.db"
