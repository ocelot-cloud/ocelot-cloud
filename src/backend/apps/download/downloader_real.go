package download

import (
	"ocelot/backend/apps/vars"
	"ocelot/backend/tools"
	"os/exec"
	"sync"
)

type DownloadManager interface {
	GetDownloadStates() map[string]DownloadState
	Download(appName string)
}

var logger = tools.Logger

type DownloadState int

const (
	Ongoing DownloadState = iota
	finished
	failure
)

func (s *DownloadState) toString() string {
	return [...]string{"Ongoing", "Finished", "Error"}[*s]
}

type stackDownloadState struct {
	stackName string
	State     DownloadState
}

type DownloaderReal struct {
	mu                      sync.Mutex
	downloadStates          []*stackDownloadState
	downloadProcessProvider DownloadProcessProvider
}

func ProvideDownloaderReal() *DownloaderReal {
	return &DownloaderReal{downloadProcessProvider: &DownloadProcessProviderReal{}}
}

func (s *DownloaderReal) GetDownloadStates() map[string]DownloadState {
	s.mu.Lock()
	defer s.mu.Unlock()
	downloadStateClone := make(map[string]DownloadState)
	for _, v := range s.downloadStates {
		downloadStateClone[v.stackName] = v.State
	}
	return downloadStateClone
}

func (s *DownloaderReal) Download(stackName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, state := range s.downloadStates {
		if state.stackName == stackName {
			state.State = Ongoing
			s.downloadProcessProvider.StartDownloadProcessAndSetStateWhenFinished(state)
			return
		}
	}

	newDownloadState := &stackDownloadState{stackName: stackName, State: Ongoing}
	s.downloadStates = append(s.downloadStates, newDownloadState)
	s.downloadProcessProvider.StartDownloadProcessAndSetStateWhenFinished(newDownloadState)
}

type DownloadProcessProvider interface {
	StartDownloadProcessAndSetStateWhenFinished(stackDownloadState *stackDownloadState)
}

type DownloadProcessProviderMock struct {
	stackDownloadState *stackDownloadState
}

func (d *DownloadProcessProviderMock) StartDownloadProcessAndSetStateWhenFinished(stackDownloadState *stackDownloadState) {
	d.stackDownloadState = stackDownloadState
}

type DownloadProcessProviderReal struct{}

func (d *DownloadProcessProviderReal) StartDownloadProcessAndSetStateWhenFinished(stackDownloadState *stackDownloadState) {
	go func() {
		stackDockerComposePath := vars.AppFileDir + "/" + stackDownloadState.stackName + "/docker-compose.yml"
		pullCmd := exec.Command("docker", "compose", "-f", stackDockerComposePath, "pull")
		err := pullCmd.Run()
		if err != nil {
			logger.Error("Error executing command '%s': %v\n", pullCmd.String(), err)
			stackDownloadState.State = failure
			return
		}

		buildCmd := exec.Command("docker", "compose", "-f", stackDockerComposePath, "build", "--pull")
		err = buildCmd.Run()
		if err == nil {
			logger.Debug("Successfully downloaded images for stack: %s", stackDownloadState.stackName)
			stackDownloadState.State = finished
		} else {
			logger.Error("Error executing command '%s': %v\n", buildCmd.String(), err)
			stackDownloadState.State = failure
		}
	}()
}
