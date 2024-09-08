package image_download

import (
	"github.com/ocelot-cloud/shared/assert"
	"ocelot/backend/tools"
	"testing"
)

func TestNginxDownloadShouldTriggerDownloadState(t *testing.T) {
	stackName := tools.NginxDownloading
	manager := ProvideDownloaderMock()
	assert.Equal(t, 0, len(manager.GetDownloadStates()))

	manager.Download(stackName)

	result := manager.GetDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, Ongoing, result[stackName])

	result = manager.GetDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, finished, result[stackName])

	manager.Download(stackName)
	result = manager.GetDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, Ongoing, result[stackName])
}

func TestNginxDefaultDownloadShouldFinishImmediately(t *testing.T) {
	stackName := tools.NginxDefault
	manager := ProvideDownloaderMock()
	assert.Equal(t, 0, len(manager.GetDownloadStates()))

	manager.Download(stackName)

	result := manager.GetDownloadStates()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, finished, result[stackName])
}
