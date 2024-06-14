package tools

import (
	"github.com/ocelot-cloud/shared"
	"testing"
)

func TestEvaluateLogLevel(t *testing.T) {
	testCases := []struct {
		name             string
		profile          BackendComponentMode
		userLogLevel     string
		expectedLogLevel shared.LogLevelValue // Assuming LogLevelValue is a type you've defined
	}{
		{"Default LogLevelValue of ProdMode", ProdWithGui, "notSet", shared.INFO},
		{"Default LogLevelValue of DevMode", DevelopmentSetup, "notSet", shared.DEBUG},
		{"Default LogLevelValue of DevMockedMode", DependenciesMocked, "notSet", shared.INFO},
		{"LogLevelValue of DevelopmentMode", DevelopmentSetup, "notSet", shared.DEBUG},
		{"LogLevelValue of ProdMode with Trace", ProdWithGui, "trace", shared.TRACE},
		{"LogLevelValue of ProdMode with Debug", ProdWithGui, "debug", shared.DEBUG},
		{"LogLevelValue of DevMode with Info", DevelopmentSetup, "info", shared.INFO},
		{"LogLevelValue of ProdMode with Warn", ProdWithGui, "warn", shared.WARN},
		{"LogLevelValue of ProdMode with Error", ProdWithGui, "error", shared.ERROR},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualLogLevel := EvaluateLogLevelBasedOn(tc.profile, tc.userLogLevel)
			shared.AssertEqual(t, tc.expectedLogLevel, actualLogLevel)
		})
	}
}

func TestPanicForInvalidLogLevel(t *testing.T) {
	shared.AssertPanics(t, func() {
		EvaluateLogLevelBasedOn(ProdWithGui, "invalid value")
	})
}

func TestGlobalConfig(t *testing.T) {
	testCases := []struct {
		name           string
		profile        BackendComponentMode
		useMock        bool
		isGuiEnabled   bool
		isCorsDisabled bool
	}{
		{"Prod Profile", ProdWithGui, false, true, false},
		{"Dev Mocked Profile", DependenciesMocked, true, false, false},
		{"Dev Setup Profile", DevelopmentSetup, false, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := SetGlobalConfig(tc.profile, "notSet", false, false)
			shared.AssertEqual(t, config.AreMocksEnabled, tc.useMock)
			shared.AssertEqual(t, config.IsGuiEnabled, tc.isGuiEnabled)
			shared.AssertEqual(t, config.AreCrossOriginRequestsAllowed, tc.isCorsDisabled)
		})
	}
}
