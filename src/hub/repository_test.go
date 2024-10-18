//go:build unit

package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"os"
	"sort"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	InitializeDatabaseWithSource(":memory:")
	code := m.Run()
	os.Exit(code)
}

func TestCreateRepoUser(t *testing.T) {
	defer repo.WipeDatabase()
	assert.False(t, repo.DoesUserExist(sampleUser))
	assert.Nil(t, repo.CreateUser(sampleForm))
	assert.True(t, repo.DoesUserExist(sampleUser))

	assert.Nil(t, repo.DeleteUser(sampleUser))
	assert.False(t, repo.DoesUserExist(sampleUser))
}

func TestCreateRepoApp(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	// TODO Assert app table is empty instead?
	appId, err := repo.CreateApp(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.True(t, repo.DoesAppExist(appId))
}

func TestDeleteAppCascadingThroughUser(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	appId, err := repo.CreateApp(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.True(t, repo.DoesAppExist(appId))
	assert.Nil(t, repo.DeleteApp(appId))
	assert.False(t, repo.DoesAppExist(appId))
}

func TestDeleteAppDirectly(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	appId, err := repo.CreateApp(sampleUser, sampleApp)
	assert.Nil(t, err)
	assert.True(t, repo.DoesAppExist(appId))
	assert.Nil(t, repo.DeleteUser(sampleUser))
	assert.False(t, repo.DoesAppExist(appId))
}

func TestCantCreateUserTwice(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	assert.NotNil(t, repo.CreateUser(sampleForm))
}

func TestCantCreateAppTwiceForSameUser(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	_, err := repo.CreateApp(sampleUser, sampleApp)
	assert.Nil(t, err)
	_, err = repo.CreateApp(sampleUser, sampleApp)
	assert.NotNil(t, err)
}

func TestCantCreateAppWithoutUser(t *testing.T) {
	defer repo.WipeDatabase()
	_, err := repo.CreateApp(sampleUser, sampleApp)
	assert.NotNil(t, err)
}

func TestTolerateSamePasswordForTwoUsers(t *testing.T) {
	defer repo.WipeDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, repo.CreateUser(sampleForm))
	newForm := *sampleForm
	newForm.User = user2
	assert.Nil(t, repo.CreateUser(&newForm))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, samplePassword))
	assert.True(t, repo.IsPasswordCorrect(user2, samplePassword))
}

func TestTolerateSameAppsForTwoUsers(t *testing.T) {
	defer repo.WipeDatabase()
	user2 := sampleUser + "2"
	assert.Nil(t, repo.CreateUser(sampleForm))
	newForm := *sampleForm
	newForm.User = user2
	assert.Nil(t, repo.CreateUser(&newForm))
	appId1, err := repo.CreateApp(sampleUser, sampleApp)
	assert.Nil(t, err)
	appId2, err := repo.CreateApp(user2, sampleApp)
	assert.Nil(t, err)

	assert.True(t, repo.DoesAppExist(appId1))
	assert.True(t, repo.DoesAppExist(appId2))

	assert.Nil(t, repo.DeleteApp(appId1))
	assert.False(t, repo.DoesAppExist(appId1))
	assert.True(t, repo.DoesAppExist(appId2))
}

func TestPasswordVerification(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, samplePassword))
	assert.False(t, repo.IsPasswordCorrect(sampleUser, samplePassword+"x"))
}

func TestSearch(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	app1 := "prefix_myapp_suffix"
	app2 := "prefix_another-app_suffix"
	_, err := repo.CreateApp(sampleUser, app1)
	assert.Nil(t, err)
	_, err = repo.CreateApp(sampleUser, app2)
	assert.Nil(t, err)

	foundApps, err := repo.FindApps("app")
	assert.Nil(t, err)

	sort.Slice(foundApps, func(i, j int) bool {
		return foundApps[i].Name < foundApps[j].Name
	})

	assert.Equal(t, 2, len(foundApps))
	assert.Equal(t, sampleUser, foundApps[0].Maintainer)
	assert.Equal(t, sampleUser, foundApps[1].Maintainer)
	assert.Equal(t, app2, foundApps[0].Name)
	assert.Equal(t, app1, foundApps[1].Name)
}

func TestSearchNegative(t *testing.T) {
	defer repo.WipeDatabase()
	repo.CreateUser(sampleForm)
	app := "prefix_myapp_suffix"
	repo.CreateApp(sampleUser, app)

	a, err := repo.FindApps("some")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(a))
}

func TestCookieExpiration(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	_, err := repo.GetUserWithCookie("")
	assert.NotNil(t, err)

	assert.True(t, repo.IsCookieExpired("non-existing-cookie"))

	timeIn30Days := utils.GetTimeIn30Days()
	cookie, _ := utils.GenerateCookie()
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, cookie.Value, timeIn30Days))
	assert.False(t, repo.IsCookieExpired(cookie.Value))

	past := time.Now().Add(-1 * time.Second)
	assert.Nil(t, repo.HashAndSaveCookie(sampleUser, cookie.Value, past))
	assert.True(t, repo.IsCookieExpired(cookie.Value))

	user, err := repo.GetUserWithCookie(cookie.Value)
	assert.Nil(t, err)
	assert.Equal(t, sampleUser, user)
}

