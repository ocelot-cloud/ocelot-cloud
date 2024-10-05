package security

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

var (
	sampleGroup            = "mygroup"
	sampleMaintainerAndApp = MaintainerAndApp{sampleMaintainer, sampleApp}
)

func TestGroupLifecycle(t *testing.T) {
	defer dbRepo.WipeDatabase()
	groups, err := groupRepo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(groups))

	assert.Nil(t, groupRepo.CreateGroup(sampleGroup))
	groups, err = groupRepo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(groups))
	assert.Equal(t, sampleGroup, groups[0])

	assert.Nil(t, groupRepo.DeleteGroup(sampleGroup))
	groups, err = groupRepo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(groups))
}

func TestListAllUsers(t *testing.T) {
	defer dbRepo.WipeDatabase()
	users, err := groupRepo.ListAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(users))

	assert.Nil(t, userRepo.CreateUser(sampleUser, samplePassword, false))
	users, err = groupRepo.ListAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, sampleUser, users[0])
}

func TestAddUserToGroup(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, userRepo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, groupRepo.CreateGroup(sampleGroup))

	members, err := groupRepo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(members))

	assert.Nil(t, groupRepo.AddUserToGroup(sampleUser, sampleGroup))

	members, err = groupRepo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(members))
	assert.Equal(t, sampleUser, members[0])

	assert.Nil(t, groupRepo.RemoveUserFromGroup(sampleUser, sampleGroup))

	members, err = groupRepo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(members))
	assert.Nil(t, groupRepo.AddUserToGroup(sampleUser, sampleGroup))
}

func TestGiveGroupAccessToApp(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, groupRepo.CreateGroup(sampleGroup))
	assert.Nil(t, appRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))

	accessList, err := groupRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(accessList))

	assert.Nil(t, groupRepo.GiveGroupAccessToApp(sampleGroup, sampleMaintainerAndApp))

	accessList, err = groupRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(accessList))
	assert.Equal(t, sampleMaintainer, accessList[0].Maintainer)
	assert.Equal(t, sampleApp, accessList[0].App)

	assert.Nil(t, groupRepo.RemoveGroupsAccessToApp(sampleGroup, sampleMaintainerAndApp))

	accessList, err = groupRepo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(accessList))
}

func TestAppAccessDeletionCascading(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, appRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, groupRepo.CreateGroup(sampleGroup))

	// TODO repo.IsGroupAccessToAppTableEmpty()
}

// TODO After finishing the persistence layer, I should add services with business logic, which handle all unhappy path cases.
// TODO Add assertion function that all tables are empty? Can be used to test deletion cascading.
