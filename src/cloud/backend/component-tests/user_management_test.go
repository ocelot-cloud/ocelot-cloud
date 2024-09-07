package component_tests

import (
	"encoding/json"
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"ocelot/backend/security"
	"ocelot/backend/tools"
	"testing"
	"time"
)

type CloudClient struct {
	parent         utils.ComponentClient
	appToOperateOn string
	apps           []tools.AppInfo
}

// TODO Add a "wipe" endpoint that stops all stacks and it also deletes all users except "admin"
// TODO replace existing component-tests request logic with the CloudClient
// TODO test /api/check-auth, get user name and isAdmin == true
// TODO user registration, authorization and authentication etc
func TestLogin(t *testing.T) {
	cloud := getCloud()
	assert.Nil(t, cloud.parent.Cookie)
	assert.Nil(t, cloud.login())
	cookie := cloud.parent.Cookie
	assert.NotNil(t, cookie)
	assert.Equal(t, 64, len(cookie.Value))
	assert.True(t, cookie.Expires.After(time.Now().AddDate(0, 0, 29)))
	assert.True(t, cookie.Expires.Before(time.Now().AddDate(0, 0, 31)))
}

func TestAppManagement(t *testing.T) {
	cloud := getClientAndLogin(t)

	assert.Nil(t, cloud.startApp())
	apps, err := cloud.readApps()
	assert.Nil(t, err)
	assert.NotNil(t, apps)
	// TODO check app in more detail
	assert.Nil(t, cloud.stopApp())
	apps, err = cloud.readApps()
	assert.Nil(t, err)
	assert.NotNil(t, apps)
	// TODO check app in more detail
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

func getClientAndLogin(t *testing.T) *CloudClient {
	cloud := getCloud()
	assert.Nil(t, cloud.login())
	assert.Nil(t, cloud.wipeData())
	return &cloud
}

func getCloud() CloudClient {
	return CloudClient{
		utils.ComponentClient{
			"admin",
			"password",
			"password" + "x",
			"http://localhost:8080",
			nil,
			true,
			true,
			"http://localhost:8080",
		},
		"nginx-default",
		nil,
	}
}
