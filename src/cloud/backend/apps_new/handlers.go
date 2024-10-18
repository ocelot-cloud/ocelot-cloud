package apps_new

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/utils"
	"io"
	"net/http"
	"ocelot/backend/repo"
	"ocelot/backend/tools"
)

var Logger = tools.Logger

// TODO There are three kind of endpoints that must be distinguished: public (like login), user-level (like readApps), admin-level (like start/stop apps) -> should be specific functions for registration.
//   registerPublicEndpoint("path", handler), registerUserEndpoint(...), registerAdminEndpoint(...)

// TODO Return error messaged like in the hub handlers
// TODO Re-use the approach to read dto's from requests like it was done in

// TODO The handlers should require ID's where appropriate instead of long data structure, like in the repos.

// TODO Should directly use appId
func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	userAndApp, err := ReadBody[tools.UserAndApp](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	apps, err := hubClient.SearchApps(userAndApp.App)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app := (*apps)[0]

	tags, err := hubClient.GetTags(app.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	utils.SendJsonResponse(w, *tags)
}

func AppSearchHandler(w http.ResponseWriter, r *http.Request) {
	searchTermSingleString, err := ReadBody[utils.SingleString](r)
	if err != nil {
		Logger.Info("Failed to read search term: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	apps, err := hubClient.SearchApps(searchTermSingleString.Value)
	if err != nil {
		Logger.Info("Failed to search apps: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	utils.SendJsonResponse(w, *apps)
}

// TODO Should directly use tagId
func TagDownloadHandler(w http.ResponseWriter, r *http.Request) {
	tagInfo, err := ReadBody[tools.TagInfo](r) // TODO Should read TagId from request
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	apps, err := hubClient.SearchApps(tagInfo.Tag)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app := (*apps)[0]
	tags, err := hubClient.GetTags(app.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tag := (*tags)[0]

	err = DownloadTag(tag.Id)
	if err != nil {
		Logger.Info("Failed to download tag: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AppStartHandler(w http.ResponseWriter, r *http.Request) {
	singleInt, err := ReadBody[tools.SingleInt](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = StartContainer(singleInt.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AppStopHandler(w http.ResponseWriter, r *http.Request) {
	appIdStruct, err := ReadBody[tools.SingleInt](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = StopContainer(appIdStruct.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// TODO Database: add "apps" table column "should_be_running", is by default false on creation
//   in-memory table: all apps that should be running are being conducted a port scan each second, depending on the success the status isAvailable is set to true or false
//   table "apps" should also contain the port to check, is set when extracting and running the container by reading app.yml, default is 80.
//   On start, there should be a process that does the health checks?
//   "apps" table: active_tag

// TODO add memory variable: map[appId int]IsAvailable bool -> implemented empty or with just false values, make a separate go routine checking that via "docker compose ls" (see old "apps" module) which sets to "true" if container is healthy

func AppReadHandler(w http.ResponseWriter, r *http.Request) {
	apps, err := repo.AppRepo.ListApps()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var appInfos []tools.AppInfoNew
	for _, app := range apps {
		// TODO "port" and "path" must be read from the app.yml files and stored in memory. Whenever an app is started, its config must be read from zip.
		// TODO "isAvailable" must be determined via healthchecks, and stored in memory
		appInfos = append(appInfos, tools.AppInfoNew{app, "80", "/", false})
	}

	utils.SendJsonResponse(w, appInfos)
}

// TODO readAppHandler, home page -> users can only see available apps and open them, no start or stop visible or allowed by backend.
//   Home page must distinguish between users and admins.

// TODO consider putting the duplicate method in hub to the shared package
func ReadBody[T any](r *http.Request) (*T, error) {
	var result T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read request body: %w", err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return &result, nil
}
