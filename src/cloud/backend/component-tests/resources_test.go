package component_tests

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"ocelot/backend/security"
	"ocelot/backend/tools"
	"os"
	"testing"
	"time"
)

var logger = tools.Logger

const (
	backendBaseUrl = "http://localhost:8080"
	TestProfile    = "TEST"
	ProdProfile    = "PROD"
)

type CloudClient struct {
	parent         utils.ComponentClient
	appToOperateOn string
	apps           []tools.AppInfo
}

func (c *CloudClient) startApp() error {
	data := utils.SingleString{c.appToOperateOn}
	_, err := c.parent.DoRequest("/api/stacks/deploy", data, "")
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudClient) readApps() (*[]tools.AppInfo, error) {
	responseBody, err := c.parent.DoRequest("/api/stacks/read", nil, "")
	if err != nil {
		return nil, err
	}

	var apps []tools.AppInfo
	err = json.Unmarshal(responseBody.([]byte), &apps)
	if err != nil {
		return nil, err
	}

	return &apps, nil
}

func (c *CloudClient) stopApp() error {
	data := utils.SingleString{c.appToOperateOn}
	_, err := c.parent.DoRequest("/api/stacks/stop", data, "")
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudClient) login() error {
	loginCredentials := security.Credentials{c.parent.User, c.parent.Password}
	// TODO should be an expected message like in the hub
	_, err := c.parent.DoRequest("/api/login", loginCredentials, "")
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudClient) wipeData() error {
	_, err := c.parent.DoRequest("/api/stacks/wipe-data", nil, "")
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudClient) assertState(expectedState string) error {
	const maxAttempts = 300
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		appInfo, err := c.readApps()
		if err != nil {
			return err
		}

		actualState := getState(c.appToOperateOn, appInfo)
		if actualState == expectedState {
			return nil
		}

		if attempt%10 == 0 {
			logger.Info("%v: App '%s' has state '%s' instead of expected state '%s'. Re-try in one second...", attempt/10, c.appToOperateOn, actualState, expectedState)
		}

		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("state not reached within time range")
}

func getState(stackName string, appInfo *[]tools.AppInfo) string {
	for _, singleInfo := range *appInfo {
		if singleInfo.Name == stackName {
			return singleInfo.State
		}
	}
	return ""
}

func getClientAndLogin(t *testing.T) *CloudClient {
	cloud := getCloud()
	assert.Nil(t, cloud.login())
	assert.Nil(t, cloud.wipeData())
	return cloud
}

func getCloud() *CloudClient {
	cloud := &CloudClient{
		utils.ComponentClient{
			"admin",
			"password",
			"password" + "x",
			"http://ocelot-cloud.localhost",
			nil,
			true,
			true,
			"http://ocelot-cloud.localhost",
		},
		"nginx-default",
		nil,
	}

	if os.Getenv("PROFILE") == "TEST" {
		cloud.parent.Origin = "http://localhost:8080"
		cloud.parent.RootUrl = "http://localhost:8080"
	}

	return cloud
}
