package main

import (
	"fmt"
	"github.com/ocelot-cloud/task-runner"
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
	Short: "Build docker image",
	Long:  "Builds the whole project from scratch and produces a production docker image",
	Run: func(cmd *cobra.Command, args []string) {
		src.Build(src.DockerImage)
		tr.ColoredPrintln("\nSuccess! Build worked.\n")
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes processes and docker artifacts",
	Long:  "Removes processes and docker artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		tr.Cleanup()
		tr.ColoredPrintln("\nSuccess! Cleanup worked.\n")
	},
}

var hubTestTypes = map[string]func(){
	"unit":        func() { src.TestHubUnits() },
	"backend":     func() { src.TestHubBackend() },
	"acceptance":  func() { src.TestHubAcceptance() },
	"all":         func() { src.TestHubAll() },
	"persistence": func() { src.TestHubPersistence() },
}

var cloudTestTypes = map[string]func(){
	"acceptance":  func() { src.TestCloudAcceptance() },
	"frontend":    func() { src.TestFrontend() },
	"backend":     func() { src.TestBackend() },
	"integration": func() { src.TestIntegration() },
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run various tests",
	Long:  "Run different types of tests for cloud, hub, ci, or schedule.",
}

var downloadDependenciesCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads application dependencies",
	Long:  "Downloads all necessary dependencies for the application. This step must be performed once at the beginning of development to set up the environment. This process is separated from other commands so that they do not check dependencies on each run, making them faster.",
	Run: func(cmd *cobra.Command, args []string) {
		src.DownloadDependencies()
	},
}

var cloudCmd = &cobra.Command{
	Use:   "cloud [" + strings.Join(getKeys(cloudTestTypes), ", ") + "]",
	Short: "Run cloud-related tests",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputTestType := args[0]
		if _, exists := cloudTestTypes[inputTestType]; !exists {
			tr.ColoredPrintln("\nerror: unknown cloud test type: %s\n", inputTestType)
			tr.ColoredPrintln("valid args: %s\n", strings.Join(getKeys(cloudTestTypes), ", "))
			os.Exit(1)
		} else {
			cloudTestTypes[inputTestType]()
		}
		tr.ColoredPrintln("\nSuccess! Cloud tests passed.\n")
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
		tr.ColoredPrintln("\nSuccess! CI tests passed.\n")
	},
}

// TODO There are flags for every command, even though they have no effect.

var hubCmd = &cobra.Command{
	Use:   "hub [" + strings.Join(getKeys(hubTestTypes), "/") + "]",
	Short: "Run hub-related tests",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputTestType := args[0]
		if _, exists := hubTestTypes[inputTestType]; !exists {
			tr.ColoredPrintln("\nerror: unknown hub test type: %s\n", inputTestType)
			tr.ColoredPrintln("valid args: %s\n", strings.Join(getKeys(hubTestTypes), ", "))
			os.Exit(1)
		} else {
			hubTestTypes[inputTestType]()
		}
		tr.ColoredPrintln("\nSuccess! Hub tests passed.\n")
	},
}

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Run scheduled tests",
	Run: func(cmd *cobra.Command, args []string) {
		src.RunScheduledTests()
		tr.ColoredPrintln("\nSuccess! Scheduled tests passed.\n")
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the ocelot-cloud docker container",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var deployContainerProdCmd = &cobra.Command{
	Use:   "prod",
	Short: "Deploy the ocelot-cloud production container",
	Run: func(cmd *cobra.Command, args []string) {
		src.DeployContainer()
		tr.ColoredPrintln("\nSuccess! Deploy worked.\n")
	},
}

var deployContainerProdWithDummiesCmd = &cobra.Command{
	Use:   "with-dummies",
	Short: "Deploy the ocelot-cloud container with dummy stacks",
	Run: func(cmd *cobra.Command, args []string) {
		src.DeployContainerWithDummies()
		tr.ColoredPrintln("\nSuccess! Deploy worked.\n")
	},
}

func main() {
	go tr.HandleSignals()
	tr.DefaultEnvs = []string{"LOG_LEVEL=DEBUG"}
	tr.CustomCleanupFunc = src.CustomCleanup

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
	pf := rootCmd.PersistentFlags()
	pf.BoolVarP(&src.SkipBackendBuild, "skip-backend-build", "b", false, "Skip building the backend")
	pf.BoolVarP(&src.SkipFrontendBuild, "skip-frontend-build", "f", false, "Skip building the frontend")
	pf.BoolVarP(&src.SkipDockerImageBuild, "skip-docker-build", "d", false, "Skip building the Docker container, including skipping building the backend and frontend")
	pf.BoolVarP(&src.Quick, "quick", "q", false, "Quick execution, only unit tests and mocked components")

	src.ComponentBuilds[src.Backend].SkipBuild = src.SkipBackendBuild
	src.ComponentBuilds[src.Frontend].SkipBuild = src.SkipFrontendBuild
	src.ComponentBuilds[src.DockerImage].SkipBuild = src.SkipDockerImageBuild

	testCmd.AddCommand(cloudCmd, ciCmd, hubCmd, scheduleCmd)
	deployCmd.AddCommand(deployContainerProdCmd, deployContainerProdWithDummiesCmd)
	rootCmd.AddCommand(buildCmd, testCmd, deployCmd, cleanCmd, downloadDependenciesCmd)

	if shouldDoPreChecks() {
		tr.Cleanup()
		failIfRequiredPortsAreAlreadyInUse()
		failIfThereAreExistingDockerContainers()
	}

	if err := rootCmd.Execute(); err != nil {
		tr.ColoredPrintln("\nError during execution: %s\n", err.Error())
		tr.CleanupAndExitWithError()
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
