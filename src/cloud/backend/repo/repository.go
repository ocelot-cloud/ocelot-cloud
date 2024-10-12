package repo

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ocelot-cloud/shared"
	"ocelot/backend/tools"
	"time"
)

var UserRepo UserRepository = &UserRepositoryImpl{}
var AppRepo AppRepository = &AppRepositoryImpl{}
var GroupRepo GroupRepository = &GroupRepositoryImpl{}
var dbRepo DatabaseRepository = &DatabaseRepositoryImpl{}
var AccessRepo AccessRepository = &AccessRepositoryImpl{}

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

	initializeTables()
	Logger.Info("database initialized")
}

func initializeTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_name VARCHAR(255) NOT NULL UNIQUE,
			hashed_password VARCHAR(255) NOT NULL UNIQUE,
			cookie_value VARCHAR(255) UNIQUE,
			cookie_expiration_date VARCHAR(255),
			is_admin BOOLEAN NOT NULL,
			secret VARCHAR(255) UNIQUE
		);
	`)
	if err != nil {
		Logger.Fatal("Failed to create users table: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS apps (
			app_id INTEGER PRIMARY KEY AUTOINCREMENT,
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
			tag_id INTEGER PRIMARY KEY AUTOINCREMENT,
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
			group_id INTEGER PRIMARY KEY AUTOINCREMENT,
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

type MaintainerAndApp struct {
	Maintainer string
	App        string
}

type TagAndBlob struct {
	Tag  string
	Blob []byte
}

type DatabaseRepository interface {
	WipeDatabase()
	// TODO Add functions like IsTableXEmpty() or so.
}

// TODO Check if all methods listed here were already tested.
// TODO Put all data strcutures together and check if they could be simplified
type UserRepository interface {
	CreateUser(user, password string, isAdmin bool) error
	IsPasswordCorrect(user, password string) bool
	DeleteUser(user string) error
	SaveCookie(user, cookieValue string, cookieExpirationDate time.Time) error
	Logout(user string) error
	DoesUserExist(user string) bool
	GetUserViaCookie(cookieValue string) (*tools.Authorization, error)
	DoesAnyAdminUserExist() bool
	ChangePassword(user, newPassword string) error

	GenerateSecret(user string) (string, error)
	GetAssociatedCookieValueAndDeleteSecret(secret string) (string, error)
}

// TODO Put ID's first in the structs.
type App struct {
	Maintainer string
	Name       string
	Id         int
}

type Tag struct {
	Name string
	Id   int
}

type AppRepository interface {
	CreateAppWithTag(maintainer, app, tag string, blob []byte) error // TODO Maybe these args should be a single data structure?
	ListApps() ([]App, error)
	ListTagsOfApp(appId int) ([]Tag, error)
	LoadTagBlob(appId int) ([]byte, error)
	DeleteApp(appId int) error
	DeleteTag(maintainer, app, tag string) error
	GetAppId(maintainer, app string) (int, error)
	GetTagId(appId int, tag string) (int, error)
	DoesTagExist(tagInfo tools.TagInfo) bool
}

type GroupRepository interface {
	CreateGroup(group string) error
	ListGroups() ([]string, error)
	DeleteGroup(group string) error
	GetGroupId(group string) (int, error)

	ListAllUsers() ([]string, error)
	AddUserToGroup(user, group string) error
	ListMembersOfGroup(group string) ([]string, error)
	RemoveUserFromGroup(user, group string) error
}

type AccessRepository interface {
	DoesUserHaveAccessToApp(user string, app MaintainerAndApp) bool
	GiveGroupAccessToApp(group string, app MaintainerAndApp) error
	ListAppAccessesOfGroup(group string) ([]MaintainerAndApp, error)
	RemoveGroupsAccessToApp(group string, app MaintainerAndApp) error
}

type DatabaseRepositoryImpl struct{}
type UserRepositoryImpl struct{}
type AppRepositoryImpl struct{}
type GroupRepositoryImpl struct{}
type AccessRepositoryImpl struct{}

// TODO for the handlers: admins should be able to delete an account. But should users be able to delete their own account? I think not. This can cause many troubles if a user does it accidentally. Maybe a feature that is disabled by default, but which can be enabled manually.
// TODO idea: by default create a group "anonymous" which cant be deleted. Access to an app for members of anonymous means, that any user, even without account can access an app.
// TODO if an app is deleted, all its tags must be deleted. If all tags of an app are deleted, the app must be deleted as well.
// TODO in hub, check if I consistently use: "INTEGER PRIMARY KEY AUTOINCREMENT" for the ID's. If not, apply it.
// TODO Delete duplicated argument types in functions
// TODO Also add deletion tests. For example, when deleting OR app, in both cases tha group to app relation must be deleted as well. maybe: assert.True(t, isTableEmpty("app-to-group")) or so.
// TODO Also check stuff like: user has access to app via group "x". Delete group so that user loses access and re-create with same name. Assert that user has no longer access.
// TODO function: List apps that can be accessed by user.
