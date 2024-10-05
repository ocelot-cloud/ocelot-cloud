package security

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

func TestAppAccessDeletionCascading(t *testing.T) {
	defer dbRepo.WipeDatabase()
	assert.Nil(t, AppRepo.CreateAppWithTag(sampleMaintainer, sampleApp, sampleTag, sampleBlob))
	assert.Nil(t, GroupRepo.CreateGroup(sampleGroup))

	// TODO repo.IsGroupAccessToAppTableEmpty()
}
