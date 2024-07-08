package main

import "github.com/ocelot-cloud/shared"

var (
	dataDir      = "data"
	usersDir     = dataDir + "/users"
	databaseFile = dataDir + "/sqlite.db"

	Logger             shared.Logger
	tagPath            = "/tags"
	downloadPath       = tagPath + "/"
	userPath           = "/users"
	appPath            = "/apps"
	loginPath          = "/login"
	registrationPath   = "/registration"
	port               = "8082"
	rootUrl            = "http://localhost:" + port
	cookieName         = "auth"
	changePasswordPath = userPath + "/password"
)
