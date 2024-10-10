package src

import (
	"fmt"
	"ocelot/ci-runner/cli"
)

var (
	SkipBackendBuild     bool
	SkipFrontendBuild    bool
	SkipDockerImageBuild bool
	Quick                bool
)

type COMPONENT int

const (
	Backend COMPONENT = iota
	Frontend
	DockerImage
	Acceptance
)

type Component struct {
	name      string
	SkipBuild bool
	build     func()
}

var ComponentBuilds = map[COMPONENT]*Component{
	Backend: {"backend", false, func() {
		cli.ExecuteInDir(backendDir, "go build")
	}},
	Frontend: {"frontend", false, func() {
		cli.ExecuteInDir(frontendDir, "npm run build")
	}},
	DockerImage: {"docker image", false, func() {
		// The flags make it executable in Docker containers
		cli.ExecuteInDir(backendDir, "go build -ldflags '-extldflags \"-static\"'")
		cli.ExecuteInDir(frontendDir, "npm run build")
		cli.ExecuteInDir(projectDir, "docker rm -f ocelotcloud/ocelotcloud")
		cli.ExecuteInDir(projectDir, "bash -c 'docker network create ocelot-net || true'")
		cli.ExecuteInDir(projectDir, "bash -c 'if [ -z \"$(docker images -q alpine:3.18.6)\" ]; then docker pull alpine:3.18.6; fi'")
		cmd := fmt.Sprintf("docker build -t ocelotcloud/ocelotcloud:local -f src/ci-runner/Dockerfile .")
		cli.ExecuteInDir(projectDir, cmd)
	}},
}

func Build(comp COMPONENT) {
	if SkipBackendBuild {
		ComponentBuilds[Backend].SkipBuild = true
	}
	if SkipFrontendBuild {
		ComponentBuilds[Frontend].SkipBuild = true
	}
	if SkipDockerImageBuild {
		ComponentBuilds[Backend].SkipBuild = true
		ComponentBuilds[Frontend].SkipBuild = true
		ComponentBuilds[DockerImage].SkipBuild = true
	}
	component := ComponentBuilds[comp]
	if component.SkipBuild {
		cli.ColoredPrintln(component.name + " build skipped")
	} else {
		component.build()
		component.SkipBuild = true
	}
}

func DownloadDependencies() {
	cli.PrintTaskDescription("downloading dependencies")
	cli.ExecuteInDir(acceptanceTestsDir, "npm install")
	cli.ExecuteInDir(frontendDir, "npm install")
	cli.ExecuteInDir(backendDir, "go mod tidy")
	cli.ExecuteInDir(hubDir, "go mod tidy")
}
