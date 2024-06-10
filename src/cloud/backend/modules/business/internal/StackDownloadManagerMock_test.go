package internal

import (
	"ocelot/tools"
	"testing"
)

func TestNginxDownloadShouldTriggerDownloadState(t *testing.T) {
	stackName := tools.NginxDownloading
	manager := ProvideDownloadManagerMock()
	tools.AssertEqual(t, 0, len(manager.GetStackDownloadStates()))

	manager.DownloadStack(stackName)

	result := manager.GetStackDownloadStates()
	tools.AssertEqual(t, 1, len(result))
	tools.AssertEqual(t, Ongoing, result[stackName])

	result = manager.GetStackDownloadStates()
	tools.AssertEqual(t, 1, len(result))
	tools.AssertEqual(t, Finished, result[stackName])

	manager.DownloadStack(stackName)
	result = manager.GetStackDownloadStates()
	tools.AssertEqual(t, 1, len(result))
	tools.AssertEqual(t, Ongoing, result[stackName])
}

func TestNginxDefaultDownloadShouldFinishImmediately(t *testing.T) {
	stackName := tools.NginxDefault
	manager := ProvideDownloadManagerMock()
	tools.AssertEqual(t, 0, len(manager.GetStackDownloadStates()))

	manager.DownloadStack(stackName)

	result := manager.GetStackDownloadStates()
	tools.AssertEqual(t, 1, len(result))
	tools.AssertEqual(t, Finished, result[stackName])
}
