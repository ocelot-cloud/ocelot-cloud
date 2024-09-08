package apps

import (
	"os/exec"
	"sync"
)

type downloadState int

const (
	ongoing downloadState = iota
	finished
	failure
)

func (s *downloadState) toString() string {
	return [...]string{"Ongoing", "Finished", "Error"}[*s]
}

type stackDownloadState struct {
	stackName string
	State     downloadState
}

type downloaderReal struct {
	mu                      sync.Mutex
	downloadStates          []*stackDownloadState
	downloadProcessProvider DownloadProcessProvider
}

func provideDownloaderReal() *downloaderReal {
	return &downloaderReal{downloadProcessProvider: &DownloadProcessProviderReal{}}
}

func (s *downloaderReal) getDownloadStates() map[string]downloadState {
	s.mu.Lock()
	defer s.mu.Unlock()
	downloadStateClone := make(map[string]downloadState)
	for _, v := range s.downloadStates {
		downloadStateClone[v.stackName] = v.State
	}
	return downloadStateClone
}

func (s *downloaderReal) download(stackName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, state := range s.downloadStates {
		if state.stackName == stackName {
			state.State = ongoing
			s.downloadProcessProvider.StartDownloadProcessAndSetStateWhenFinished(state)
			return
		}
	}

	newDownloadState := &stackDownloadState{stackName: stackName, State: ongoing}
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
		stackDockerComposePath := appFileDir + "/" + stackDownloadState.stackName + "/docker-compose.yml"
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
