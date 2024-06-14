package internal

import (
	"github.com/ocelot-cloud/shared"
	"ocelot/backend/config"
	"testing"
)

func TestNginxDownloadShouldTriggerDownloadState(t *testing.T) {
	stackName := tools.NginxDownloading
	manager := ProvideDownloadManagerMock()
	shared.AssertEqual(t, 0, len(manager.GetStackDownloadStates()))

	manager.DownloadStack(stackName)

	result := manager.GetStackDownloadStates()
	shared.AssertEqual(t, 1, len(result))
	shared.AssertEqual(t, Ongoing, result[stackName])

	result = manager.GetStackDownloadStates()
	shared.AssertEqual(t, 1, len(result))
	shared.AssertEqual(t, Finished, result[stackName])

	manager.DownloadStack(stackName)
	result = manager.GetStackDownloadStates()
	shared.AssertEqual(t, 1, len(result))
	shared.AssertEqual(t, Ongoing, result[stackName])
}

func TestNginxDefaultDownloadShouldFinishImmediately(t *testing.T) {
	stackName := tools.NginxDefault
	manager := ProvideDownloadManagerMock()
	shared.AssertEqual(t, 0, len(manager.GetStackDownloadStates()))

	manager.DownloadStack(stackName)

	result := manager.GetStackDownloadStates()
	shared.AssertEqual(t, 1, len(result))
	shared.AssertEqual(t, Finished, result[stackName])
}
