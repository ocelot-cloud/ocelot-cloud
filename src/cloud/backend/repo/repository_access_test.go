package repo

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestGiveGroupAccessToApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, GroupRepo.CreateGroup(sampleGroup))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))

	appsToWhichGroupHasAccess, err := AccessRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(appsToWhichGroupHasAccess))

	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(sampleGroup, appId))

	appsToWhichGroupHasAccess, err = AccessRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(appsToWhichGroupHasAccess))
	assert.Equal(t, sampleMaintainer, appsToWhichGroupHasAccess[0].Maintainer)
	assert.Equal(t, sampleApp, appsToWhichGroupHasAccess[0].Name)

	assert.Nil(t, AccessRepo.RemoveGroupsAccessToApp(sampleGroup, appId))

	appsToWhichGroupHasAccess, err = AccessRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(appsToWhichGroupHasAccess))
	appId, err = AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(sampleGroup, appId))
}

func TestUserAccessToApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, GroupRepo.CreateGroup(sampleGroup))
	groupId, err := GroupRepo.GetGroupId(sampleGroup)
	assert.Nil(t, err)
	userId, err := UserRepo.GetUserId(sampleUser)
	assert.Nil(t, err)
	assert.Nil(t, GroupRepo.AddUserToGroup(userId, groupId))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))

	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, appId))
	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(sampleGroup, appId))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, appId))
	assert.Nil(t, AccessRepo.RemoveGroupsAccessToApp(sampleGroup, appId))
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, appId))

	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(sampleGroup, appId))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, appId))
	assert.Nil(t, GroupRepo.RemoveUserFromGroup(userId, groupId))
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, appId))
	assert.Nil(t, GroupRepo.AddUserToGroup(userId, groupId))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, appId))
}

func TestAppAccessDeletionCascading(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, appId))
}

// TODO Admins should always have access to all apps, not matter which groups they are in or not in
