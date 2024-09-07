package component_tests

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/utils"
	"ocelot/backend/security"
	"testing"
)

type CloudClient struct {
	parent utils.ComponentClient
	stack  string
}

// TODO user registration, login?, authorization and authentication etc
func TestLogin(t *testing.T) {
	client := getDefaultCloudClient()
	println(client.stack)
	assert.Nil(t, client.parent.Cookie)
	assert.Nil(t, client.login())
	assert.NotNil(t, client.parent.Cookie)
	// TODO Further assertions
	// TODO also make check-auth, get user name and isAdmin == true
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
