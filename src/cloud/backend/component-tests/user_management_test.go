package component_tests

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"ocelot/backend/security"
	"testing"
	"time"
)

type CloudClient struct {
	parent utils.ComponentClient
	stack  string
}

// TODO Add a "wipe" endpoint that stops all stacks and it also deletes all users except "admin"
// TODO replace existing component-tests request logic with the CloudClient
// TODO test /api/check-auth, get user name and isAdmin == true
// TODO user registration, authorization and authentication etc
func TestLogin(t *testing.T) {
	cloud := getDefaultCloudClient()
	assert.Nil(t, cloud.parent.Cookie)
	assert.Nil(t, cloud.login())
	cookie := cloud.parent.Cookie
	assert.NotNil(t, cookie)
	assert.Equal(t, 64, len(cookie.Value))
	assert.True(t, cookie.Expires.After(time.Now().AddDate(0, 0, 29)))
	assert.True(t, cookie.Expires.Before(time.Now().AddDate(0, 0, 31)))
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

func getDefaultCloudClient() CloudClient {
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
	}
}
