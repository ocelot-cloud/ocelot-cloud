package main

import (
	"github.com/ocelot-cloud/shared"
	"os"
)

var (
	currentSchemaVersion = "0.1.0"
	databaseFile         = shared.DataDir + "/sqlite.db" // TODO it is greyed out, but shouldn't that be used somewhere?
	Logger               shared.Logger
	cookieName           = "auth"
	profile              = getProfile()
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
