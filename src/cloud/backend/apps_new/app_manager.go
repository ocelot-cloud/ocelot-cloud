package apps_new

import (
	"ocelot/backend/tools"
	"time"
)

var logger = tools.Logger

// TODO First implement the app database and build the logic here on top of it.
// TODO At start, I need to take the stacks, put them into a tar.gz and load them in the database.
// TODO I need to create a high-level structure, which contains the AppContexts of all installed apps.

type State int

const (
	Uninitialized State = iota
	Downloading
	Starting
	Available
	Stopping
	Error // TODO In case some operation went wrong
)

func createAppState() *AppContext {
	return &AppContext{nil, Uninitialized}
}

type AppDetails struct {
	maintainer string
	app        string
	version    string
}

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
