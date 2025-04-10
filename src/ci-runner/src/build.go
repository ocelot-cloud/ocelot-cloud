package src

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
		ExecuteInDir(scriptsDir, "./build.sh")
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
		ColoredPrint(component.name + " build skipped")
	} else {
		component.build()
		component.SkipBuild = true
	}
}
