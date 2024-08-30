package apps

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/tools"
	"testing"
)

func TestNginxDownloadShouldTriggerDownloadState(t *testing.T) {
	stackName := tools.NginxDownloading
	manager := provideDownloaderMock()
	assert.Equal(t, 0, len(manager.getDownloadStates()))

	manager.download(stackName)

	result := manager.getDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, ongoing, result[stackName])

	result = manager.getDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, finished, result[stackName])

	manager.download(stackName)
	result = manager.getDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, ongoing, result[stackName])
}

func TestNginxDefaultDownloadShouldFinishImmediately(t *testing.T) {
	stackName := tools.NginxDefault
	manager := provideDownloaderMock()
	assert.Equal(t, 0, len(manager.getDownloadStates()))

	manager.download(stackName)

	result := manager.getDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, finished, result[stackName])
}
