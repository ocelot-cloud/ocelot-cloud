package internal

import "os"

type StackStateService struct {
	stacks map[string]StackState
}

func ProvideStackStateService() *StackStateService {
	service := StackStateService{}
	stackNames, _ := stackNamesInDirectory()
	for _, v := range stackNames {
		service.stacks[v] = Uninitialized
	}
	return &service
}

func stackNamesInDirectory() ([]string, error) {
	files, err := os.ReadDir(StackFileDir)
	if err != nil {
		Logger.Warn("Could not read stack from directory '" + StackFileDir + "': " + err.Error())
		return nil, err
	}

	var stackNames []string
	for _, f := range files {
		if f.IsDir() {
			stackNames = append(stackNames, f.Name())
		}
	}
	return stackNames, nil
}

func DeployStack(stackName string) error {
	// if deploy is running right now, then abort second deploy, allow stopping however
	// start downloading
	// state = downloading
	// wait until download finished
	// start stack
	// state = Starting
	// wait until healthy
	// state = available
	return nil
}
