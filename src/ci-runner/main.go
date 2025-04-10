package main

// TODO Simplify cloud ci logic:
/*
	There are two modes:
		1) TEST profile, native backend deploy (allow CORS, no origin checks, mocks) and frontend run in parallel
		2) PROD profile, backend/frontend deployed via docker, no mocks by default (maybe can be enabled manually?)
	frontend: PROD and TEST profile, no mocked frontend any longer
	testing:
		1) fast backend testing: unit tests + mocked TEST backend with API tests
		2) frontend + TEST backend + cypress
		3) not sure: run API tests against PROD container?
		4) PROD with cypress
*/
// TODO Copying artifacts is not necessary. I did this initially since "Dockerfile" can only address folder below its path. But when I go to the cloud directory I can simply use docker build -f "path to docker file" and the "Dockerfile" takes its resources from there or so.

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"ocelot/ci-runner/src"
	"os"
	"os/exec"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "ci-runner",
	Short: "CI Runner CLI",
	Long:  `CI Runner CLI to build, test, and deploy projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build frontend and backend",
	Long:  "Builds frontend and backend and can detect possible compilation errors.",
	Run: func(cmd *cobra.Command, args []string) {
		src.BuildBackendAndFrontend()
		src.ColoredPrint("\nSuccess! Build worked.\n")
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes processes and docker artifacts",
	Long:  "Removes processes and docker artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		src.Cleanup()
		src.ColoredPrint("\nSuccess! Cleanup worked.\n")
	},
}

var cloudTestTypes = map[string]func(){
	"backend":    func() { src.TestBackendComponent(src.Quick) },
	"frontend":   func() { src.TestCloudFrontendFast() },
	"acceptance": func() { src.TestCloudAcceptance() },
	"all":        func() { src.TestCloudAll() },
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run various tests",
	Long:  "Run different types of tests for cloud, hub, ci, or schedule.",
}

var cloudCmd = &cobra.Command{
	Use:   "cloud [backend/frontend/acceptance]",
	Short: "Run cloud-related tests",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputTestType := args[0]
		if _, exists := cloudTestTypes[inputTestType]; !exists {
			src.ColoredPrint("\nerror: unknown cloud test type: %s\n", inputTestType)
			src.ColoredPrint("valid args: %s\n", strings.Join(getKeys(cloudTestTypes), ", "))
			os.Exit(1)
		} else {
			cloudTestTypes[inputTestType]()
		}
		src.ColoredPrint("\nSuccess! Cloud tests passed.\n")
	},
}

func getKeys(m map[string]func()) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Run CI-related tests",
	Run: func(cmd *cobra.Command, args []string) {
		src.TestCi()
		src.ColoredPrint("\nSuccess! CI tests passed.\n")
	},
}

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Run scheduled tests",
	Run: func(cmd *cobra.Command, args []string) {
		src.RunScheduledTests()
		src.ColoredPrint("\nSuccess! Scheduled tests passed.\n")
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the server",
	Long:  `Starts the server in production mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		src.DeployLocally()
		src.ColoredPrint("\nSuccess! Deploy worked.\n")
	},
}

func main() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
	pf := rootCmd.PersistentFlags()
	pf.BoolVarP(&src.SkipBackendBuild, "skip-backend-build", "b", false, "Skip building the backend")
	pf.BoolVarP(&src.SkipFrontendBuild, "skip-frontend-build", "f", false, "Skip building the frontend")
	pf.BoolVarP(&src.SkipDockerImageBuild, "skip-docker-build", "d", false, "Skip building the Docker container, including skipping building the backend and frontend")
	pf.BoolVarP(&src.Quick, "quick", "q", false, "Quick execution, only unit tests and mocked components")

	src.ComponentBuilds[src.Backend].SkipBuild = src.SkipBackendBuild
	src.ComponentBuilds[src.Frontend].SkipBuild = src.SkipFrontendBuild
	src.ComponentBuilds[src.DockerImage].SkipBuild = src.SkipDockerImageBuild

	testCmd.AddCommand(cloudCmd)
	testCmd.AddCommand(ciCmd)
	testCmd.AddCommand(scheduleCmd)

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(cleanCmd)

	if shouldDoPreChecks() {
		src.Cleanup()
		failIfRequiredPortsAreAlreadyInUse()
		failIfThereAreExistingDockerContainers()
	}

	if err := rootCmd.Execute(); err != nil {
		src.ColoredPrint("\nError during execution: %s\n", err.Error())
		src.CleanupAndExitWithError()
	}
}

func shouldDoPreChecks() bool {
	if len(os.Args) == 1 {
		return false
	} else if len(os.Args) > 1 && (os.Args[1] == "completion" || os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help") {
		return false
	} else {
		return true
	}
}

func failIfRequiredPortsAreAlreadyInUse() {
	ports := []string{"8080", "8081", "8082"}

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
