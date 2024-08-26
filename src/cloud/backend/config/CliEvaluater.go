package tools

import (
	"flag"
	"fmt"
	"github.com/ocelot-cloud/shared"
	"os"
	"strings"
)

var logger = shared.ProvideLogger()

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

	var isOidcAuthenticationDisabled bool
	flag.BoolVar(&isOidcAuthenticationDisabled, "disable-security", false, "disable security, such as authentication via OIDC")

	var useDummyStacksCliArgument bool
	flag.BoolVar(&useDummyStacksCliArgument, "enable-dummy-stacks", false, "disable security, such as authentication via OIDC")

	flag.Parse()

	var useDummyStacks = shallDummyStacksBeUsed(useDummyStacksCliArgument)

	// TODO Replace each backendMode step by step:
	if PROFILE == PROD {
		// TODO Using dummy stacks?
		return SetGlobalConfig(logLevelStr, !isOidcAuthenticationDisabled, true)
	} else if PROFILE == TEST {
		return SetGlobalConfig(logLevelStr, !isOidcAuthenticationDisabled, useDummyStacks)
	} else {
		errMsg := fmt.Sprintf("Unknown profile: %v", PROFILE)
		panic(errMsg)
	}
	// TODO Test cases to handle in ci-runner: backend mocked, backend full, frontend + backend mocked
}

func shallDummyStacksBeUsed(useDummyStacksCliArgument bool) bool {
	useDummyStacksEnvVariable := os.Getenv("USE_DUMMY_STACKS")
	return useDummyStacksCliArgument ||
		strings.ToLower(useDummyStacksEnvVariable) == "true" // TODO there should be only the ENV variable
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
		partialConfig = PartialConfig{"localhost", false, true, false}
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

	shared.LogLevel = EvaluateLogLevelBasedOn(logLevelStr)
	logger = shared.ProvideLogger()
	logGlobalConfig(config)
	return &config
}

func logGlobalConfig(config GlobalConfig) {
	logger.Info("Profile is: %s", PROFILE.String())
	logger.Info("Log level is: %s", shared.LogLevel.String())
	logger.Debug("Is web GUI enabled? -> %v", config.IsGuiEnabled)
	logger.Debug("Is security enabled? -> %v", config.IsSecurityEnabled)
	logger.Debug("Is the CORS policy relaxed by explicitly allowing cross-origin requests by setting specific response headers? -> %v", config.AreCrossOriginRequestsAllowed)
	if config.AreCrossOriginRequestsAllowed {
		logger.Warn("The CORS policy is relaxed and cross-origin requests are allowed. This is for development, so don't use this option in production!")
	}
	logger.Debug("Are mocks enabled for faster testing? -> %v", config.AreMocksEnabled)
	logger.Debug("Use dummy stacks? -> %v", config.UseDummyStacks)
}

func EvaluateLogLevelBasedOn(levelStr string) shared.LogLevelValue {
	switch strings.ToLower(levelStr) {
	case "trace":
		return shared.TRACE
	case "debug":
		return shared.DEBUG
	case "info":
		return shared.INFO
	case "warn":
		return shared.WARN
	case "error":
		return shared.ERROR
	default:
		return shared.INFO
	}
}
