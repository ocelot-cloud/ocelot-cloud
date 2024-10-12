package component_tests

import (
	"encoding/json"
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"ocelot/backend/setup"
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

// TODO abstract paths

func (c *CloudClient) getSecret() (string, error) {
	body, err := c.parent.DoRequest("/api/secret", nil, "")
	if err != nil {
		return "", err
	}

	var secret string
	err = json.Unmarshal(body.([]byte), &secret)
	if err != nil {
		return "", err
	}

	return secret, nil
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
	loginCredentials := setup.Credentials{c.parent.User, c.parent.Password}
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

// TODO Subsequently are the new app functions. When implemented, delete the old ones.
func (c *CloudClient) searchHubApps() (*[]tools.UserAndApp, error) {
	responseBody, err := c.parent.DoRequest("/api/apps/search", utils.SingleString{"sample"}, "")
	if err != nil {
		return nil, err
	}
	var userAndAppList []tools.UserAndApp
	err = json.Unmarshal(responseBody.([]byte), &userAndAppList)
	if err != nil {
		return nil, err
	}

	return &userAndAppList, nil
}

func (c *CloudClient) getHubTags(userAndApp tools.UserAndApp) (*[]string, error) {
	responseBody, err := c.parent.DoRequest("/api/tags/list", userAndApp, "")
	if err != nil {
		return nil, err
	}
	var tags []string
	err = json.Unmarshal(responseBody.([]byte), &tags)
	if err != nil {
		return nil, err
	}

	return &tags, nil
}

func (c *CloudClient) downloadTagFromHub(tagInfo tools.TagInfo) error {
	_, err := c.parent.DoRequest("/api/tags/download", tagInfo, "")
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudClient) startAppNew(tagInfo tools.TagInfo) error {
	_, err := c.parent.DoRequest("/api/apps/start", tagInfo, "")
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudClient) stopAppNew(tagInfo tools.TagInfo) error {
	_, err := c.parent.DoRequest("/api/apps/stop", tagInfo, "")
	if err != nil {
		return err
	}
	return nil
}