func TestGetTagList(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	appId, err := repo.CreateApp(sampleUser, sampleApp)
	assert.Nil(t, err)
	foundTags, err := repo.GetTagList(appId)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundTags))
	tagId, err := repo.GetTagId(appId, sampleTag)
	assert.NotNil(t, err)
	assert.False(t, repo.DoesTagExist(tagId))

	assert.Nil(t, repo.CreateTag(appId, sampleTag, []byte("asdf")))
	foundTags, err = repo.GetTagList(appId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, sampleTag, foundTags[0].Name)
	tagId, err = repo.GetTagId(appId, sampleTag)
	assert.Nil(t, err)
	assert.True(t, repo.DoesTagExist(tagId))
	data, err := repo.GetTagContent(tagId)
	assert.Nil(t, err)
	assert.Equal(t, []byte("asdf"), data)

	assert.Nil(t, repo.DeleteTag(tagId))
	foundTags, err = repo.GetTagList(appId)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(foundTags))
	assert.False(t, repo.DoesTagExist(tagId))

	assert.Nil(t, repo.CreateTag(appId, sampleTag, []byte("asdf")))
	foundTags, err = repo.GetTagList(appId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(foundTags))
	assert.Equal(t, sampleTag, foundTags[0].Name)
	tagId, err = repo.GetTagId(appId, sampleTag)
	assert.Nil(t, err)
	assert.True(t, repo.DoesTagExist(tagId))
	assert.Nil(t, repo.DeleteUser(sampleUser))
	tags, err := repo.GetTagList(appId) // TODO not sure, should it fail when app or user do not exist?
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tags))
	assert.False(t, repo.DoesTagExist(tagId))
}

// TODO Make checks that GetAppId gives the same results as GetAppList[0].Id, same for tags

func TestChangeRepoPassword(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, samplePassword))
	newPassword := samplePassword + "x"
	assert.Nil(t, repo.ChangePassword(sampleUser, newPassword))
	assert.False(t, repo.IsPasswordCorrect(sampleUser, sampleForm.Password))
	assert.True(t, repo.IsPasswordCorrect(sampleUser, newPassword))
}

func TestChangeRepoOrigin(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	assert.False(t, repo.IsOriginCorrect(sampleUser, sampleOrigin))
	assert.Nil(t, repo.SetOrigin(sampleUser, sampleOrigin))
	assert.True(t, repo.IsOriginCorrect(sampleUser, sampleOrigin))
}

func TestUsedSpace(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	space, err := repo.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, space)
	appId, err := repo.CreateApp(sampleUser, sampleApp)
	assert.Nil(t, err)

	bytes := []byte("hello")
	bytes2 := []byte(" world")
	assert.Nil(t, repo.CreateTag(appId, sampleTag, bytes))
	space, err = repo.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 5, space)

	assert.Nil(t, repo.CreateTag(appId, sampleTag+"x", bytes2))
	space, err = repo.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 11, space)

	tagId, err := repo.GetTagId(appId, sampleTag)
	assert.Nil(t, err)
	assert.Nil(t, repo.DeleteTag(tagId))
	space, err = repo.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 6, space)

	assert.Nil(t, repo.CreateTag(appId, sampleTag, bytes2))
	space, err = repo.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 12, space)

	assert.Nil(t, repo.DeleteApp(appId))
	space, err = repo.GetUsedSpaceInBytes(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, space)
}

func TestRepoLogout(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	sampleCookie := "asdasdasd"
	err := repo.HashAndSaveCookie(sampleUser, sampleCookie, time.Now().Add(1*time.Hour))
	assert.Nil(t, err)
	assert.False(t, repo.IsCookieExpired(sampleCookie))
	assert.Nil(t, repo.Logout(sampleUser))
	assert.True(t, repo.IsCookieExpired(sampleCookie))
}

func TestGetAppListRepo(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	list, err := repo.GetAppList(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))
	_, err = repo.CreateApp(sampleUser, sampleApp)
	assert.Nil(t, err)
	_, err = repo.CreateApp(sampleUser, sampleApp+"x")
	assert.Nil(t, err)
	list, err = repo.GetAppList(sampleUser)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, sampleApp, list[0].Name)
	assert.Equal(t, sampleApp+"x", list[1].Name)
}

func TestConcurrencyRobustness(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleForm))
	for i := 0; i < 1000; i++ {
		go func() {
			_, err := repo.GetAppList(sampleUser)
			assert.Nil(t, err)
		}()
	}
}
