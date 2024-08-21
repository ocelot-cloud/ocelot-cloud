package main

import (
	"github.com/ocelot-cloud/shared"
	"os"
)

var (
	currentSchemaVersion = "0.1.0"
	databaseFile         = shared.DataDir + "/sqlite.db"

	Logger             shared.Logger
	tagPath            = "/tags"
	tagUploadPath      = tagPath + "/upload"
	tagDeletePath      = tagPath + "/delete"
	getTagsPath        = tagPath + "/get-tags"
	downloadPath       = tagPath + "/"
	logoutPath         = "/logout"
	appPath            = "/apps"
	appCreationPath    = appPath + "/create"
	appGetListPath     = appPath + "/get-list"
	appDeletePath      = appPath + "/delete"
	searchAppsPath     = appPath + "/search"
	wipeDataPath       = "/wipe-data"
	authCheckPath      = "/auth-check"
	port               = "8082"
	rootUrl            = "http://localhost:" + port
	cookieName         = "auth"
	userPath           = "/user"
	deleteUserPath     = userPath + "/delete"
	changePasswordPath = userPath + "/password"
	loginPath          = "/login"
	registrationPath   = "/registration"
	profile            = getProfile()
)

type PROFILE int

const (
	PROD PROFILE = iota
	TEST
)

func getProfile() PROFILE {
	envProfile := os.Getenv("PROFILE")
	if envProfile == "TEST" {
		return TEST
	} else {
		return PROD
	}
}
