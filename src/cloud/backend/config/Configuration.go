package tools

import (
	"flag"
	"fmt"
	"github.com/ocelot-cloud/shared"
	"os"
)

var logger = shared.ProvideLogger("info")

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

	// TODO get rid of all "disable-security" and "enable-dummy-stacks"

	var useDummyStacks = os.Getenv("USE_DUMMY_STACKS") == "true"

	// TODO Replace each backendMode step by step:
	if PROFILE == PROD {
		// TODO Should I only use dummy stacks in PROD? Or just real stacks?
		// TODO security/auth should always be enabled
		disableSecurity := os.Getenv("DISABLE_SECURITY") == "true"
		return SetGlobalConfig(logLevelStr, !disableSecurity, true)
	} else if PROFILE == TEST {
		return SetGlobalConfig(logLevelStr, false, useDummyStacks)
	} else {
		errMsg := fmt.Sprintf("Unknown profile: %v", PROFILE)
		panic(errMsg)
	}
	// TODO Test cases to handle in ci-runner: backend mocked, backend full, frontend + backend mocked
}

type PartialConfig struct {
	RootDomain                    string
	IsGuiEnabled                  bool
	AreCrossOriginRequestsAllowed bool
	IsOidcAuthenticationEnabled   bool
}

func SetGlobalConfig(logLevelStr string, isOidcAuthenticationEnabled bool, useDummyStacks bool) *GlobalConfig {
	config := GlobalConfig{}
	partialConfig := PartialConfig{}

	var areMocksEnabled bool
	// TODO PROD should take the root domain from ENV variable, if not present, fail
	// TODO TEST should be default localhost address
	if PROFILE == PROD {
		partialConfig = PartialConfig{"localhost", true, false, true}
		areMocksEnabled = true
	} else if PROFILE == TEST {
		partialConfig = PartialConfig{"localhost", false, true, true} // TODO both auth vars are always true, can be removed
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
		isOidcAuthenticationEnabled,
		useDummyStacks,
		"http",
		partialConfig.RootDomain,
		"8080",
	}

	logger = shared.ProvideLogger(logLevelStr)
	logGlobalConfig(config)
	return &config
}

func logGlobalConfig(config GlobalConfig) {
	logger.Info("Profile is: %s", PROFILE.String())
	logger.Info("Log level is: %s", shared.GetLogLevel())
	logger.Debug("Is web GUI enabled? -> %v", config.IsGuiEnabled)
	logger.Debug("Is security enabled? -> %v", config.IsSecurityEnabled)
	logger.Debug("Is the CORS policy relaxed by explicitly allowing cross-origin requests by setting specific response headers? -> %v", config.AreCrossOriginRequestsAllowed)
	if config.AreCrossOriginRequestsAllowed {
		logger.Warn("The CORS policy is relaxed and cross-origin requests are allowed. This is for development, so don't use this option in production!")
	}
	logger.Debug("Are mocks enabled for faster testing? -> %v", config.AreMocksEnabled)
	logger.Debug("Use dummy stacks? -> %v", config.UseDummyStacks)
}
