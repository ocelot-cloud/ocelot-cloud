package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"ocelot/tools"
)

var Logger = tools.ProvideLogger()

func DoSomeDataBaseStuff() {
	db, err := sql.Open("sqlite3", "ocelot-cloud.db")
	if err != nil {
		Logger.Fatal("Failed to open database: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS foo (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT)")
	if err != nil {
		Logger.Debug("Failed to create table: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO foo (name) VALUES (?)", "gopher")
	if err != nil {
		Logger.Debug("Failed to insert into table: %v", err)
		return
	}

	rows, err := db.Query("SELECT id, name FROM foo")
	if err != nil {
		Logger.Debug("Failed to query table: %v", err)
		return
	}
	defer rows.Close()

	Logger.Debug("Current rows in 'foo' table:")
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			Logger.Debug("Failed to scan row: %v", err)
			continue
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		Logger.Debug("Error occurred during row iteration: %v", err)
	}
}
