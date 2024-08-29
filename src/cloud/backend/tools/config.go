package tools

import (
	"fmt"
	"github.com/ocelot-cloud/shared"
	"net/url"
	"os"
	"strings"
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
	hostParams, err := getHostParams(profile, os.Getenv("HOST"))
	if err != nil {
		Logger.Fatal("Failed to get host params: %v", err)
		// TODO Give examples/explanation how to fix it. "https://my-domain.com"
	}
	config.HttpScheme = hostParams.Scheme
	config.RootDomain = hostParams.Domain
	config.DockerContainerPort = hostParams.Port
	config.BackendExecutablePort = "8080"

	config.UseDummyStacks = os.Getenv("USE_DUMMY_STACKS") == "true"

	// TODO security/auth should always be enabled. Remove "IsSecurityEnabled" and "DISABLE_SECURITY" from everywhere in the project.
	if profile == PROD {
		config.IsGuiEnabled = true
		config.AreCrossOriginRequestsAllowed = false
		config.IsSecurityEnabled = os.Getenv("DISABLE_SECURITY") != "true"
		config.AreMocksEnabled = false
	} else {
		config.IsGuiEnabled = false
		config.AreCrossOriginRequestsAllowed = true
		config.IsSecurityEnabled = false
		config.AreMocksEnabled = true
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
	Logger.Info("Profile is: %s", profile.String())
	Logger.Info("Log level is: %s", shared.GetLogLevel())
	Logger.Debug("Is web GUI enabled? -> %v", config.IsGuiEnabled)
	Logger.Debug("Is security enabled? -> %v", config.IsSecurityEnabled)
	Logger.Debug("Is the CORS policy relaxed by explicitly allowing cross-origin requests by setting specific response headers? -> %v", config.AreCrossOriginRequestsAllowed)
	if config.AreCrossOriginRequestsAllowed {
		Logger.Warn("The CORS policy is relaxed and cross-origin requests are allowed. This is for development, so don't use this option in production!")
	}
	Logger.Debug("Are mocks enabled? -> %v", config.AreMocksEnabled)
	Logger.Debug("Use dummy stacks? -> %v", config.UseDummyStacks)
}
