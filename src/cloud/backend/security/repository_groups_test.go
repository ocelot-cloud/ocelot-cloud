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
