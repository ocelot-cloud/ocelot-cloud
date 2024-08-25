package src

import "fmt"

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
		ExecuteInDir(backendDir, "go build")
	}},
	Frontend: {"backend", false, func() {
		ExecuteInDir(frontendDir, "npm install")
		ExecuteInDir(frontendDir, "npm run build")
	}},
	DockerImage: {"docker image", false, func() {
		// The flags make it executable in Docker containers
		ExecuteInDir(backendDir, "go build -ldflags '-extldflags \"-static\"'")
		// TODO Consider installing "vite"?
		ExecuteInDir(frontendDir, "npm install")
		ExecuteInDir(frontendDir, "npm run build")
		ExecuteInDir(projectDir, "docker rm -f ocelotcloud/ocelotcloud")
		ExecuteInDir(projectDir, "bash -c 'docker network create ocelot-net || true'")
		ExecuteInDir(projectDir, "bash -c 'if [ -z \"$(docker images -q alpine:3.18.6)\" ]; then docker pull alpine:3.18.6; fi'")
		cmd := fmt.Sprintf("docker build -t ocelotcloud/ocelotcloud:local -f src/ci-runner/Dockerfile .")
		ExecuteInDir(projectDir, cmd)
	}},
	Acceptance: {"backend", false, func() {
		ExecuteInDir(acceptanceTestsDir, "npm install")
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
		ColoredPrintln(component.name + " build skipped")
	} else {
		component.build()
		component.SkipBuild = true
	}
}
