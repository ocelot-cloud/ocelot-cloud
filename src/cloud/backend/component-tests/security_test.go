//go:build security

package component_tests

import (
	"github.com/ocelot-cloud/shared/assert"
	"net/http"
	"ocelot/backend/tools"
	"testing"
	"time"
)

// These tests can only be run against a Docker container because only Docker containers can proxy to application containers.
func TestAppAccess(t *testing.T) {
	cloud2 := getCloud() // TODO This should directly return a reference
	cloud := &cloud2
	cloud.parent.RootUrl = "http://ocelot-cloud.localhost"
	assert.Nil(t, cloud.login())
	assert.Nil(t, cloud.startApp())

	waitForContainer(t, cloud)

	resp, err := http.Get("http://nginx-default.localhost")
	if err != nil {
		logger.Error("app request failed %v: ", err)
		t.Fail()
	}
	assert.Equal(t, 401, resp.StatusCode)
	defer resp.Body.Close()

	/* TODO
	cloud.parent.Cookie = nil
	_, err = cloud.parent.DoRequest("", nil, "TODO")
	assert.NotNil(t, err)
	assert.Equal(t, utils.GetErrMsg(302, "asdf"), err.Error())
	*/

	// TODO set cookie to nil, try to access app -> denied
	// TODO secret as query param -> redirect
	// TODO wrong cookie -> denied
}

func waitForContainer(t *testing.T, cloud *CloudClient) {
	start := time.Now()

	for {
		apps, err := cloud.readApps()
		if err != nil {
			return
		}

		var testApp tools.AppInfo
		for _, app := range *apps {
			if app.Name == "nginx-default" {
				testApp = app
				break
			}
		}

		if testApp.State == "Available" {
			break
		}

		if time.Since(start) > 10*time.Second {
			t.Fail()
			return
		}

		time.Sleep(100 * time.Millisecond)
	}
}
