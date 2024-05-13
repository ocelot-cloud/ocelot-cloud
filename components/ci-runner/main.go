package main

import (
	"flag"
	"fmt"
	"net"
	"ocelot/ci-runner/src"
	"os"
	"os/exec"
	"strings"
)

// TODO In all modules, wrap the "testify" framework.
// TODO Add CLI options to skip frontend building, skip backend building, or both (maybe dont stop production server at the end to immediately rerun acceptance tests), and executing specific acceptance tests only. Fail if acceptance test not found.

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("    build                           : Build the project\n")
		fmt.Printf("    test-backend-testing-api        : Run backend testing API tests (fast)\n")
		fmt.Printf("    test-backend-component          : Run backend component tests\n")
		fmt.Printf("    test-backend-component-mocked   : Run backend component tests using mocks (fast)\n")
		fmt.Printf("    test-backend-fast               : Run all fast backend tests\n")
		fmt.Printf("    test-acceptance                 : Run acceptance tests\n")
		fmt.Printf("    test-ci                         : Run all CI tests, i.e. all tests excluding scheduled tests\n")
		fmt.Printf("    test-frontend-fast              : Runs mocked frontend without backend and acceptance tests\n")
		fmt.Printf("    test-development-setup          : Runs backend, frontend and acceptance test like a dev would do it\n")
		fmt.Printf("    run-scheduled-tests             : Runs scheduled tests, i.e. slow tests that are usually executed once a day by the CI pipeline.\n")
		fmt.Printf("    test-backend-image-download     : Testing whether the download of docker images and associated state changes work.\n")
		fmt.Printf("    run-prod          			   : Starts a server in production mode.\n")
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	src.Cleanup()
	failIfRequiredPortsAreAlreadyInUse()
	failIfThereAreExistingDockerContainers()

	command := flag.Arg(0)
	switch command {
	case "build":
		src.BuildBackendAndFrontend()
	case "test-backend-testing-api":
		src.TestBackendTestingApi()
	case "test-backend-component":
		src.TestBackendComponent()
	case "test-backend-component-mocked":
		src.TestBackendComponentMocked()
	case "test-backend-fast":
		src.TestBackendFast()
	case "test-acceptance":
		src.TestAcceptance()
	case "test-ci":
		src.TestCi()
	case "test-frontend-fast":
		src.TestFrontendFast()
	case "test-development-setup":
		src.TestComponentsInDevelopmentSetupMode()
	case "run-scheduled-tests":
		src.RunScheduledTests()
	case "test-backend-image-download":
		src.TestBackendImageDownload()
	case "run-prod":
		src.RunProduction()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		flag.Usage()
		os.Exit(1)
	}
	// Failed tests abort the program. Since the program was not aborted, we can assume that all tests passed.
	src.ColoredPrint("\nSuccess! All tests passed.\n")
}

func failIfRequiredPortsAreAlreadyInUse() {
	ports := []string{"8080", "8081"}

	for _, port := range ports {
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Printf("Error: Port %s is already in use. Exiting.\n", port)
			os.Exit(1)
		} else {
			err := listener.Close()
			if err != nil {
				fmt.Printf("Could not close listener on port %s.\n", port)
				os.Exit(1)
			}
			fmt.Printf("Port %s is available.\n", port)
		}
	}
}

func failIfThereAreExistingDockerContainers() {
	cmd := exec.Command("docker", "ps", "-a", "-q")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing docker command: %v\n", err)
		os.Exit(1)
	}
	if len(strings.TrimSpace(string(output))) > 0 {
		fmt.Println("Error: There are existing Docker containers. Please destroy them and try again.")
		os.Exit(1)
	} else {
		fmt.Println("As required for DevOps jobs, no Docker containers are deployed.")
	}
}
