package main

import (
	"database/sql"
	"errors"
	"log"
)

func EnsureSchemaVersionTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create database_schema_version table: %v", err)
	}

	var version string
	err = db.QueryRow("SELECT version FROM schema_version").Scan(&version)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec("INSERT INTO schema_version (version) VALUES (?)", currentSchemaVersion)
		if err != nil {
			log.Fatalf("Failed to set initial schema version: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Failed to query schema version: %v", err)
	} else if version != currentSchemaVersion {
		log.Fatalf("Schema version mismatch: expected %s, got %s", currentSchemaVersion, version)
	}
}
