package src

import (
	"bytes"
	"fmt"
	"log"
	"ocelot/ci-runner/cli"
	"os/exec"
	"strings"
)

func CustomCleanup() {
	killPotentiallyDisturbingPreExistingComponentProcesses()
	pruneDockerToEmptySetup()
	assertThatNoProcessesSurvived()
}

func assertThatNoProcessesSurvived() {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}
	for _, line := range strings.Split(out.String(), "\n") {
		if strings.Contains(line, "./backend") {
			cli.ColoredPrintln("The backend daemon process was not killed after tests were completed.\n")
			cli.CleanupAndExitWithError()
		} else if strings.Contains(line, "vue-service") || strings.Contains(line, "vue-cli-service") {
			cli.ColoredPrintln("The frontend daemon process was not killed after tests were completed.\n")
			cli.CleanupAndExitWithError()
		}
	}
}

func killPotentiallyDisturbingPreExistingComponentProcesses() {
	potentiallyPreExistingProcesses := []string{
		"./backend",
		"vue-cli-service",
		"vue-service",
		"./hub",
		"vite",
	}
	processKillCommandTemplate := "pgrep -f %s | xargs -I %% kill -9 %%"
	var processKillCommands []string
	for _, process := range potentiallyPreExistingProcesses {
		command := fmt.Sprintf(processKillCommandTemplate, process)
		processKillCommands = append(processKillCommands, command)
	}
	runShellCommand(processKillCommands)
}

func runShellCommand(commands []string) {
	for _, command := range commands {
		_ = exec.Command("/bin/sh", "-c", command).Run()
	}
}

func pruneDockerToEmptySetup() {
	dockerPruningCommands := []string{
		"docker rm $(docker ps -a -q) -f",
		"docker network prune -f",
		"docker volume prune -a -f",
		"docker image prune -f"}
	runShellCommand(dockerPruningCommands)
}
