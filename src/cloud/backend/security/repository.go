package security

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ocelot-cloud/shared"
	"github.com/ocelot-cloud/shared/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var Repo Repository = &MyRepository{}

var DB *sql.DB
var DatabaseFile = shared.DataDir + "/sqlite.db"

func InitializeDatabaseWithSource(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		Logger.Fatal("Failed to open database: %v\n", err)
	}

	// Prevents concurrency problems. My guess is that the sqlite client does not handle concurrency correctly.
	// Approach works, but may be too slow. Another client or DB may be needed in the future.
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	// TODO EnsureSchemaVersionTable()

	initializeTables()
	Logger.Info("Database initialized")
}

// TODO Initial design for the data model, during implementation/testing you should check if correct.
func initializeTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			user_name VARCHAR(255) NOT NULL UNIQUE,
			hashed_password VARCHAR(255) NOT NULL UNIQUE,
			hashed_cookie_value VARCHAR(255) UNIQUE,
			cookie_expiration_date VARCHAR(255),
			is_admin BOOLEAN NOT NULL
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create users table: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS apps (
			app_id SERIAL PRIMARY KEY,
			maintainer VARCHAR(255) NOT NULL,
			app VARCHAR(255) NOT NULL,
			UNIQUE (maintainer, app)
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create apps table: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			app_id INT REFERENCES apps(app_id) ON DELETE CASCADE,
			tag_id SERIAL PRIMARY KEY,
			tag VARCHAR(255) NOT NULL,
			blob BYTEA NOT NULL,
			UNIQUE (app_id, tag)
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create tags table: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS groups (
			group_id SERIAL PRIMARY KEY,
			group_name VARCHAR(255) NOT NULL UNIQUE 
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create groups table: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS user_to_group (
			user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
			group_id INT REFERENCES groups(group_id) ON DELETE CASCADE,
			PRIMARY KEY (user_id, group_id)
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create user_to_group table: %v", err)
	}

	_, err = DB.Exec(`
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
	User    string
	IsAdmin bool
}

type MaintainerAndApp struct {
	Maintainer string
	App        string
}

type TagAndBlob struct {
	Tag  string
	Blob []byte
}

// TODO To be implemented
type Repository interface {
	CreateUser(user string, password string, isAdmin bool) error
	WipeDatabase()
	IsPasswordCorrect(user string, password string) bool
	DeleteUser(user string) error
	HashAndSaveCookie(user string, cookieValue string, cookieExpirationDate time.Time) error
	Logout(user string) error
	DoesUserExist(user string) bool
	GetUserViaCookie(cookieValue string) (*Authorization, error)
	DoesAnyAdminUserExist() bool

	ChangePassword(user string, newPassword string) error

	CreateAppWithTag(maintainer string, app string, tag string, blob []byte) error
	ListAppInfo() ([]MaintainerAndApp, error)
	ListTagsOfApp(maintainer string, app string) ([]string, error)
	LoadTagBlob(maintainer string, app string, tag string) ([]byte, error)
	/*
		// TODO Material from hub which might be an inspiration. If not used, please delete.
		// TODO Instead of appForm/appEntry, just use app_id, much easier. Add: GetAppId(appForm).
		type appForm struct: {maintainer, app, tag}
		delete app: func(appForm) err
		load app: func(appForm) (blob, error)
		type AppInfo struct : {maintainer, app}

		CreateGroup(group string) error
		DeleteGroup(group) error
		ListUsers() ([]string, error)
		AddUserToGroup(user, group) error
		RemoveUserFromGroup(user, group) error
		GiveGroupAccessToApp(group, app), error
		RemoveGroupsAccessToApp(group, app), error
		DoesUserHaveAccessToApp(user, appForm) bool

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

func (r *MyRepository) DoesAnyAdminUserExist() bool {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE is_admin = ?)", true).Scan(&exists)
	if err != nil {
		Logger.Error("Failed to check if there is any admin user: %v", err)
		return false
	}
	return exists
}

func (r *MyRepository) CreateUser(user string, password string, isAdmin bool) error {
	hashedPassword, err := utils.SaltAndHash(password)
	if err != nil {
		return err
	}
	_, err = DB.Exec("INSERT INTO users (user_name, hashed_password, is_admin) VALUES (?, ?, ?)", user, hashedPassword, isAdmin)
	if err != nil {
		Logger.Warn("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user")
	}
	return nil
}

// TODO shift to shared module

func (r *MyRepository) WipeDatabase() {
	_, err := DB.Exec("DELETE FROM users")
	if err != nil {
		Logger.Fatal("Database wipe failed: %v", err)
	}
}

func (r *MyRepository) IsPasswordCorrect(user string, password string) bool {
	var hashedPassword string
	err := DB.QueryRow("SELECT hashed_password FROM users WHERE user_name = ?", user).Scan(&hashedPassword)
	if err != nil {
		Logger.Error("Failed to fetch hashed password: %v", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (r *MyRepository) DeleteUser(user string) error {
	_, err := DB.Exec("DELETE FROM users WHERE user_name = ?", user)
	if err != nil {
		Logger.Warn("Failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user")
	}
	return nil
}

func (r *MyRepository) HashAndSaveCookie(user string, cookieValue string, cookieExpirationDate time.Time) error {
	hashedCookieValue, err := utils.Hash(cookieValue)
	if err != nil {
		return err
	}

	_, err = DB.Exec("UPDATE users SET hashed_cookie_value = ?, cookie_expiration_date = ? WHERE user_name = ?", hashedCookieValue, cookieExpirationDate.Format(time.RFC3339), user)
	if err != nil {
		Logger.Warn("Failed to update cookie of user '%s': %v", user, err)
		return fmt.Errorf("failed to update cookie")
	}
	return nil
}

func (r *MyRepository) Logout(user string) error {
	_, err := DB.Exec("UPDATE users SET hashed_cookie_value = ?, cookie_expiration_date = ? WHERE user_name = ?", "", "", user)
	if err != nil {
		Logger.Error("Failed to delete cookie of user '%s': %v", user, err)
		return fmt.Errorf("failed to delete cookie")
	}
	return nil
}

func (r *MyRepository) DoesUserExist(user string) bool {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_name = ?)", user).Scan(&exists)
	if err != nil {
		Logger.Error("Failed to check if user exists: %v", err)
		return false
	}
	return exists
}

// TODO Test if isAdmin is correct in authorization.
func (r *MyRepository) GetUserViaCookie(cookieValue string) (*Authorization, error) {
	hashedCookieValue, err := utils.Hash(cookieValue)
	if err != nil {
		return nil, err
	}

	var user string
	var isAdmin bool
	err = DB.QueryRow("SELECT user_name, is_admin FROM users WHERE hashed_cookie_value = ?", hashedCookieValue).Scan(&user, &isAdmin)
	if err != nil {
		Logger.Error("Failed to fetch user data: %v", err)
		return nil, fmt.Errorf("failed to fetch user data")
	}
	return &Authorization{user, isAdmin}, nil
}

func (r *MyRepository) ChangePassword(user string, newPassword string) error {
	hashedNewPassword, err := utils.SaltAndHash(newPassword)
	if err != nil {
		return err
	}

	_, err = DB.Exec("UPDATE users SET hashed_password = ? WHERE user_name = ?", hashedNewPassword, user)
	if err != nil {
		Logger.Error("Failed to update password of user '%s': %v", user, err)
		return fmt.Errorf("failed to update password")
	}
	return nil
}

// TODO is app already exists,
func (r *MyRepository) CreateAppWithTag(maintainer string, app string, tag string, blob []byte) error {
	_, err := DB.Exec("INSERT INTO apps (maintainer, app) VALUES (?, ?)", maintainer, app)
	if err != nil {
		Logger.Error("Failed to create app: %v", err)
		return fmt.Errorf("failed to create app")
	}

	appId, err := r.getAppId(maintainer, app)
	if err != nil {
		// TODO
	}

	_, err = DB.Exec("INSERT INTO tags (app_id, tag, blob) VALUES (?, ?, ?)", appId, tag, blob)
	if err != nil {
		// TODO
	}

	return nil
}

func (r *MyRepository) getAppId(maintainer string, app string) (int, error) {
	var appId int
	err := DB.QueryRow("SELECT app_id FROM apps WHERE maintainer = ? AND app = ?", maintainer, app).Scan(&appId)
	if err != nil {
		// TODO
		return 0, err
	}
	return appId, nil
}

func (r *MyRepository) ListAppInfo() ([]MaintainerAndApp, error) {
	rows, err := DB.Query("SELECT maintainer, app FROM apps")
	if err != nil {
		Logger.Error("Failed to fetch app list: %v", err)
		return nil, fmt.Errorf("failed to fetch app list")
	}
	defer rows.Close()

	var result []MaintainerAndApp
	for rows.Next() {
		var maintainer, app string
		if err = rows.Scan(&maintainer, &app); err != nil {
			Logger.Error("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row")
		}
		result = append(result, MaintainerAndApp{Maintainer: maintainer, App: app})
	}

	if err = rows.Err(); err != nil {
		Logger.Error("Rows error: %v", err)
		return nil, fmt.Errorf("rows error")
	}

	return result, nil
}

func (r *MyRepository) ListTagsOfApp(maintainer string, app string) ([]string, error) {
	// TODO duplication with ListAppInfo
	appId, err := r.getAppId(maintainer, app)
	if err != nil {
		// TODO
	}

	rows, err := DB.Query("SELECT tag FROM tags WHERE app_id = ?", appId)
	if err != nil {
		Logger.Error("Failed to fetch tag list of app: %s/%s, %v", maintainer, app, err)
		return nil, fmt.Errorf("failed to fetch tag list")
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var singleTag string
		if err = rows.Scan(&singleTag); err != nil {
			Logger.Error("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row")
		}
		result = append(result, singleTag)
	}

	if err = rows.Err(); err != nil {
		Logger.Error("Rows error: %v", err)
		return nil, fmt.Errorf("rows error")
	}

	return result, nil
}

func (r *MyRepository) LoadTagBlob(maintainer string, app string, tag string) ([]byte, error) {
	appId, err := r.getAppId(maintainer, app)
	if err != nil {
		// TODO
	}

	var blob []byte
	err = DB.QueryRow("SELECT blob FROM tags WHERE app_id = ? AND tag = ?", appId, tag).Scan(&blob)
	if err != nil {
		// TODO
	}

	return blob, nil
}

// TODO for the handlers: admins should be able to delete an account. But should users be able to delete their own account? I think not. This can cause many troubles if a user does it accidentally. Maybe a feature that is disabled by default, but which can be enabled manually.
// TODO idea: by default create a group "anonymous" which cant be deleted. Access to an app for members of anonymous means, that any user, even without account can access an app.
// TODO if an app is deleted, all its tags must be deleted. If all tags of an app are deleted, the app must be deleted as well.
