package security

import (
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"testing"
	"time"
)

var repo Repository = &MyRepository{}

func TestMain(m *testing.M) {
	InitializeDatabaseWithSource(":memory:")
	repo.WipeDatabase()
	code := m.Run()
	os.Exit(code)
}

var (
	sampleUser     = "user"
	samplePassword = "password"
	sampleCookie   = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	sampleMaintainer = "maintainer"
	sampleApp        = "app"
	sampleTag        = "1.0"
	sampleBlob       = []byte("hello")
)

// TODO Finish SQLite Client Implementation And Tests
func TestSqliteClient(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, repo.IsPasswordCorrect(sampleUser, samplePassword+"x"))

	assert.NotNil(t, repo.CreateUser(sampleUser, samplePassword+"x", false))

	assert.Nil(t, repo.DeleteUser(sampleUser))
	assert.False(t, repo.IsPasswordCorrect(sampleUser, samplePassword))
}

func TestDoesUserExist(t *testing.T) {
	defer repo.WipeDatabase()
	assert.False(t, repo.DoesUserExist(sampleUser))
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.True(t, repo.DoesUserExist(sampleUser))
	assert.Nil(t, repo.DeleteUser(sampleUser))
	assert.False(t, repo.DoesUserExist(sampleUser))
}

func TestGetUserWithCookie(t *testing.T) {
	defer repo.WipeDatabase()
	_, err := repo.GetUserViaCookie(sampleCookie)
	assert.NotNil(t, err)

	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err := repo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)
	assert.False(t, auth.IsAdmin)
	repo.WipeDatabase()

	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err = repo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)
	assert.True(t, auth.IsAdmin)
}

func TestDoesAnyAdminUserExist(t *testing.T) {
	defer repo.WipeDatabase()
	assert.False(t, repo.DoesAnyAdminUserExist())
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.False(t, repo.DoesAnyAdminUserExist())
	assert.Nil(t, repo.CreateUser(sampleUser+"x", samplePassword, true))
	assert.True(t, repo.DoesAnyAdminUserExist())
}

func TestLogout(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now()))
	auth, err := repo.GetUserViaCookie(sampleCookie)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, auth.User)

	assert.Nil(t, repo.Logout(sampleUser))
	assert.True(t, repo.DoesUserExist(sampleUser))
	auth, err = repo.GetUserViaCookie(sampleCookie)
	assert.NotNil(t, err)
	assert.Nil(t, auth)
}

func TestChangePassword(t *testing.T) {
	defer repo.WipeDatabase()
	oldPassword := samplePassword
	newPassword := samplePassword + "x"
	assert.Nil(t, repo.CreateUser(sampleUser, oldPassword, false))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, oldPassword))

	assert.Nil(t, repo.ChangePassword(sampleUser, newPassword))
	assert.False(t, repo.IsPasswordCorrect(sampleUser, oldPassword))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, newPassword))
}

// TODO Trigger already existing and not existing errors.
func TestAppLifecycle(t *testing.T) {
	defer repo.WipeDatabase()
	assertEmptyAppAndTags(t)

	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	maintainersAndApps, err := repo.ListAppInfo()
	assert.Nil(t, err)
	assert.NotNil(t, maintainersAndApps)
	assert.Equal(t, 1, len(maintainersAndApps))
	assert.Equal(t, sampleMaintainer, maintainersAndApps[0].Maintainer)
	assert.Equal(t, sampleApp, maintainersAndApps[0].App)

	tags, err := repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, sampleTag, tags[0])

	blob, err := repo.LoadTagBlob(sampleMaintainer, sampleApp, sampleTag)
	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, sampleBlob, blob)

	// TODO Deleting the only tag left should also delete the app.
	// TODO Test creating a second app with different tag, as this should not cause collisions if handled correctly.
}

func assertEmptyAppAndTags(t *testing.T) {
	maintainersAndApps, err := repo.ListAppInfo()
	assert.Nil(t, err)
	assert.Nil(t, maintainersAndApps)

	tags, err := repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.NotNil(t, err)
	assert.Nil(t, tags)

	blob, err := repo.LoadTagBlob(sampleMaintainer, sampleApp, sampleTag)
	assert.NotNil(t, err)
	assert.Nil(t, blob)
}

func TestDeleteApp(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, repo.DeleteApp(sampleMaintainer, sampleApp))

	apps, err := repo.ListAppInfo()
	assert.Nil(t, err)
	assert.Nil(t, apps)

	tags, err := repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.NotNil(t, err)
	assert.Nil(t, tags)
}

func TestCreatingTwoTagsInApp(t *testing.T) {
	defer repo.WipeDatabase()
	sampleTag2 := "2.0"
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag2, sampleBlob))

	app, err := repo.ListAppInfo()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(app))
	assert.Equal(t, sampleMaintainer, app[0].Maintainer)
	assert.Equal(t, sampleApp, app[0].App)

	tags, err := repo.ListTagsOfApp(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tags))

	assert.True(t, contain(tags, sampleTag))
	assert.True(t, contain(tags, sampleTag2))
}

func contain(tags []string, expectedTag string) bool {
	for _, actualTag := range tags {
		if actualTag == expectedTag {
			return true
		}
	}
	return false
}

func TestDeleteTag(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, repo.DeleteTag(sampleMaintainer, sampleApp, sampleTag))
	// TODO assertEmptyAppAndTags(t)
}

// TODO check if expiration is working
// TODO can't set a cookie without user
// TODO all inconsistencies should be handled in this layer -> user does not exist, user already existing etc.
// TODO error: user already exists
// TODO SetCookie, DeleteCookie, IsCookieValid
// TODO the DB interface appears to grow quite large when all all use cases are implemented. Check if could be split up.
// TODO Test deletion cascading, e.g. deleting user should also delete his group memberships etc.
// TODO Replace "SERIAL" in the schemes by "INT"
