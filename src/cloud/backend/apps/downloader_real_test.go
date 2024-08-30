package apps

import (
	"github.com/ocelot-cloud/shared/assert"
	"os"
	"os/exec"
	"testing"
	"time"
)

var testStack = "test-stack"
var testStack2 = "test-stack2"

var downloadProcessProviderMock *DownloadProcessProviderMock
var stackDownloadManager StackDownloadManager

func setup() {
	downloadProcessProviderMock = &DownloadProcessProviderMock{}
	stackDownloadManager = &downloaderReal{downloadProcessProvider: downloadProcessProviderMock}
}

func TestDownloadStack_InitialState(t *testing.T) {
	setup()

	downloadStates := stackDownloadManager.getDownloadStates()
	assert.Equal(t, 0, len(downloadStates))
}

func TestDownloadStack_SingleDownload(t *testing.T) {
	setup()

	stackDownloadManager.download(testStack)
	downloadStates := stackDownloadManager.getDownloadStates()

	assert.Equal(t, 1, len(downloadStates))
	assert.Equal(t, downloadStates[testStack], ongoing)
}

func TestDownloadStack_DuplicateDownloadDoesNotCreateNewDownloadState(t *testing.T) {
	setup()

	stackDownloadManager.download(testStack)
	stackDownloadManager.download(testStack)
	downloadStates := stackDownloadManager.getDownloadStates()

	assert.Equal(t, 1, len(downloadStates))
	assert.Equal(t, downloadStates[testStack], ongoing)
}

func TestDownloadStack_FinishedDownloadState(t *testing.T) {
	setup()

	stackDownloadManager.download(testStack)
	downloadProcessProviderMock.stackDownloadState.State = finished
	downloadStates := stackDownloadManager.getDownloadStates()

	assert.Equal(t, 1, len(downloadStates))
	assert.Equal(t, downloadStates[testStack], finished)
}

func TestDownloadStack_AllowDownloadSecondTime(t *testing.T) {
	setup()

	stackDownloadManager.download(testStack)
	downloadProcessProviderMock.stackDownloadState.State = finished
	stackDownloadManager.download(testStack)
	downloadStates := stackDownloadManager.getDownloadStates()

	assert.Equal(t, 1, len(downloadStates))
	assert.Equal(t, downloadStates[testStack], ongoing)
}

func TestDownloadStack_ErrorState(t *testing.T) {
	setup()

	stackDownloadManager.download(testStack)
	downloadProcessProviderMock.stackDownloadState.State = failure
	downloadStates := stackDownloadManager.getDownloadStates()

	assert.Equal(t, 1, len(downloadStates))
	assert.Equal(t, downloadStates[testStack], failure)
}

func TestDownloadStack_MultipleDownloads(t *testing.T) {
	setup()

	stackDownloadManager.download(testStack)
	stackDownloadManager.download(testStack2)
	downloadStates := stackDownloadManager.getDownloadStates()

	assert.Equal(t, 2, len(downloadStates))
	assert.Equal(t, downloadStates[testStack], ongoing)
}

func TestDownloadProcessProviderReal(t *testing.T) {
	if os.Getenv("IS_IMAGE_DOWNLOAD_TEST") != "true" {
		t.Skip()
		return
	}

	cmd := exec.Command("docker", "rmi", "-f", "nginx:alpine3.17")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to delete docker image nginx:alpine3.17: %v", err)
	}

	downloader := DownloadProcessProviderReal{}
	downloadState := &stackDownloadState{"nginx-download", ongoing}

	downloader.StartDownloadProcessAndSetStateWhenFinished(downloadState)

	timeout := time.After(60 * time.Second)
	tick := time.Tick(20 * time.Millisecond)
	for {
		select {
		case <-timeout:
			t.Fatal("Test failed: Download did not finish in time.")
		case <-tick:
			if downloadState.State == finished {
				logger.Info("Download finished successfully in time.")
				return
			}
		}
	}
}
