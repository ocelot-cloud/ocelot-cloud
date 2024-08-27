package tools

import (
	"flag"
	"github.com/ocelot-cloud/shared"
	"os"
)

var Logger = shared.ProvideLogger(os.Getenv("LOG_LEVEL"))

// TODO Most of the stuff in this file can be deleted as soon as I simplified the PROFILE setup.
var PROFILE Profile

type Profile int

const (
	PROD Profile = iota
	TEST
)

func (p *Profile) String() string {
	return [...]string{"PROD", "TEST"}[*p]
}

type BackendComponentMode int

func GenerateGlobalConfiguration() *GlobalConfig {
	profile := os.Getenv("PROFILE")
	if profile == "TEST" {
		PROFILE = TEST
	} else {
		PROFILE = PROD
	}

	// TODO add env variables: LOG_LEVEL, USE_MOCKS, USE_DUMMY_STACKS (not sure here, do I still need this?)

	var logLevelStr string
	flag.StringVar(&logLevelStr, "log-level", "notSet", "set log level (trace, debug, info, warn, error)")

	flag.Parse()

	return SetGlobalConfig(logLevelStr)
	// TODO get rid of all "disable-security" and "enable-dummy-stacks"

	// TODO Test cases to handle in ci-runner: backend mocked, backend full, frontend + backend mocked
}

type PartialConfig struct {
	IsGuiEnabled                  bool
	AreCrossOriginRequestsAllowed bool
	UseDummyStacks                bool
	IsOidcAuthenticationEnabled   bool
}

func SetGlobalConfig(logLevelStr string) *GlobalConfig {
	// TODO Should I only use dummy stacks in PROD? Or just real stacks?
	// TODO security/auth should always be enabled

	config := GlobalConfig{}
	partialConfig := PartialConfig{}

	var useDummyStacks = os.Getenv("USE_DUMMY_STACKS") == "true"
	var areMocksEnabled bool
	// TODO PROD should take the root domain from ENV variable, if not present, fail
	// TODO TEST should be default localhost address
	if PROFILE == PROD {
		disableSecurity := os.Getenv("DISABLE_SECURITY") == "true"
		partialConfig = PartialConfig{true, false, true, !disableSecurity}
		areMocksEnabled = true
	} else if PROFILE == TEST {
		partialConfig = PartialConfig{false, true, useDummyStacks, false}
		areMocksEnabled = false
	}

	enableMocksEnv := os.Getenv("ENABLE_MOCKS")
	if enableMocksEnv != "" {
		if enableMocksEnv == "true" {
			areMocksEnabled = true
		} else if enableMocksEnv == "false" {
			areMocksEnabled = false
		}
	}

	config = GlobalConfig{
		partialConfig.AreCrossOriginRequestsAllowed,
		areMocksEnabled,
		partialConfig.IsGuiEnabled,
		partialConfig.IsOidcAuthenticationEnabled,
		partialConfig.UseDummyStacks,
		"http",
		"localhost",
		"8080",
	}

	Logger = shared.ProvideLogger(logLevelStr)
	logGlobalConfig(config)
	return &config
}

func logGlobalConfig(config GlobalConfig) {
	Logger.Info("Profile is: %s", PROFILE.String())
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
