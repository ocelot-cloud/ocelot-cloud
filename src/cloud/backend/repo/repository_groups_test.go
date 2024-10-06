package repo

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
	groups, err := GroupRepo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(groups))

	assert.Nil(t, GroupRepo.CreateGroup(sampleGroup))
	groups, err = GroupRepo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(groups))
	assert.Equal(t, sampleGroup, groups[0])

	assert.Nil(t, GroupRepo.DeleteGroup(sampleGroup))
	groups, err = GroupRepo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(groups))
}

func TestListAllUsers(t *testing.T) {
	defer dbRepo.WipeDatabase()
	users, err := GroupRepo.ListAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(users))

	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, false))
	users, err = GroupRepo.ListAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, sampleUser, users[0])
}

func TestAddUserToGroup(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, UserRepo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, GroupRepo.CreateGroup(sampleGroup))

	members, err := GroupRepo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(members))

	assert.Nil(t, GroupRepo.AddUserToGroup(sampleUser, sampleGroup))

	members, err = GroupRepo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(members))
	assert.Equal(t, sampleUser, members[0])

	assert.Nil(t, GroupRepo.RemoveUserFromGroup(sampleUser, sampleGroup))

	members, err = GroupRepo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(members))
	assert.Nil(t, GroupRepo.AddUserToGroup(sampleUser, sampleGroup))
}

// TODO After finishing the persistence layer, I should add services with business logic, which handle all unhappy path cases.
// TODO Add assertion function that all tables are empty? Can be used to test deletion cascading.
