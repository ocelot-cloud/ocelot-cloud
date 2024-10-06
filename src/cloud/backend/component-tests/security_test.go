//go:build security

package component_tests

import (
	"github.com/ocelot-cloud/shared/assert"
	"net/http"
	"ocelot/backend/tools"
	"testing"
	"time"
)

func TestAppAccess(t *testing.T) {
	cloud := getCloud()
	assert.Nil(t, cloud.login())
	cookieValue := cloud.parent.Cookie.Value

	assert.Nil(t, cloud.startApp())
	waitForContainer(t, cloud)

	cloud.parent.Cookie = nil
	assertUnauthorizedAppAccess(t)
	checkIfRedirectViaSecretWorks(t, cookieValue)

	// TODO wrong cookie -> denied
}

func assertUnauthorizedAppAccess(t *testing.T) {
	resp, err := http.Get("http://nginx-default.localhost")
	if err != nil {
		logger.Error("app request failed %v: ", err)
		t.Fail()
	}
	assert.Equal(t, 401, resp.StatusCode)
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
	responseCookie := resp.Cookies()[0]
	assert.Equal(t, tools.CookieName, responseCookie.Name)
	assert.Equal(t, cookieValue, responseCookie.Value)
	assert.Equal(t, "/", resp.Header.Get("Location"))
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

// TODO make thisRootURL = ENV("ROOT_URL"), default to "http://ocelot-cloud.localhost:8080"?
func TestSecretGeneration(t *testing.T) {
	cloud := getClientAndLogin(t)
	secret, err := cloud.getSecret()
	assert.Nil(t, err)
	assert.Equal(t, 64, len(secret))

	// TODO extra test?
	secret2, _ := cloud.getSecret()
	assert.NotEqual(t, secret, secret2)
}

// TODO Test that a secret only works once. After exchanging it against a cookie, do it a second time and it should fail.
// TODO Make an assertion that cookie value is different from secret
