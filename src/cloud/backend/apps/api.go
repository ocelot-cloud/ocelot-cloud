package apps

import (
	"github.com/ocelot-cloud/shared/utils"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ocelot/backend/security"
	"strings"
)

// TODO Make sure to remove the ocelot cookie before proxying a request to the service behind, so that it can't read/steal it.
func ProxyRequestToTheDockerContainer(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(security.CookieName) == "" {
		http.Error(w, "cookie not found", http.StatusUnauthorized)
		return
	}

	logger.Trace("Proxying request with target host %s", r.Host)
	host, _, _ := net.SplitHostPort(r.Host)
	if host == "" {
		host = r.Host
	}
	targetContainer := strings.TrimSuffix(host, "."+config.RootDomain)
	targetPort := stackConfigService.GetAppConfig(targetContainer).Port
	targetURL, err := url.Parse("http://" + targetContainer + ":" + targetPort)
	logger.Debug("proxying to target URL: %s", targetURL)
	if err != nil {
		logger.Error("error when parsing URL, %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO add tests
	urlSecret := r.URL.Query().Get("secret")
	/* TODO
	var secret string
	headerSecret := r.Header.Get("ocelot-auth")
	if urlSecret != "" {
		secret = urlSecret
	} else if {
		secret = headerSecret
	} else {
		respond: not authenticated
		return
	}
	*/

	// TODO somewhere here I need to  1) add the cookie authentication and 2) add the group based authorization logic.

	if urlSecret != "" {
		cookie, err := utils.GenerateCookie()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		cookie.Name = security.CookieName // TODO Duplication with other cookie generation call.
		cookie.Value = urlSecret
		http.SetCookie(w, cookie)

		redirectURL := *r.URL
		redirectURL.RawQuery = ""
		http.Redirect(w, r, redirectURL.String(), http.StatusFound) // TODO write a test for that redirect.
		return
	}

	// TODO Should be "ocelot-auth" to avoid conflicts. Also abstract.
	// TODO Tell the browser to re-do the request but without "secret" query param? Maybe via redirecting to the same URL? I dont want the secret to be exposed so long in the URL.

	// the path of original request is preserved
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.URL.Host = targetContainer
	r.URL.Scheme = "http"
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.URL.Query().Del("secret") // TODO to be tested
	proxy.ServeHTTP(w, r)

	// TODO Cookie updates should be done
}
