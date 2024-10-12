package apps_new

import (
	"ocelot/backend/tools"
	"time"
)

var logger = tools.Logger

// TODO !!! I just leave the states away. Maybe simply have Not Available and Available. I can add new states later.
// TODO The health check should be done by ocelot -> make a simple port scan. If the port is open, the app is available.

// TODO First implement the app database and build the logic here on top of it.
// TODO At start, I need to take the stacks, put them into a tar.gz and load them in the database.
// TODO I need to create a high-level structure, which contains the AppContexts of all installed apps.

type State int

const (
	// TODO I also need a state for: "Stopped" implying there is still data. Stopped + "action: prune" -> Uninitialized
	// TODO "prune" deletes everything: image, network, volume, container
	Uninitialized State = iota
	Downloading
	Starting
	Available
	Stopping
	Error // TODO In case some operation went wrong, but how to fix it from a users perspective? maybe use stop/prune?
)

func createAppState() *AppContext {
	return &AppContext{nil, Uninitialized}
}

type AppDetails struct {
	maintainer string
	app        string
	version    string
}

// TODO For the database/docker service I need more details for each app: download domain (e.g. "hub.ocelot-cloud.com"), maintainer, and tag version in addition to the app name.
// TODO Find better names, AppDetails and AppInfo is too similar and not clear enough
type AppContext struct {
	app   *AppDetails // TODO Needed to know on which app to operate on.
	state State
	// TODO Add fields: Downloader, DockerService, YamlService; or add a single field which contains all of them combined
}

func (a *AppContext) Deploy() {
	if a.state == Uninitialized {
		a.state = Downloading
		go func() {
			time.Sleep(100 * time.Millisecond)
			a.state = Starting
			time.Sleep(100 * time.Millisecond)
			a.state = Available
		}()

		// TODO Start download
		// TODO when finished, trigger next event: startApp + state = Starting
		// TODO when finished, trigger next event: state = Available
	} else {
		// TODO
		logger.Info("action xxx was triggered. Since state is yyy, the action is ignored")
	}
}

func (a *AppContext) Stop() {
	if a.state == Uninitialized {
		// TODO Do nothing?
	} else {
		a.state = Stopping
		go func() {
			time.Sleep(100 * time.Millisecond)
			a.state = Uninitialized
		}()
		// TODO conduct stopping
		// TODO Abort download possible?
	}
}

// TODO What do I do in case Ocelot is restarted? Or when the PC is restarted? State should remain accurate. -> Save some information in the database?
