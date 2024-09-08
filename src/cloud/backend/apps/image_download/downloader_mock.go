package image_download

import "ocelot/backend/tools"

type DownloaderMock struct {
	downloadStates map[string]DownloadState
}

func ProvideDownloaderMock() *DownloaderMock {
	return &DownloaderMock{make(map[string]DownloadState)}
}

func (s *DownloaderMock) GetDownloadStates() map[string]DownloadState {
	downloadStatesClone := make(map[string]DownloadState)
	for key, value := range s.downloadStates {
		downloadStatesClone[key] = value
	}
	s.updateStates()
	return downloadStatesClone
}

func (s *DownloaderMock) updateStates() {
	for key, value := range s.downloadStates {
		if value == Ongoing {
			s.downloadStates[key] = finished
		}
	}
}

func (s *DownloaderMock) Download(stackName string) {
	if stackName == tools.NginxDownloading {
		s.downloadStates[stackName] = Ongoing
	} else {
		s.downloadStates[stackName] = finished
	}

}
