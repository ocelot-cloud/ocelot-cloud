package tools

import (
	"github.com/ocelot-cloud/shared"
	"os"
)

var Logger = shared.ProvideLogger(os.Getenv("LOG_LEVEL"))

type BackendProfile int

const (
	PROD BackendProfile = iota
	TEST
)

func (p *BackendProfile) String() string {
	return [...]string{"PROD", "TEST"}[*p]
}

type BackendComponentMode int

func GenerateGlobalConfiguration() *GlobalConfig {
	Profile := getProfile()
	config := getGlobalConfigBasedOnProfile(Profile)
	logGlobalConfig(config, Profile)
	return &config
}

func getProfile() BackendProfile {
	profile := os.Getenv("PROFILE")
	if profile == "TEST" {
		return TEST
	} else {
		return PROD
	}
}

func getGlobalConfigBasedOnProfile(profile BackendProfile) GlobalConfig {
	config := GlobalConfig{}
	// TODO PROD should take the root domain from ENV variable, if not present, fail
	config.Scheme = "http"
	config.RootDomain = "localhost"
	config.Port = "8080"

	config.UseDummyStacks = os.Getenv("USE_DUMMY_STACKS") == "true"
	config.AreMocksEnabled = os.Getenv("ENABLE_MOCKS") == "true"

	// TODO security/auth should always be enabled. Remove "IsSecurityEnabled" and "DISABLE_SECURITY" from everywhere in the project.
	if profile == PROD {
		config.IsGuiEnabled = true
		config.AreCrossOriginRequestsAllowed = false
		config.IsSecurityEnabled = os.Getenv("DISABLE_SECURITY") != "true"
	} else {
		config.IsGuiEnabled = false
		config.AreCrossOriginRequestsAllowed = true
		config.IsSecurityEnabled = false
	}
	return config
}

func logGlobalConfig(config GlobalConfig, profile BackendProfile) {
	Logger.Info("Profile is: %s", profile.String())
	Logger.Info("Log level is: %s", shared.GetLogLevel())
	Logger.Debug("Is web GUI enabled? -> %v", config.IsGuiEnabled)
	Logger.Debug("Is security enabled? -> %v", config.IsSecurityEnabled)
	Logger.Debug("Is the CORS policy relaxed by explicitly allowing cross-origin requests by setting specific response headers? -> %v", config.AreCrossOriginRequestsAllowed)
	if config.AreCrossOriginRequestsAllowed {
		Logger.Warn("The CORS policy is relaxed and cross-origin requests are allowed. This is for development, so don't use this option in production!")
	}
	Logger.Debug("Are mocks enabled for faster testing? -> %v", config.AreMocksEnabled)
	Logger.Debug("Use dummy stacks? -> %v", config.UseDummyStacks)
}
