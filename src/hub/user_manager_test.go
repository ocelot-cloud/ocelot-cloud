package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"testing"
)

var samplePassword = "mypassword"
var userManager UserManager

func init() {
	userManager = &UserManagerSqlite{}
}

// TODO Finalize.
func TestStuff(t *testing.T) {
	defer resetDatabase(t)
	assert.False(t, userManager.DoesUserExist(sampleUser))
	err := userManager.CreateRepoUser(sampleUser, samplePassword)
	assert.Nil(t, err)
	// TODO assert.True(t, a.DoesUserExist(sampleUser))
}

func resetDatabase(t *testing.T) {
	if err := os.Remove(databaseFile); err != nil && !os.IsNotExist(err) {
		assert.Fail(t, err.Error())
	}
	initializeDatabase()
}
