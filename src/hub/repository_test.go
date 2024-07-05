//go:build unit

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"sort"
	"testing"
	"time"
)

var samplePassword = "mypassword"
var um Repository = &SqliteRepository{}

func TestMain(m *testing.M) {
	initializeDatabase(":memory:")
	code := m.Run()
	os.Exit(code)
}

func TestCreateRepoUser(t *testing.T) {
	defer cleanupDatabase()
	assert.False(t, um.DoesUserExist(sampleUser))
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.True(t, um.DoesUserExist(sampleUser))

	assert.Nil(t, um.DeleteUser(sampleUser))
	assert.False(t, um.DoesUserExist(sampleUser))
}

func TestCreateRepoApp(t *testing.T) {
	defer cleanupDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestDeleteAppCascadingThroughUser(t *testing.T) {
	defer cleanupDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.DeleteApp(sampleUser, sampleApp))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestDeleteAppDirectly(t *testing.T) {
	defer cleanupDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.DeleteUser(sampleUser))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestCantCreateUserTwice(t *testing.T) {
	defer cleanupDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.NotNil(t, um.CreateUser(sampleUser, samplePassword))
}

func TestCantCreateAppTwiceForSameUser(t *testing.T) {
	defer cleanupDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.NotNil(t, um.CreateApp(sampleUser, sampleApp))
}

func TestCantCreateAppWithoutUser(t *testing.T) {
	defer cleanupDatabase()
	assert.NotNil(t, um.CreateApp(sampleUser, sampleApp))
}

func TestTolerateSamePasswordForTwoUsers(t *testing.T) {
	defer cleanupDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateUser(user2, samplePassword))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.True(t, um.IsPasswordCorrect(user2, samplePassword))
}

func TestTolerateSameAppsForTwoUsers(t *testing.T) {
	defer cleanupDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.Nil(t, um.CreateUser(user2, samplePassword))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(user2, sampleApp))

	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(user2, sampleApp))

	assert.Nil(t, um.DeleteApp(sampleUser, sampleApp))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(user2, sampleApp))
}

func TestPasswordVerification(t *testing.T) {
	defer cleanupDatabase()
	assert.Nil(t, um.CreateUser(sampleUser, samplePassword))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, um.IsPasswordCorrect(sampleUser, samplePassword+"x"))
}

func TestSearch(t *testing.T) {
	defer cleanupDatabase()
	um.CreateUser(sampleUser, samplePassword)
	app1 := "prefix_myapp_suffix"
	app2 := "prefix_another-app_suffix"
	um.CreateApp(sampleUser, app1)
	um.CreateApp(sampleUser, app2)

	a, err := um.FindApps("app")
	assert.Nil(t, err)

	sort.Slice(a, func(i, j int) bool {
		return a[i].AppName < a[j].AppName
	})

	assert.Equal(t, 2, len(a))
	assert.Equal(t, sampleUser, a[0].Username)
	assert.Equal(t, sampleUser, a[1].Username)
	assert.Equal(t, app2, a[0].AppName)
	assert.Equal(t, app1, a[1].AppName)
}

func TestSearchNegative(t *testing.T) {
	defer cleanupDatabase()
	um.CreateUser(sampleUser, samplePassword)
	app := "prefix_myapp_suffix"
	um.CreateApp(sampleUser, app)

	a, err := um.FindApps("some")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(a))
}

func cleanupDatabase() {
	users := getUsers()
	for _, v := range users {
		um.DeleteUser(v)
	}
}

func getUsers() []string {
	rows, _ := db.Query("SELECT user_name FROM users")
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var userName string
		rows.Scan(&userName)
		usernames = append(usernames, userName)
	}
	return usernames
}

func TestCookieExpiration(t *testing.T) {
	defer cleanupDatabase()
	um.CreateUser(sampleUser, samplePassword)

	assert.True(t, um.IsCookieValid("non-existing-cookie"))

	timeIn30Days := getTimeIn30Days()
	cookie, _ := generateCookie()
	assert.Nil(t, um.SetCookie(sampleUser, cookie.Value, timeIn30Days))
	assert.False(t, um.IsCookieValid(cookie.Value))

	past := time.Now().Add(-1 * time.Second)
	assert.Nil(t, um.SetCookie(sampleUser, cookie.Value, past))
	assert.True(t, um.IsCookieValid(cookie.Value))
}
