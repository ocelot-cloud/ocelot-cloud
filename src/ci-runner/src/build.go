package src

import (
	"fmt"
	"github.com/ocelot-cloud/task-runner"
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
		tr.ExecuteInDir(backendDir, "go build")
	}},
	Frontend: {"frontend", false, func() {
		tr.ExecuteInDir(frontendDir, "npm run build")
	}},
	DockerImage: {"docker image", false, func() {
		// The flags make it executable in Docker containers
		tr.ExecuteInDir(backendDir, "go build -ldflags '-extldflags \"-static\"'")
		tr.ExecuteInDir(frontendDir, "npm run build")
		tr.ExecuteInDir(projectDir, "docker rm -f ocelotcloud/ocelotcloud")
		tr.ExecuteInDir(projectDir, "bash -c 'docker network create ocelot-net || true'")
		tr.ExecuteInDir(projectDir, "bash -c 'if [ -z \"$(docker images -q alpine:3.18.6)\" ]; then docker pull alpine:3.18.6; fi'")
		cmd := fmt.Sprintf("docker build -t ocelotcloud/ocelotcloud:local -f src/ci-runner/Dockerfile .")
		tr.ExecuteInDir(projectDir, cmd)
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
		tr.ColoredPrintln(component.name + " build skipped")
	} else {
		component.build()
		component.SkipBuild = true
	}
}

func DownloadDependencies() {
	tr.PrintTaskDescription("downloading dependencies")
	tr.ExecuteInDir(acceptanceTestsDir, "npm install")
	tr.ExecuteInDir(frontendDir, "npm install")
	tr.ExecuteInDir(backendDir, "go mod tidy")
	tr.ExecuteInDir(hubDir, "go mod tidy")
}
