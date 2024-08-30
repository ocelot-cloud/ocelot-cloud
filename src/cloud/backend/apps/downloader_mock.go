package apps

import "ocelot/backend/tools"

type downloaderMock struct {
	downloadStates map[string]downloadState
}

func provideDownloaderMock() *downloaderMock {
	return &downloaderMock{make(map[string]downloadState)}
}

func (s *downloaderMock) getDownloadStates() map[string]downloadState {
	downloadStatesClone := make(map[string]downloadState)
	for key, value := range s.downloadStates {
		downloadStatesClone[key] = value
	}
	s.updateStates()
	return downloadStatesClone
}

func (s *downloaderMock) updateStates() {
	for key, value := range s.downloadStates {
		if value == ongoing {
			s.downloadStates[key] = finished
		}
	}
}

func (s *downloaderMock) download(stackName string) {
	if stackName == tools.NginxDownloading {
		s.downloadStates[stackName] = ongoing
	} else {
		s.downloadStates[stackName] = finished
	}

}
