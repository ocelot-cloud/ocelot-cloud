package src

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const Scheme = "http"
const RootDomain = "localhost"

const ocelotUrl = Scheme + "://ocelot-cloud." + RootDomain
const frontendServerUrl = Scheme + "://localhost:8081"

func ExecuteInDir(dir string, commandStr string, envs ...string) {
	shortDir := strings.Replace(dir, srcDir, "", -1)
	ColoredPrint("\nIn directory '.%s', executing '%s'\n", shortDir, commandStr)

	cmd := buildCommand(dir, commandStr)
	if len(envs) != 0 {
		cmd.Env = append(os.Environ(), envs...)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutMulti := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderrMulti := io.MultiWriter(os.Stderr, &stderrBuf)
	cmd.Stdout = stdoutMulti
	cmd.Stderr = stderrMulti

	startTime := time.Now()
	err := cmd.Run()
	elapsed := time.Since(startTime)
	elapsedStr := fmt.Sprintf("%.3f", elapsed.Seconds())

	output := stdoutBuf.String() + stderrBuf.String()
	elapsedTimeSummary := fmt.Sprintf("Time taken: %s seconds.", elapsedStr)
	if err != nil {
		ColoredPrint(" => Command failed with error: %v; %s\n", err, elapsedTimeSummary)
		CleanupAndExitWithError()
	} else {
		if strings.Contains(output, "no test files") {
			ColoredPrint(" => Testing failed because no tests were found. %s\n", elapsedTimeSummary)
			CleanupAndExitWithError()
		} else if strings.Contains(commandStr, "go test") && !strings.Contains(output, "PASS:") && !containsOkLine(output) {
			ColoredPrint(" => Testing failed because no tests were actually executed; all tests were either skipped or not included. %s\n", elapsedTimeSummary)
			CleanupAndExitWithError()
		} else if strings.Contains(commandStr, "go test") && strings.Contains(output, "testing: warning: no tests to run") {
			ColoredPrint(" => Testing failed because no tests were actually executed. %s\n", elapsedTimeSummary)
			CleanupAndExitWithError()
		} else {
			ColoredPrint(" => Command successful. %s\n", elapsedTimeSummary)
		}
	}
}

func containsOkLine(output string) bool {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "ok") {
			return true
		}
	}
	return false
}

func ColoredPrint(format string, a ...interface{}) {
	colorReset := "\033[0m"
	colorCode := "\033[31m"
	fmt.Printf(colorCode+format+colorReset, a...)
}

func buildCommand(dir string, commandStr string) *exec.Cmd {
	parts, err := ParseCommand(commandStr)
	if err != nil {
		fmt.Printf("Error parsing command: %s\n", err)
		CleanupAndExitWithError()
	}
	if len(parts) == 0 {
		fmt.Println("Error: no command provided")
		CleanupAndExitWithError()
	}
	command := parts[0]
	args := parts[1:]

	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	return cmd
}

func StartBackendDaemon(profile string) {
	StartDaemon(backendDir, "./backend -enable-dummy-stacks -disable-security -log-level=debug -profile="+profile)
	WaitUntilPortIsReady("localhost:8080")
}

func WaitUntilPortIsReady(address string) {
	retryOperation(func() (bool, error) {
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			conn.Close()
			return true, nil
		}
		return false, err
	}, "Port", address, 30)
}

func retryOperation(operation func() (bool, error), description, target string, maxAttempts int) {
	attempt := 0
	for attempt < maxAttempts {
		success, err := operation()
		if success && err == nil {
			fmt.Printf("%s was requested successfully at %s\n", description, target)
			return
		} else {
			if attempt%5 == 0 {
				fmt.Printf("Attempt %v/%v: %s is not yet reachable at %s. error: %v. Trying again...\n", attempt, maxAttempts, description, target, err)
			}
			attempt++
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Printf("Error: %s could not be reached in time at %s. Cleanup and exit...\n", description, target)
	CleanupAndExitWithError()
}

func WaitForIndexPageToBeReady(url string) {
	retryOperation(func() (bool, error) {
		response, err := http.Get(url)
		if err == nil && response.StatusCode == 200 {
			return true, nil
		}
		return false, err
	}, "Index page", url, 30)
}
