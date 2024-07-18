//go:build unit

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

var um Repository = &SqliteRepository{}

func TestMain(m *testing.M) {
	initializeDatabaseWithSource(":memory:")
	code := m.Run()
	os.Exit(code)
}

func TestCreateRepoUser(t *testing.T) {
	defer um.WipeDatabase()
	assert.False(t, um.DoesUserExist(sampleUser))
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.True(t, um.DoesUserExist(sampleUser))

	assert.Nil(t, um.DeleteUser(sampleUser))
	assert.False(t, um.DoesUserExist(sampleUser))
}

func TestCreateRepoApp(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestDeleteAppCascadingThroughUser(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.DeleteApp(sampleUser, sampleApp))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestDeleteAppDirectly(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.Nil(t, um.DeleteUser(sampleUser))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
}

func TestCantCreateUserTwice(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.NotNil(t, um.CreateUser(sampleForm))
}

func TestCantCreateAppTwiceForSameUser(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.NotNil(t, um.CreateApp(sampleUser, sampleApp))
}

func TestCantCreateAppWithoutUser(t *testing.T) {
	defer um.WipeDatabase()
	assert.NotNil(t, um.CreateApp(sampleUser, sampleApp))
}

func TestTolerateSamePasswordForTwoUsers(t *testing.T) {
	defer um.WipeDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, um.CreateUser(sampleForm))
	newForm := *sampleForm
	newForm.Username = user2
	assert.Nil(t, um.CreateUser(&newForm))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.True(t, um.IsPasswordCorrect(user2, samplePassword))
}

func TestTolerateSameAppsForTwoUsers(t *testing.T) {
	defer um.WipeDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, um.CreateUser(sampleForm))
	newForm := *sampleForm
	newForm.Username = user2
	assert.Nil(t, um.CreateUser(&newForm))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	assert.Nil(t, um.CreateApp(user2, sampleApp))

	assert.True(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(user2, sampleApp))

	assert.Nil(t, um.DeleteApp(sampleUser, sampleApp))
	assert.False(t, um.DoesAppExist(sampleUser, sampleApp))
	assert.True(t, um.DoesAppExist(user2, sampleApp))
}

func TestPasswordVerification(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, um.IsPasswordCorrect(sampleUser, samplePassword+"x"))
}

func TestSearch(t *testing.T) {
	defer um.WipeDatabase()
	um.CreateUser(sampleForm)
	app1 := "prefix_myapp_suffix"
	app2 := "prefix_another-app_suffix"
	um.CreateApp(sampleUser, app1)
	um.CreateApp(sampleUser, app2)

	a, err := um.FindApps("app")
	assert.Nil(t, err)

	sort.Slice(a, func(i, j int) bool {
		return a[i].App < a[j].App
	})

	assert.Equal(t, 2, len(a))
	assert.Equal(t, sampleUser, a[0].User)
	assert.Equal(t, sampleUser, a[1].User)
	assert.Equal(t, app2, a[0].App)
	assert.Equal(t, app1, a[1].App)
}

func TestSearchNegative(t *testing.T) {
	defer um.WipeDatabase()
	um.CreateUser(sampleForm)
	app := "prefix_myapp_suffix"
	um.CreateApp(sampleUser, app)

	a, err := um.FindApps("some")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(a))
}

// TODO handle the unhandled errors
func TestCookieExpiration(t *testing.T) {
	defer um.WipeDatabase()
	um.CreateUser(sampleForm)
	_, err := um.GetUserWithCookie("")
	assert.NotNil(t, err)

	assert.True(t, um.IsCookieExpired("non-existing-cookie"))

	timeIn30Days := getTimeIn30Days()
	cookie, _ := generateCookie()
	assert.Nil(t, um.SetCookie(sampleUser, cookie.Value, timeIn30Days))
	assert.False(t, um.IsCookieExpired(cookie.Value))

	past := time.Now().Add(-1 * time.Second)
	assert.Nil(t, um.SetCookie(sampleUser, cookie.Value, past))
	assert.True(t, um.IsCookieExpired(cookie.Value))

	user, err := um.GetUserWithCookie(cookie.Value)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, user)
}

func TestGetTagList(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))
	foundTags, err := um.GetTagList(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundTags))
	assert.False(t, um.DoesTagExist(sampleUser, sampleApp, sampleTag))

	assert.Nil(t, um.CreateTag(sampleUser, sampleApp, sampleTag, []byte("asdf")))
	foundTags, err = um.GetTagList(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, sampleTag, foundTags[0])
	assert.True(t, um.DoesTagExist(sampleUser, sampleApp, sampleTag))
	data, err := um.GetTagContent(sampleUser, sampleApp, sampleTag)
	assert.Nil(t, err)
	assert.Equal(t, []byte("asdf"), data)

	assert.Nil(t, um.DeleteTag(sampleUser, sampleApp, sampleTag))
	foundTags, err = um.GetTagList(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundTags))
	assert.False(t, um.DoesTagExist(sampleUser, sampleApp, sampleTag))

	assert.Nil(t, um.CreateTag(sampleUser, sampleApp, sampleTag, []byte("asdf")))
	foundTags, err = um.GetTagList(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, sampleTag, foundTags[0])
	assert.True(t, um.DoesTagExist(sampleUser, sampleApp, sampleTag))
	assert.Nil(t, um.DeleteUser(sampleUser))
	_, err = um.GetTagList(sampleUser, sampleApp)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "user not found"))
	assert.False(t, um.DoesTagExist(sampleUser, sampleApp, sampleTag))
}

func TestChangeRepoPassword(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.True(t, um.IsPasswordCorrect(sampleUser, samplePassword))
	newPassword := samplePassword + "x"
	assert.Nil(t, um.ChangePassword(sampleUser, newPassword))
	assert.False(t, um.IsPasswordCorrect(sampleUser, sampleForm.Password))
	assert.True(t, um.IsPasswordCorrect(sampleUser, newPassword))
}

func TestChangeRepoOrigin(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	assert.True(t, um.IsOriginCorrect(sampleUser, sampleOrigin))
	newOrigin := "http://my-new-domain.com:8080"
	assert.Nil(t, um.ChangeOrigin(sampleUser, newOrigin))
	assert.False(t, um.IsOriginCorrect(sampleUser, sampleForm.Origin))
	assert.True(t, um.IsOriginCorrect(sampleUser, newOrigin))
}

// TODO Input validation for such bytes? -> I should test that when I unpack this tar.gz, there should be at least a docker-compose.yml, optionally a app.yml, nothing else. How to do security checks on that file?
func TestUsedSpace(t *testing.T) {
	defer um.WipeDatabase()
	assert.Nil(t, um.CreateUser(sampleForm))
	space, err := um.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, space)

	assert.Nil(t, um.CreateApp(sampleUser, sampleApp))

	bytes := []byte("hello")
	bytes2 := []byte(" world")
	assert.Nil(t, um.CreateTag(sampleUser, sampleApp, sampleTag, bytes))
	space, err = um.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 5, space)

	assert.Nil(t, um.CreateTag(sampleUser, sampleApp, sampleTag+"x", bytes2))
	space, err = um.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 11, space)

	assert.Nil(t, um.DeleteTag(sampleUser, sampleApp, sampleTag))
	space, err = um.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 6, space)

	assert.Nil(t, um.CreateTag(sampleUser, sampleApp, sampleTag, bytes2))
	space, err = um.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 12, space)

	assert.Nil(t, um.DeleteApp(sampleUser, sampleApp))
	space, err = um.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, space)
}
