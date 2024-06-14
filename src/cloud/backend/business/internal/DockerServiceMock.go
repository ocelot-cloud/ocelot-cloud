package internal

import (
	"fmt"
	"ocelot/backend/config"
)

type DockerServiceMock struct {
	stackStates                  map[string]StackState
	hasWaitedToPassDownloadState bool
}

func ProvideServiceMock() *DockerServiceMock {
	return &DockerServiceMock{stackStates: make(map[string]StackState), hasWaitedToPassDownloadState: false}
}

func (d *DockerServiceMock) DeployStack(stackName string) error {
	if stackName == "not-existing-stack" {
		return logAndCreateStackNotFoundError(stackName)
	} else if stackName == tools.NginxSlowStart || stackName == tools.NginxDownloading {
		d.stackStates[stackName] = Starting
	} else {
		d.stackStates[stackName] = Available
	}
	state := d.stackStates[stackName]
	Logger.Debug("Mock pretends to have deployed stack '%s' with state %s.", stackName, state.String())
	return nil
}

func (d *DockerServiceMock) StopStack(stackName string) error {
	if _, ok := d.stackStates[stackName]; ok {
		d.stackStates[stackName] = Uninitialized
	} else {
		return fmt.Errorf("error, stack %s does not exist in mock", stackName)
	}
	Logger.Debug("Mock pretends to have stopped stack '%s'", stackName)
	return nil
}

func (d *DockerServiceMock) GetRunningStackStateInfo() (map[string]StackDetails, error) {
	Logger.Trace("Mock return stack state info of virtually managed stacks")

	clonedStates := make(map[string]StackDetails)
	for stackName, stackState := range d.stackStates {
		clonedStates[stackName] = StackDetails{stackState, "/"}
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
