package main

import (
	"github.com/ocelot-cloud/shared"
	"os"
)

var (
	currentSchemaVersion = "0.1.0"
	databaseFile         = shared.DataDir + "/sqlite.db" // TODO it is greyed out, but shouldn't that be used somewhere?
	Logger               shared.Logger
	port                 = "8082"
	rootUrl              = "http://localhost:" + port
	cookieName           = "auth"
	profile              = getProfile()

	registrationPath = "/registration"
	loginPath        = "/login"
	logoutPath       = "/logout"
	authCheckPath    = "/auth-check"
	wipeDataPath     = "/wipe-data"

	userPath           = "/user"
	deleteUserPath     = userPath + "/delete"
	changePasswordPath = userPath + "/password"

	tagPath       = "/tags"
	tagUploadPath = tagPath + "/upload"
	tagDeletePath = tagPath + "/delete"
	getTagsPath   = tagPath + "/get-tags"
	downloadPath  = tagPath + "/"

	appPath         = "/apps"
	appCreationPath = appPath + "/create"
	appGetListPath  = appPath + "/get-list"
	appDeletePath   = appPath + "/delete"
	searchAppsPath  = appPath + "/search"
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
