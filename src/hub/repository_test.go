//go:build unit

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/hub"
	"github.com/ocelot-cloud/shared/utils"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	initializeDatabaseWithSource(":memory:")
	code := m.Run()
	os.Exit(code)
}

func TestCreateRepoUser(t *testing.T) {
	defer repo.WipeDatabase()
	assert.False(t, repo.DoesUserExist(hub.SampleUser))
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.True(t, repo.DoesUserExist(hub.SampleUser))

	assert.Nil(t, repo.DeleteUser(hub.SampleUser))
	assert.False(t, repo.DoesUserExist(hub.SampleUser))
}

func TestCreateRepoApp(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.False(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
	assert.True(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
}

func TestDeleteAppCascadingThroughUser(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
	assert.True(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
	assert.Nil(t, repo.DeleteApp(hub.SampleUser, hub.SampleApp))
	assert.False(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
}

func TestDeleteAppDirectly(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.False(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
	assert.True(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
	assert.Nil(t, repo.DeleteUser(hub.SampleUser))
	assert.False(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
}

func TestCantCreateUserTwice(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.NotNil(t, repo.CreateUser(hub.SampleForm))
}

func TestCantCreateAppTwiceForSameUser(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
	assert.NotNil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
}

func TestCantCreateAppWithoutUser(t *testing.T) {
	defer repo.WipeDatabase()
	assert.NotNil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
}

func TestTolerateSamePasswordForTwoUsers(t *testing.T) {
	defer repo.WipeDatabase()
	user2 := hub.SampleUser + "2"
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	newForm := *hub.SampleForm
	newForm.User = user2
	assert.Nil(t, repo.CreateUser(&newForm))
	assert.True(t, repo.IsPasswordCorrect(hub.SampleUser, hub.SamplePassword))
	assert.True(t, repo.IsPasswordCorrect(user2, hub.SamplePassword))
}

func TestTolerateSameAppsForTwoUsers(t *testing.T) {
	defer repo.WipeDatabase()
	user2 := hub.SampleUser + "2"
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	newForm := *hub.SampleForm
	newForm.User = user2
	assert.Nil(t, repo.CreateUser(&newForm))
	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
	assert.Nil(t, repo.CreateApp(user2, hub.SampleApp))

	assert.True(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
	assert.True(t, repo.DoesAppExist(user2, hub.SampleApp))

	assert.Nil(t, repo.DeleteApp(hub.SampleUser, hub.SampleApp))
	assert.False(t, repo.DoesAppExist(hub.SampleUser, hub.SampleApp))
	assert.True(t, repo.DoesAppExist(user2, hub.SampleApp))
}

func TestPasswordVerification(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.True(t, repo.IsPasswordCorrect(hub.SampleUser, hub.SamplePassword))
	assert.False(t, repo.IsPasswordCorrect(hub.SampleUser, hub.SamplePassword+"x"))
}

func TestSearch(t *testing.T) {
	defer repo.WipeDatabase()
	repo.CreateUser(hub.SampleForm)
	app1 := "prefix_myapp_suffix"
	app2 := "prefix_another-app_suffix"
	repo.CreateApp(hub.SampleUser, app1)
	repo.CreateApp(hub.SampleUser, app2)

	a, err := repo.FindApps("app")
	assert.Nil(t, err)

	sort.Slice(a, func(i, j int) bool {
		return a[i].App < a[j].App
	})

	assert.Equal(t, 2, len(a))
	assert.Equal(t, hub.SampleUser, a[0].User)
	assert.Equal(t, hub.SampleUser, a[1].User)
	assert.Equal(t, app2, a[0].App)
	assert.Equal(t, app1, a[1].App)
}

func TestSearchNegative(t *testing.T) {
	defer repo.WipeDatabase()
	repo.CreateUser(hub.SampleForm)
	app := "prefix_myapp_suffix"
	repo.CreateApp(hub.SampleUser, app)

	a, err := repo.FindApps("some")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(a))
}

func TestCookieExpiration(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	_, err := repo.GetUserWithCookie("")
	assert.NotNil(t, err)

	assert.True(t, repo.IsCookieExpired("non-existing-cookie"))

	timeIn30Days := utils.GetTimeIn30Days()
	cookie, _ := utils.GenerateCookie()
	assert.Nil(t, repo.HashAndSaveCookie(hub.SampleUser, cookie.Value, timeIn30Days))
	assert.False(t, repo.IsCookieExpired(cookie.Value))

	past := time.Now().Add(-1 * time.Second)
	assert.Nil(t, repo.HashAndSaveCookie(hub.SampleUser, cookie.Value, past))
	assert.True(t, repo.IsCookieExpired(cookie.Value))

	user, err := repo.GetUserWithCookie(cookie.Value)
	assert.Nil(t, err)
	assert.Equal(t, hub.SampleUser, user)
}

func TestGetTagList(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
	foundTags, err := repo.GetTagList(hub.SampleUser, hub.SampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundTags))
	assert.False(t, repo.DoesTagExist(hub.SampleUser, hub.SampleApp, hub.SampleTag))

	assert.Nil(t, repo.CreateTag(hub.SampleUser, hub.SampleApp, hub.SampleTag, []byte("asdf")))
	foundTags, err = repo.GetTagList(hub.SampleUser, hub.SampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, hub.SampleTag, foundTags[0])
	assert.True(t, repo.DoesTagExist(hub.SampleUser, hub.SampleApp, hub.SampleTag))
	data, err := repo.GetTagContent(hub.SampleUser, hub.SampleApp, hub.SampleTag)
	assert.Nil(t, err)
	assert.Equal(t, []byte("asdf"), data)

	assert.Nil(t, repo.DeleteTag(hub.SampleUser, hub.SampleApp, hub.SampleTag))
	foundTags, err = repo.GetTagList(hub.SampleUser, hub.SampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundTags))
	assert.False(t, repo.DoesTagExist(hub.SampleUser, hub.SampleApp, hub.SampleTag))

	assert.Nil(t, repo.CreateTag(hub.SampleUser, hub.SampleApp, hub.SampleTag, []byte("asdf")))
	foundTags, err = repo.GetTagList(hub.SampleUser, hub.SampleApp)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, hub.SampleTag, foundTags[0])
	assert.True(t, repo.DoesTagExist(hub.SampleUser, hub.SampleApp, hub.SampleTag))
	assert.Nil(t, repo.DeleteUser(hub.SampleUser))
	_, err = repo.GetTagList(hub.SampleUser, hub.SampleApp)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "user not found"))
	assert.False(t, repo.DoesTagExist(hub.SampleUser, hub.SampleApp, hub.SampleTag))
}

func TestChangeRepoPassword(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.True(t, repo.IsPasswordCorrect(hub.SampleUser, hub.SamplePassword))
	newPassword := hub.SamplePassword + "x"
	assert.Nil(t, repo.ChangePassword(hub.SampleUser, newPassword))
	assert.False(t, repo.IsPasswordCorrect(hub.SampleUser, hub.SampleForm.Password))
	assert.True(t, repo.IsPasswordCorrect(hub.SampleUser, newPassword))
}

func TestChangeRepoOrigin(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	assert.False(t, repo.IsOriginCorrect(hub.SampleUser, hub.SampleOrigin))
	assert.Nil(t, repo.SetOrigin(hub.SampleUser, hub.SampleOrigin))
	assert.True(t, repo.IsOriginCorrect(hub.SampleUser, hub.SampleOrigin))
}

func TestUsedSpace(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	space, err := repo.GetUsedSpaceInBytes(hub.SampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, space)

	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))

	bytes := []byte("hello")
	bytes2 := []byte(" world")
	assert.Nil(t, repo.CreateTag(hub.SampleUser, hub.SampleApp, hub.SampleTag, bytes))
	space, err = repo.GetUsedSpaceInBytes(hub.SampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 5, space)

	assert.Nil(t, repo.CreateTag(hub.SampleUser, hub.SampleApp, hub.SampleTag+"x", bytes2))
	space, err = repo.GetUsedSpaceInBytes(hub.SampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 11, space)

	assert.Nil(t, repo.DeleteTag(hub.SampleUser, hub.SampleApp, hub.SampleTag))
	space, err = repo.GetUsedSpaceInBytes(hub.SampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 6, space)

	assert.Nil(t, repo.CreateTag(hub.SampleUser, hub.SampleApp, hub.SampleTag, bytes2))
	space, err = repo.GetUsedSpaceInBytes(hub.SampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 12, space)

	assert.Nil(t, repo.DeleteApp(hub.SampleUser, hub.SampleApp))
	space, err = repo.GetUsedSpaceInBytes(hub.SampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, space)
}

func TestRepoLogout(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	sampleCookie := "asdasdasd"
	err := repo.HashAndSaveCookie(hub.SampleUser, sampleCookie, time.Now().Add(1*time.Hour))
	assert.Nil(t, err)
	assert.False(t, repo.IsCookieExpired(sampleCookie))
	assert.Nil(t, repo.Logout(hub.SampleUser))
	assert.True(t, repo.IsCookieExpired(sampleCookie))
}

func TestGetAppListRepo(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	list, err := repo.GetAppList(hub.SampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))
	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp))
	assert.Nil(t, repo.CreateApp(hub.SampleUser, hub.SampleApp+"x"))
	list, err = repo.GetAppList(hub.SampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, hub.SampleApp, list[0])
	assert.Equal(t, hub.SampleApp+"x", list[1])
}

func TestConcurrencyRobustness(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(hub.SampleForm))
	for i := 0; i < 1000; i++ {
		go func() {
			_, err := repo.GetAppList(hub.SampleUser)
			assert.Nil(t, err)
		}()
	}
}
