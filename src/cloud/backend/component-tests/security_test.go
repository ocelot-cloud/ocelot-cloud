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
	assert.Nil(t, cloud.login())
	cookieValue := cloud.parent.Cookie.Value

	assert.Nil(t, cloud.startApp())
	waitForContainer(t, cloud)

	/*
		// TODO replace with DoRequest function or so? "RequestApp()"
		resp, err := http.Get("http://nginx-default.localhost")
		if err != nil {
			logger.Error("app request failed %v: ", err)
			t.Fail()
		}
		assert.Equal(t, 401, resp.StatusCode)
		defer resp.Body.Close()
	*/

	checkIfRedirectViaSecretWorks(t, cookieValue)

	// TODO wrong cookie -> denied
}

func checkIfRedirectViaSecretWorks(t *testing.T, cookieValue string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", "http://nginx-default.localhost?secret="+cookieValue, nil)
	if err != nil {
		logger.Error("app request failed %v: ", err)
		t.Fail()
	}
	req.AddCookie(&http.Cookie{Name: "ocelot-auth", Value: cookieValue})
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("app request failed %v: ", err)
		t.Fail()
	}
	assert.Equal(t, 302, resp.StatusCode)
	assert.Equal(t, cookieValue, resp.Cookies()[0].Value)
	assert.Equal(t, "/", resp.Header.Get("Location"))
	defer resp.Body.Close()
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
