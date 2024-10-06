package repo

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestGiveGroupAccessToApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, GroupRepo.CreateGroup(sampleGroup))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))

	accessList, err := AccessRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(accessList))

	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(sampleGroup, sampleMaintainerAndApp))

	accessList, err = AccessRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(accessList))
	assert.Equal(t, sampleMaintainer, accessList[0].Maintainer)
	assert.Equal(t, sampleApp, accessList[0].App)

	assert.Nil(t, AccessRepo.RemoveGroupsAccessToApp(sampleGroup, sampleMaintainerAndApp))

	accessList, err = AccessRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(accessList))
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(sampleGroup, sampleMaintainerAndApp))
}

func TestUserAccessToApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	assert.Nil(t, GroupRepo.CreateGroup(sampleGroup))
	assert.Nil(t, GroupRepo.AddUserToGroup(sampleUser, sampleGroup))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))

	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, sampleMaintainerAndApp))
	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(sampleGroup, sampleMaintainerAndApp))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, sampleMaintainerAndApp))
	assert.Nil(t, AccessRepo.RemoveGroupsAccessToApp(sampleGroup, sampleMaintainerAndApp))
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, sampleMaintainerAndApp))

	assert.Nil(t, AccessRepo.GiveGroupAccessToApp(sampleGroup, sampleMaintainerAndApp))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, sampleMaintainerAndApp))
	assert.Nil(t, GroupRepo.RemoveUserFromGroup(sampleUser, sampleGroup))
	assert.False(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, sampleMaintainerAndApp))
	assert.Nil(t, GroupRepo.AddUserToGroup(sampleUser, sampleGroup))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, sampleMaintainerAndApp))
}

func TestAppAccessDeletionCascading(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.True(t, AccessRepo.DoesUserHaveAccessToApp(sampleUser, sampleMaintainerAndApp))
}

// TODO Admins should always have access to all apps, not matter which groups they are in or not in
