package internal

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/config"
	"testing"
)

func TestNginxDownloadShouldTriggerDownloadState(t *testing.T) {
	stackName := tools.NginxDownloading
	manager := ProvideDownloadManagerMock()
	assert.Equal(t, 0, len(manager.GetStackDownloadStates()))

	manager.DownloadStack(stackName)

	result := manager.GetStackDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, Ongoing, result[stackName])

	result = manager.GetStackDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, Finished, result[stackName])

	manager.DownloadStack(stackName)
	result = manager.GetStackDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, Ongoing, result[stackName])
}

func TestNginxDefaultDownloadShouldFinishImmediately(t *testing.T) {
	stackName := tools.NginxDefault
	manager := ProvideDownloadManagerMock()
	assert.Equal(t, 0, len(manager.GetStackDownloadStates()))

	manager.DownloadStack(stackName)

	result := manager.GetStackDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, Finished, result[stackName])
}
