package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var CustomCleanupFunc = func() {
	println("For your information: Custom cleanup method not set.")
}

var idsOfDaemonProcessesCreatedDuringThisRun []int

func StartDaemon(dir string, commandStr string, envs ...string) {
	var cmd *exec.Cmd
	cmd = buildCommand(dir, commandStr)

	appendEnvsToCommand(cmd, envs)

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
		ColoredPrintln("Started daemon with ID '%v' using command '%s'\n", cmd.Process.Pid, commandStr)
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
	ColoredPrintln("Cleanup method called.\n")
	killDaemonProcessesCreateDuringThisRun()
	CustomCleanupFunc()
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
}

func CleanupAndExitWithError() {
	Cleanup()
	os.Exit(1)
}

func appendEnvsToCommand(cmd *exec.Cmd, envs []string) {
	envsWithLogLevel := append(envs, DefaultEnvs...)
	cmd.Env = append(os.Environ(), envsWithLogLevel...)
}

var DefaultEnvs []string
