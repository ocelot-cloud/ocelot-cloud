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
	defer repo.WipeDatabase()
	groups, err := repo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(groups))

	assert.Nil(t, repo.CreateGroup(sampleGroup))
	groups, err = repo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(groups))
	assert.Equal(t, sampleGroup, groups[0])

	assert.Nil(t, repo.DeleteGroup(sampleGroup))
	groups, err = repo.ListGroups()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(groups))
}

func TestListAllUsers(t *testing.T) {
	defer repo.WipeDatabase()
	users, err := repo.ListAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(users))

	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, false))
	users, err = repo.ListAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, sampleUser, users[0])
}

func TestAddUserToGroup(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateUser(sampleUser, samplePassword, true))
	assert.Nil(t, repo.CreateGroup(sampleGroup))

	members, err := repo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(members))

	assert.Nil(t, repo.AddUserToGroup(sampleUser, sampleGroup))

	members, err = repo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(members))
	assert.Equal(t, sampleUser, members[0])

	assert.Nil(t, repo.RemoveUserFromGroup(sampleUser, sampleGroup))

	members, err = repo.ListMembersOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(members))
	assert.Nil(t, repo.AddUserToGroup(sampleUser, sampleGroup))
}

func TestGiveGroupAccessToApp(t *testing.T) {
	defer repo.WipeDatabase()
	assert.Nil(t, repo.CreateGroup(sampleGroup))
	assert.Nil(t, repo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))

	accessList, err := repo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(accessList))

	assert.Nil(t, repo.GiveGroupAccessToApp(sampleGroup, sampleMaintainerAndApp))

	accessList, err = repo.ListAppAccessesOfGroup(sampleGroup)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(accessList))
	assert.Equal(t, sampleMaintainer, accessList[0].Maintainer)
	assert.Equal(t, sampleApp, accessList[0].App)
}

func TestRemoveGroupsAccessToApp(t *testing.T) {

}

// TODO After finishing the persistence layer, I should add services with business logic, which handle all unhappy path cases.
// TODO Add assertion function that all tables are empty? Can be used to test deletion cascading.
