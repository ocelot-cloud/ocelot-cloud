package tools

import (
	"fmt"
	"github.com/ocelot-cloud/shared"
	"net/url"
	"os"
	"strings"
)

const CookieName = "ocelot-auth"

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
	hostParams, err := getHostParams(profile, os.Getenv("HOST"))
	if err != nil {
		Logger.Fatal("Failed to get host params: %v", err)
		// TODO Give examples/explanation how to fix it. "https://my-domain.com"
	}
	config.HttpScheme = hostParams.Scheme
	config.RootDomain = hostParams.Domain
	config.BackendPort = hostParams.Port
	config.BackendExecutablePort = "8080"

	config.UseDummyStacks = os.Getenv("USE_DUMMY_STACKS") == "true"

	// TODO security/auth should always be enabled. Remove "IsSecurityEnabled" from everywhere in the project.
	if profile == PROD {
		config.IsGuiEnabled = true
		config.UseRealDatabase = true
		config.AreMocksEnabled = false
		config.AreCrossOriginRequestsAllowed = false
		config.CreateDefaultAdminUser = false
		config.OpenDataWipeEndpoint = os.Getenv("ENABLE_DATA_WIPE_ENDPOINT") == "true"
	} else {
		config.IsGuiEnabled = false
		config.UseRealDatabase = false
		config.AreMocksEnabled = true
		config.AreCrossOriginRequestsAllowed = true
		config.CreateDefaultAdminUser = true
		config.OpenDataWipeEndpoint = true
		Logger = shared.ProvideLogger("DEBUG")
	}
	return config
}

type HostParams struct {
	Scheme string
	Domain string
	Port   string
}

func getHostParams(profile BackendProfile, hostEnv string) (*HostParams, error) {
	if profile == PROD {
		if hostEnv == "" {
			return nil, fmt.Errorf("HOST environment variable is not set")
		}

		host, err := url.Parse(hostEnv)
		if err != nil || host == nil || !host.IsAbs() || host.Path != "" || host.Host == "" {
			return nil, fmt.Errorf("invalid HOST URL: %s", host)
		}

		var port string
		if host.Port() == "" {
			if host.Scheme == "http" {
				port = "80"
			} else if host.Scheme == "https" {
				port = "443"
			} else {
				return nil, fmt.Errorf("error when evaluating port from HOST env variable")
			}
		} else {
			port = host.Port()
		}

		domain := strings.Split(host.Host, ":")[0]
		return &HostParams{host.Scheme, domain, port}, nil
	} else {
		return &HostParams{"http", "localhost", "8080"}, nil
	}

}

func logGlobalConfig(config GlobalConfig, profile BackendProfile) {
	if profile == PROD {
		Logger.Info("Profile is: %s", profile.String())
	} else {
		Logger.Warn("Profile is: %s. It is intended for development, so do not use this profile in production!", profile.String())
	}
	Logger.Info("Log level is: %s", shared.GetLogLevel())
	Logger.Debug("Is web GUI enabled? -> %v", config.IsGuiEnabled)
	Logger.Debug("Create default admin user? -> %v", config.CreateDefaultAdminUser)
	Logger.Debug("Is the CORS policy relaxed by explicitly allowing cross-origin requests by setting specific response headers? -> %v", config.AreCrossOriginRequestsAllowed)
	if config.AreCrossOriginRequestsAllowed {
		Logger.Warn("The CORS policy is relaxed and cross-origin requests are allowed.")
	}
	Logger.Debug("Are mocks enabled? -> %v", config.AreMocksEnabled)
	Logger.Debug("Use dummy stacks? -> %v", config.UseDummyStacks)
	if !config.UseRealDatabase {
		Logger.Warn("An in-memory database is used. No data is stored persistently.")
	}
}
