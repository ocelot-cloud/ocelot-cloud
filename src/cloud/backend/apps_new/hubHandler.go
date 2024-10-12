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

func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	userAndApp, err := ReadBody[tools.UserAndApp](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tags, err := hubClient.GetTags(*userAndApp)
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

func TagDownloadHandler(w http.ResponseWriter, r *http.Request) {
	tagInfo, err := ReadBody[tools.TagInfo](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	doesTagExist := repo.AppRepo.DoesTagExist(*tagInfo)
	if !doesTagExist {
		err = DownloadTag(*tagInfo)
		if err != nil {
			Logger.Info("Failed to download tag: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func AppStartHandler(w http.ResponseWriter, r *http.Request) {
	tagInfo, err := ReadBody[tools.TagInfo](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = StartContainer(*tagInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AppStopHandler(w http.ResponseWriter, r *http.Request) {
	tagInfo, err := ReadBody[tools.TagInfo](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = StopContainer(*tagInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// TODO Database: add "apps" table column "should_be_running"
//   in-memory table: all apps that should be running are being conducted a port scan each second, depending on the success the status isAvailable is set to true or false
//   app name must be the url name, which is currently not applied, "nginxdefault" is "nginx-default" in the url
//   table "apps" should also contain the port to check

func AppReadHandler(w http.ResponseWriter, r *http.Request) {

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
