package src

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var idsOfDaemonProcessesCreatedDuringThisRun []int

func StartDaemon(dir string, commandStr string, envs ...string) {
	var cmd *exec.Cmd
	cmd = buildCommand(dir, commandStr)
	if len(envs) != 0 {
		cmd.Env = append(os.Environ(), envs...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if cmd.Process == nil {
		fmt.Printf("Error: The process was not able to start properly.\n")
	} else {
		idsOfDaemonProcessesCreatedDuringThisRun = append(idsOfDaemonProcessesCreatedDuringThisRun, cmd.Process.Pid)
	}

	if err != nil {
		fmt.Printf("Command: '%s' -> failed with error: %v\n", commandStr, err)
		CleanupAndExitWithError()
	} else {
		ColoredPrint("Started daemon with ID '%v' using command '%s'\n", cmd.Process.Pid, commandStr)
		go func() {
			err := cmd.Wait()
			if err != nil {
				fmt.Printf("Command: '%s' -> reason of stopping: %v\n", commandStr, err)
			} else {
				fmt.Printf("Command: '%s' -> stopped through casual termination\n", commandStr)
			}
		}()
	}
}

func Cleanup() {
	ColoredPrint("Cleanup method called.\n")
	killDaemonProcessesCreateDuringThisRun()
	killPotentiallyDisturbingPreExistingComponentProcesses()
	pruneDockerToEmptySetup()
	fmt.Print("\x1b[?25h") // Ensure CLI cursor is visible
	fmt.Print("\x1b[0m")   // Resets all CLI cursor attributes such as color
}

func killDaemonProcessesCreateDuringThisRun() {
	println("Killing daemon processes")
	if len(idsOfDaemonProcessesCreatedDuringThisRun) == 0 {
		fmt.Println("  No daemon processes to kill.")
		return
	}

	for _, processID := range idsOfDaemonProcessesCreatedDuringThisRun {
		fmt.Printf("  Killing process with ID '%v'\n", processID)
		processGroupID, err := syscall.Getpgid(processID)
		if err != nil {
			log.Fatalf("Failed to get process group ID of process ID '%v' because of error: %v\n", processID, err)
		}
		if err := syscall.Kill(-processGroupID, syscall.SIGKILL); err != nil {
			log.Fatalf("Failed to kill process group ID '%v' because of error: %v of \n", processID, err)
		}
	}
	idsOfDaemonProcessesCreatedDuringThisRun = make([]int, 0)
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
			ColoredPrint("The backend daemon process was not killed after tests were completed.\n")
			CleanupAndExitWithError()
		} else if strings.Contains(line, "vue-service") || strings.Contains(line, "vue-cli-service") {
			ColoredPrint("The frontend daemon process was not killed after tests were completed.\n")
			CleanupAndExitWithError()
		}
	}
}

func killPotentiallyDisturbingPreExistingComponentProcesses() {
	potentiallyPreExistingProcesses := []string{
		"./backend",
		"vue-cli-service",
		"vue-service",
		"./hub",
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

func CleanupAndExitWithError() {
	Cleanup()
	os.Exit(1)
}
