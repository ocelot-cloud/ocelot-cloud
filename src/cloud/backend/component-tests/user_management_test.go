package component_tests

import (
	"github.com/ocelot-cloud/shared/utils"
	"testing"
)

type CloudClient struct {
	parent utils.ComponentClient
	stack  string
}

// TODO user registration, login?, authorization and authentication etc
func TestAsd(t *testing.T) {
	client := getDefaultCloudClient()
	println(client.stack)
}

func getDefaultCloudClient() CloudClient {
	return CloudClient{
		utils.ComponentClient{
			"admin",
			"password",
			"password" + "x",
			"http://localhost:8081",
			nil,
			true,
			true,
			"http://localhost:8081",
		},
		"nginx-default",
	}
}
