package apps

import (
	"fmt"
	"ocelot/backend/tools"
)

type dockerServiceMock struct {
	stackStates                  map[string]appState
	hasWaitedToPassDownloadState bool
}

func provideServiceMock() *dockerServiceMock {
	return &dockerServiceMock{stackStates: make(map[string]appState), hasWaitedToPassDownloadState: false}
}

func (d *dockerServiceMock) deployApp(stackName string) error {
	if stackName == "not-existing-stack" {
		return logAndCreateAppNotFoundError(stackName)
	} else if stackName == tools.NginxSlowStart || stackName == tools.NginxDownloading {
		d.stackStates[stackName] = Starting
	} else {
		d.stackStates[stackName] = Available
	}
	state := d.stackStates[stackName]
	logger.Debug("Mock pretends to have deployed stack '%s' with state %s.", stackName, state.toString())
	return nil
}

func (d *dockerServiceMock) stopApp(stackName string) error {
	if _, ok := d.stackStates[stackName]; ok {
		d.stackStates[stackName] = Uninitialized
	} else {
		return fmt.Errorf("error, stack %s does not exist in mock", stackName)
	}
	logger.Debug("Mock pretends to have stopped stack '%s'", stackName)
	return nil
}

func (d *dockerServiceMock) getRunningAppStateInfo() (map[string]appDetailsType, error) {
	logger.Trace("Mock return stack state info of virtually managed stacks")

	clonedStates := make(map[string]appDetailsType)
	for stackName, stackState := range d.stackStates {
		clonedStates[stackName] = appDetailsType{stackState, "/"}
	}

	for key, value := range d.stackStates {
		if key == tools.NginxSlowStart {
			d.stackStates[key] = Available
		} else if key == tools.NginxDownloading {
			if !d.hasWaitedToPassDownloadState {
				d.hasWaitedToPassDownloadState = true
			} else if key == tools.NginxDownloading && value == Starting {
				d.stackStates[tools.NginxDownloading] = Available
			}
		}
	}
	return clonedStates, nil
}
