package security

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

var (
	sampleGroup = "mygroup"
)

func TestGroupLifecycle(t *testing.T) {
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
}
