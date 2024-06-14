package internal

import "ocelot/backend/config"

type StackDownloadManagerMock struct {
	downloadStates map[string]DownloadState
}

func ProvideDownloadManagerMock() *StackDownloadManagerMock {
	return &StackDownloadManagerMock{make(map[string]DownloadState)}
}

func (s *StackDownloadManagerMock) GetStackDownloadStates() map[string]DownloadState {
	downloadStatesClone := make(map[string]DownloadState)
	for key, value := range s.downloadStates {
		downloadStatesClone[key] = value
	}
	s.updateStates()
	return downloadStatesClone
}

func (s *StackDownloadManagerMock) updateStates() {
	for key, value := range s.downloadStates {
		if value == Ongoing {
			s.downloadStates[key] = Finished
		}
	}
}

func (s *StackDownloadManagerMock) DownloadStack(stackName string) {
	if stackName == tools.NginxDownloading {
		s.downloadStates[stackName] = Ongoing
	} else {
		s.downloadStates[stackName] = Finished
	}

}
