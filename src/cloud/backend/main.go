package main

import (
	"github.com/gorilla/mux"
	"ocelot/backend/security"
	"ocelot/backend/setup"
	"ocelot/backend/tools"
)

// TODO Make CI pipeline running again
// TODO Make tests with real containers only when using REST API, for GUI/Acceptance tests always use the mock since it makes trouble in CI otherwise

// TODO Update "shared" module version
// TODO Consider reusing stuff from the hub, like security (potential clash with cloud package "security"), sql logic, hub client (search for apps, download, maybe upload to keep them private?)
// TODO Implement security, there should be a policy that Origin from request header == initially defined Origin as ENV variable or default ("http://localhost:8080")
// TODO refactor table: list apps with state, but make them selectable, so that there is only a single start/stop button.
// TODO Simplify profiles: DEV + PROD, no mocked frontend anymore, no security disabling anymore.
// TODO Due to implementation of the hub I can delete alls the stacks in the cloud. Acceptance tests need to integrate hub and need to implement download of stacks at the beginning? Hub should have those default files included? -> Dummies stay in cloud, sample apps like gitea go to the hub
// TODO In the end, add deploy script which only works on my device, since I have the correct SSH keys and config.
// TODO Drop the folder structure for the stacks and store everything in an sqlite. When using dummies, just load them into database at the start if not present.
// TODO ci tests should work without internet, when all required dependencies are already downloaded
//   -> I guess either docker or node is causing the issue. Testing: Disconnect internet and then try to run tests.
// TODO  GUI should take its base URL from the current URL, e.g. http://localhost:8081 when testing, so that it is flexible
// TODO get rid of the linux specific bash code in the ci-runner, replace it by native go code.
// TODO Delete the second frontend script (the one with mocked frontend) when I simplified the setup so that there are only two profiles
// TODO in ci-runner "Build(Acceptance)" etc should not be necessary. It is not intuitive when implementing new tests. I think, ExecuteInDir and StartDaemon should have an initial function like: if argumentDir == fronendDir then Build(Frontend), analogous for acceptance
// TODO Also scheduled tests can be simplified (no development profile any longer)?
// TODO In cloud is use this line "var logger = shared.ProvideLogger()". Is this maybe no longer working with the new version as I have to set it to Info by hand? -> Maybe simplify by using: ProvideLogger("DEBUG") instead.

func main() {
	setup.VerifyCliToolInstallations()
	config := tools.GenerateGlobalConfiguration()
	setup.InitializeDatabase(config)
	router := mux.NewRouter()

	security.InitializeSecurity(router, config)
	setup.InitializeApplication(router, config)
}
