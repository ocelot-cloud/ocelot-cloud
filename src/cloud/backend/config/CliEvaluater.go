package tools

import (
	"flag"
	"fmt"
	"github.com/ocelot-cloud/shared"
	"os"
	"strings"
)

var logger = shared.ProvideLogger()

const (
	BackendModeProdWithGui        = "production"
	BackendModeDependenciesMocked = "dependencies-mocked"
	BackendModeDevelopmentSetup   = "development-setup"
)

type BackendComponentMode int

const (
	// ProdWithGui For Production. Secure and all features enabled.
	ProdWithGui BackendComponentMode = iota
	// DependenciesMocked For fast testing. Slow dependencies replaced by mocks.
	DependenciesMocked
	// DevelopmentSetup For GUI development, when backend should run in the background, testing via manually interactions with GUI talking to backend.
	DevelopmentSetup
)

func (s *BackendComponentMode) String() string {
	return [...]string{"ProdWithGui", "DependenciesMocked", "DevelopmentSetup"}[*s]
}

func GenerateGlobalConfiguration() *GlobalConfig {
	var backendModeString string
	profiles := fmt.Sprintf("%s, %s, %s", BackendModeProdWithGui, BackendModeDependenciesMocked, BackendModeDevelopmentSetup)
	flag.StringVar(&backendModeString, "profile", BackendModeProdWithGui, "The profile in which the application is run. Possible values: "+profiles)

	var logLevelStr string
	flag.StringVar(&logLevelStr, "log-level", "notSet", "set log level (trace, debug, info, warn, error)")

	var isOidcAuthenticationDisabled bool
	flag.BoolVar(&isOidcAuthenticationDisabled, "disable-security", false, "disable security, such as authentication via OIDC")

	var useDummyStacksCliArgument bool
	flag.BoolVar(&useDummyStacksCliArgument, "enable-dummy-stacks", false, "disable security, such as authentication via OIDC")

	flag.Parse()

	var backendMode BackendComponentMode
	if backendModeString == BackendModeProdWithGui {
		backendMode = ProdWithGui
	} else if backendModeString == BackendModeDependenciesMocked {
		backendMode = DependenciesMocked
	} else if backendModeString == BackendModeDevelopmentSetup {
		backendMode = DevelopmentSetup
	} else {
		panic("Backend mode not supported: " + backendModeString)
	}
	var useDummyStacks = shallDummyStacksBeUsed(useDummyStacksCliArgument, backendMode)

	return SetGlobalConfig(backendMode, logLevelStr, !isOidcAuthenticationDisabled, useDummyStacks)
}

func shallDummyStacksBeUsed(useDummyStacksCliArgument bool, backendMode BackendComponentMode) bool {
	useDummyStacksEnvVariable := os.Getenv("USE_DUMMY_STACKS")
	return useDummyStacksCliArgument ||
		strings.ToLower(useDummyStacksEnvVariable) == "true" ||
		backendMode == DevelopmentSetup
}

type PartialConfig struct {
	RootDomain                       string
	AreMocksEnabled                  bool
	IsGuiEnabled                     bool
	AreCrossOriginRequestsAllowed    bool
	WaitForSecurityBeforeOpeningPort bool
	IsOidcAuthenticationEnabled      bool
}

func SetGlobalConfig(backendMode BackendComponentMode, logLevelStr string, isOidcAuthenticationEnabled bool, useDummyStacks bool) *GlobalConfig {
	config := GlobalConfig{}
	partialConfig := PartialConfig{}

	if backendMode == ProdWithGui {
		partialConfig = PartialConfig{"localhost", false, true, false, true, true}
	} else if backendMode == DependenciesMocked {
		partialConfig = PartialConfig{"localhost", true, false, false, false, false}
	} else if backendMode == DevelopmentSetup {
		partialConfig = PartialConfig{"localhost", false, false, true, false, false}
	}

	config = GlobalConfig{
		partialConfig.AreCrossOriginRequestsAllowed,
		partialConfig.AreMocksEnabled,
		partialConfig.IsGuiEnabled,
		isOidcAuthenticationEnabled,
		backendMode,
		partialConfig.WaitForSecurityBeforeOpeningPort,
		useDummyStacks,
		"http",
		partialConfig.RootDomain,
		"8080",
	}

	shared.LogLevel = EvaluateLogLevelBasedOn(backendMode, logLevelStr)
	logger = shared.ProvideLogger()
	logGlobalConfig(config)
	return &config
}

func logGlobalConfig(config GlobalConfig) {
	logger.Info("Profile is: %s", config.BackendMode.String())
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

func EvaluateLogLevelBasedOn(BackendMode BackendComponentMode, levelStr string) shared.LogLevelValue {
	if levelStr == "notSet" {
		if BackendMode == ProdWithGui || BackendMode == DependenciesMocked {
			return shared.INFO
		} else if BackendMode == DevelopmentSetup {
			return shared.DEBUG
		}
	}

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
		panicMsg := fmt.Sprintf("Invalid log level: %s. Valid values are '-log-level=x' with x is one of these values: trace, debug, info (default), warn, error", levelStr)
		panic(panicMsg)
	}
}
