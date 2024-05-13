package internal

import (
	"os/exec"
	"sync"
)

type DownloadState int

const (
	Ongoing DownloadState = iota
	Finished
	Error
)

func (s *DownloadState) String() string {
	return [...]string{"Ongoing", "Finished", "Error"}[*s]
}

type StackDownloadState struct {
	stackName string
	State     DownloadState
}

type StackDownloadManagerReal struct {
	mu                      sync.Mutex
	downloadStates          []*StackDownloadState
	downloadProcessProvider DownloadProcessProvider
}

func ProvideStackDownloadManagerReal() *StackDownloadManagerReal {
	return &StackDownloadManagerReal{downloadProcessProvider: &DownloadProcessProviderReal{}}
}

func (s *StackDownloadManagerReal) GetStackDownloadStates() map[string]DownloadState {
	s.mu.Lock()
	defer s.mu.Unlock()
	downloadStateClone := make(map[string]DownloadState)
	for _, v := range s.downloadStates {
		downloadStateClone[v.stackName] = v.State
	}
	return downloadStateClone
}

func (s *StackDownloadManagerReal) DownloadStack(stackName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, downloadState := range s.downloadStates {
		if downloadState.stackName == stackName {
			downloadState.State = Ongoing
			s.downloadProcessProvider.StartDownloadProcessAndSetStateWhenFinished(downloadState)
			return
		}
	}

	newDownloadState := &StackDownloadState{stackName: stackName, State: Ongoing}
	s.downloadStates = append(s.downloadStates, newDownloadState)
	s.downloadProcessProvider.StartDownloadProcessAndSetStateWhenFinished(newDownloadState)
}

type DownloadProcessProvider interface {
	StartDownloadProcessAndSetStateWhenFinished(stackDownloadState *StackDownloadState)
}

type DownloadProcessProviderMock struct {
	stackDownloadState *StackDownloadState
}

func (d *DownloadProcessProviderMock) StartDownloadProcessAndSetStateWhenFinished(stackDownloadState *StackDownloadState) {
	d.stackDownloadState = stackDownloadState
}

type DownloadProcessProviderReal struct{}

func (d *DownloadProcessProviderReal) StartDownloadProcessAndSetStateWhenFinished(stackDownloadState *StackDownloadState) {
	go func() {
		stackDockerComposePath := StackFileDir + "/" + stackDownloadState.stackName + "/docker-compose.yml"
		pullCmd := exec.Command("docker", "compose", "-f", stackDockerComposePath, "pull")
		err := pullCmd.Run()
		if err != nil {
			Logger.Error("Error executing command '%s': %v\n", pullCmd.String(), err)
			stackDownloadState.State = Error
			return
		}

		buildCmd := exec.Command("docker", "compose", "-f", stackDockerComposePath, "build", "--pull")
		err = buildCmd.Run()
		if err == nil {
			Logger.Debug("Successfully downloaded images for stack: %s", stackDownloadState.stackName)
			stackDownloadState.State = Finished
		} else {
			Logger.Error("Error executing command '%s': %v\n", buildCmd.String(), err)
			stackDownloadState.State = Error
		}
	}()
}
