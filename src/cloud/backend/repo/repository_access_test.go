package repo

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestGiveGroupAccessToApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, GroupRepo.CreateGroup(sampleGroup))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))

	groupId, err := GroupRepo.GetGroupId(sampleGroup)
	assert.Nil(t, err)
	appsToWhichGroupHasAccess, err := AccessRepo.ListAppAccessesOfGroup(groupId)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(appsToWhichGroupHasAccess))

	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(groupId, appId))

	appsToWhichGroupHasAccess, err = AccessRepo.ListAppAccessesOfGroup(groupId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(appsToWhichGroupHasAccess))
	assert.Equal(t, sampleMaintainer, appsToWhichGroupHasAccess[0].Maintainer)
	assert.Equal(t, sampleApp, appsToWhichGroupHasAccess[0].Name)

	assert.Nil(t, AccessRepo.RemoveGroupsAccessToApp(sampleGroup, appId))

	appsToWhichGroupHasAccess, err = AccessRepo.ListAppAccessesOfGroup(groupId)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(appsToWhichGroupHasAccess))
	appId, err = AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)

	// TODO Not sure, but shouldn't there be a test which says that a user has the according access when he is member of a group with access rights to an app?
	err = UserRepo.CreateUser(sampleUser, samplePassword, false)
	assert.Nil(t, err)
	userId, err := UserRepo.GetUserId(sampleUser)
	assert.Nil(t, err)
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(userId, appId))
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
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(userId, appId))
	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(groupId, appId))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(userId, appId))
	assert.Nil(t, AccessRepo.RemoveGroupsAccessToApp(sampleGroup, appId))
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(userId, appId))

	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(groupId, appId))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(userId, appId))
	assert.Nil(t, GroupRepo.RemoveUserFromGroup(userId, groupId))
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(userId, appId))
	assert.Nil(t, GroupRepo.AddUserToGroup(userId, groupId))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(userId, appId))
}

func TestAppAccessDeletionCascading(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	appId, err := AppRepo.GetAppId(sampleMaintainer, sampleApp)
	assert.Nil(t, err)
	userId, err := UserRepo.GetUserId(sampleUser)
	assert.Nil(t, err)
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(userId, appId))
}

// TODO Admins should always have access to all apps, not matter which groups they are in or not in
// TODO I should also check that there are no residues in the database when deleting items. Could lead to security issues otherwise.
