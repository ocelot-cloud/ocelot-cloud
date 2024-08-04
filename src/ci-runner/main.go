package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"ocelot/ci-runner/src"
	"os"
	"os/exec"
	"strings"
)

// TODO Add option of executing specific acceptance tests only. Maybe add command which lists available options. Fail if acceptance test not found.
// TODO Try to stop a potentially running "./hub" process when cleaning up at the beginning.

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

// TODO abstract duplication: make a structure with name + func, e.g. {("backend", src.TestBackendComponent(src.Quick)), ("frontend", ...}
var validTestTypes = []string{"backend", "frontend", "acceptance", "ci", "scheduled", "hub"}
var testTypeStr = strings.Join(validTestTypes, ", ")

var testCmd = &cobra.Command{
	Use:   "test [test-type]",
	Short: "Run various tests",
	Long:  "Run different types of tests: " + testTypeStr,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputTestType := args[0]
		switch inputTestType {
		case "backend":
			src.TestBackendComponent(src.Quick)
		case "frontend":
			src.TestFrontendFast()
		case "acceptance":
			src.TestAcceptance()
		case "ci":
			src.TestCi()
		case "scheduled":
			src.RunScheduledTests()
		case "hub":
			src.TestHub()
		case "hub-acceptance":
			src.TestHubAcceptance()
		case "hub-all":
			src.TestHubAll()
		default:
			src.ColoredPrint("\nerror: unknown command: %s\n", inputTestType)
			src.ColoredPrint("valid args: %s\n", testTypeStr)
			os.Exit(1)
		}
		src.ColoredPrint("\nSuccess! All tests passed.\n")
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
