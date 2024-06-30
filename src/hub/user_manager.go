package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// TODO use sqlite to store users, apps (for search), password etc.

var databaseFile = "sqlite.db"
var db *sql.DB

func init() {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		Logger.Fatal("Failed to open database: %v", err)
		return
	}
	defer db.Close()
}

type UserManager interface {
	CreateRepoUser(user string, password string) error
	DoesUserExist(user string) bool
}

type UserManagerSqlite struct{}

func (u *UserManagerSqlite) DoesUserExist(user string) bool {
	return false
}

func (u *UserManagerSqlite) CreateRepoUser(user string, password string) error {
	return nil
}
