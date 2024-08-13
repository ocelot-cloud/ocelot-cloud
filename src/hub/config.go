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
	getTagsPath        = tagPath + "/get-tags"
	downloadPath       = tagPath + "/"
	userPath           = "/user"
	logoutPath         = "/logout"
	appPath            = "/apps"
	searchAppsPath     = "/apps/search"
	loginPath          = "/login"
	registrationPath   = "/registration"
	wipeDataPath       = "/wipe-data"
	authCheckPath      = "/auth-check"
	port               = "8082"
	rootUrl            = "http://localhost:" + port
	cookieName         = "auth"
	changePasswordPath = userPath + "/password"
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
