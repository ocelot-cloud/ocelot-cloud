package apps_new

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/utils"
	"io"
	"net/http"
)

// TODO There are three kind of endpoints that must be distinguished: public (like login), user-level (like readApps), admin-level (like start/stop apps) -> should be specific functions for registration.
//   registerPublicEndpoint("path", handler), registerUserEndpoint(...), registerAdminEndpoint(...)

// TODO Re-use the approach to read dto's from requests like it was done in

func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	userAndApp, err := ReadBody[UserAndApp](r)
	if err != nil {
		return
	}
	tags, err := hubClient.GetTags(*userAndApp)
	if err != nil {
		return
	}
	utils.SendJsonResponse(w, *tags)
}

func AppSearchHandler(w http.ResponseWriter, r *http.Request) {
	searchTermSinglestring, err := ReadBody[utils.SingleString](r)
	if err != nil {
		return
	}
	apps, err := hubClient.SearchApps(searchTermSinglestring.Value)
	if err != nil {
		return
	}
	utils.SendJsonResponse(w, *apps)
}

func TagDownloadHandler(w http.ResponseWriter, r *http.Request) {

}

func AppStartHandler(w http.ResponseWriter, r *http.Request) {

}

func AppStopHandler(w http.ResponseWriter, r *http.Request) {

}

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
