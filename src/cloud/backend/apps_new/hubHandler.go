package apps_new

import "net/http"

// TODO There are three kind of endpoints that must be distinguished: public (like login), user-level (like readApps), admin-level (like start/stop apps) -> should be specific functions for registration.
//   registerPublicEndpoint("path", handler), registerUserEndpoint(...), registerAdminEndpoint(...)

// TODO Re-use the approach to read dto's from requests like it was done in

func GetTagsHandler(w http.ResponseWriter, r *http.Request) {

}

func AppSearchHandler(w http.ResponseWriter, r *http.Request) {

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
