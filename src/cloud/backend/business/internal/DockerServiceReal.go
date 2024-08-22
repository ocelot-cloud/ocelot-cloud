package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TODO Run initial test, either "docker compose" or "docker-compose" must be installed. If not, exit. If one is installed, set it globally as dockerComposeCommand or so

type DockerServiceReal struct{}

func (d *DockerServiceReal) DeployStack(stackName string) error {
	cmdPath := getStackPath(stackName)

	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		return logAndCreateStackNotFoundError(stackName)
	}

	networkCreationBashCmd := fmt.Sprintf("docker network ls | grep -q %s-net || docker network create %s-net", stackName, stackName)
	_ = exec.Command("/bin/sh", "-c", networkCreationBashCmd).Run()

	stackDeployCmd := exec.Command("docker", "compose", "-f", cmdPath, "-p", stackName, "up", "-d")
	output, err := stackDeployCmd.CombinedOutput()
	if err != nil {
		Logger.Warn("failed to deploy stack: %v, Output: %s", err, string(output))
		return fmt.Errorf("failed stack deployment")
	} else {
		Logger.Debug("Docker service deployed stack '%s'", stackName)
		return nil
	}
}

func logAndCreateStackNotFoundError(stackName string) error {
	errorMessage := "Could not find stack: " + stackName
	Logger.Error(errorMessage)
	return fmt.Errorf(errorMessage)
}

func getStackPath(stackName string) string {
	if stackName == "ocelot-cloud" {
		return fmt.Sprintf("%s/%s/docker-compose.yml", CoreStackFileDir, stackName)
	} else {
		return fmt.Sprintf("%s/%s/docker-compose.yml", StackFileDir, stackName)
	}
}

func (d *DockerServiceReal) StopStack(stackName string) error {
	configPath := getStackPath(stackName)
	cmd := exec.Command("docker", "compose", "-p", stackName, "-f", configPath, "down")
	output, err := cmd.CombinedOutput()
	if err != nil {
		Logger.Error("Command '%s' failed to stop stack: %v, Output: %s", cmd.String(), err, output)
		return fmt.Errorf("stack stopping error")
	} else {
		Logger.Debug("Docker service stopped stack '%s'", stackName)
		return nil
	}
}

func (d *DockerServiceReal) GetRunningStackStateInfo() (map[string]StackDetails, error) {
	lines, err := getDockerComposeListLines()
	if err != nil {
		Logger.Error("error, 'docker compose' command seemed not to have worked properly: %s", err.Error())
		return nil, err
	}
	genericRunningStateStacksInfo := extractNamesOfRunningStacksFromLines(lines)
	fullStacksInfoWithMoreSpecificHealthState := setHealthStates(genericRunningStateStacksInfo)

	return fullStacksInfoWithMoreSpecificHealthState, nil
}

func setHealthStates(stackStateInfo map[string]StackDetails) map[string]StackDetails {
	resultInfo := make(map[string]StackDetails)
	for stackName, stackDetail := range stackStateInfo {
		if stackDetail.State == Running {
			stackDetail.State = getHealthStateOf(stackName)
		}
		resultInfo[stackName] = stackDetail
	}
	return resultInfo
}

func getHealthStateOf(stackName string) StackState {
	if areAllStackContainersWithHealthChecksReallyHealthy(stackName) {
		return Available
	} else {
		return Starting
	}
}

func areAllStackContainersWithHealthChecksReallyHealthy(stackName string) bool {
	dockerComposeYamlPathOfStack := getStackPath(stackName)
	stackInfoCmd := exec.Command("docker", "compose", "-f", dockerComposeYamlPathOfStack, "ps")
	var out bytes.Buffer
	stackInfoCmd.Stdout = &out
	err := stackInfoCmd.Run()
	if err != nil {
		Logger.Error("Failed to read CLI output of stack specific container info for stack '%s'", stackName)
		return false
	}

	scanner := bufio.NewScanner(&out)
	skipHeader := true
	for scanner.Scan() {
		line := scanner.Text()
		if skipHeader {
			skipHeader = false
			continue
		}
		if strings.Contains(line, "(health: starting)") || strings.Contains(line, "(unhealthy)") {
			return false
		}
	}
	if err := scanner.Err(); err != nil {
		Logger.Error("Scanner error occurred for stack '%s'", stackName)
		return false
	}
	return true
}

func getDockerComposeListLines() ([]string, error) {
	cmd := exec.Command("docker", "compose", "ls", "-a")
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		Logger.Error("Command '%s' did not work: %v. Maybe the wrong version is used.", cmd.String(), err)
		versionOutputBytes, versionErr := exec.Command("docker", "compose", "version").CombinedOutput()
		if versionErr == nil {
			Logger.Error("Docker Compose version is: %s", string(versionOutputBytes))
		}
		return nil, err
	}

	output := string(outputBytes)
	lines := strings.Split(output, "\n")
	return lines, nil
}

func extractNamesOfRunningStacksFromLines(lines []string) map[string]StackDetails {
	var resultInfos = make(map[string]StackDetails)
	for _, line := range lines {
		if isHeaderOrEmpty(line) {
			continue
		}
		fields := strings.Fields(line)
		stackName, stackDetail := transformToStackStackInfo(fields)
		resultInfos[stackName] = stackDetail
	}
	return resultInfos
}

func isHeaderOrEmpty(line string) bool {
	return strings.HasPrefix(line, "NAME") || len(strings.TrimSpace(line)) == 0
}

func transformToStackStackInfo(fields []string) (string, StackDetails) {
	name := fields[0]
	rawStatus := fields[1]
	var status StackState
	if strings.Contains(rawStatus, "running") {
		status = Running
	} else {
		status = Uninitialized
	}
	return name, StackDetails{status, "/"}
}
