package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

var samplePassword = "mypassword"
var userManager UserManager

// TODO Finalize.
func TestStuff(t *testing.T) {
	initializeDatabase()
	userManager = &UserManagerSqlite{}
	defer resetDatabase(t)
	assert.False(t, userManager.DoesUserExist(sampleUser))
	err := userManager.CreateRepoUser(sampleUser, samplePassword)
	assert.Nil(t, err)
	assert.True(t, userManager.DoesUserExist(sampleUser))
}

func resetDatabase(t *testing.T) {
	err := deleteIfExist(databaseFile)
	if err != nil {
		Logger.Error("Failed to delete database: %s, error: %v", databaseFile, err)
		t.Fail()
	}
	initializeDatabase()
}
