package main

import (
	"github.com/ocelot-cloud/shared"
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
	shared.AssertFalse(t, userManager.DoesUserExist(sampleUser))
	err := userManager.CreateRepoUser(sampleUser, samplePassword)
	shared.AssertNil(t, err)
	// TODO shared.AssertTrue(t, a.DoesUserExist(sampleUser))
}

func resetDatabase(t *testing.T) {
	if err := os.Remove("sqlite.db"); err != nil && !os.IsNotExist(err) {
		// TODO Can "shared" module provider have an "assert" package? -> assert.Equal() ...
		// TODO Should only be "Fail". And "shared" should be "assert" -> make other function names shorter as well. Also add Skip().
		shared.AssertFail(t, err.Error())
	}
	initializeDatabase()
}
