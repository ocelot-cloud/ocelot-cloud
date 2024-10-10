package apps

import (
	"github.com/ocelot-cloud/shared/utils"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ocelot/backend/repo"
	"ocelot/backend/tools"
	"strings"
)

// TODO Logic looks quite complex. Maybe extract to unit and test it?
// TODO Shouldn't this be part of security module?
// TODO Dont use cookie anymore as secret. Create a separate one.
// TODO Make sure to remove the ocelot cookie before proxying a request to the service behind, so that it can't read/steal it.
func ProxyRequestToTheDockerContainer(w http.ResponseWriter, r *http.Request) {
	logger.Trace("Proxying request with target URL %s%s", r.URL, r.URL.Path) // TODO temp? Might contain secret/cookie info
	host, _, _ := net.SplitHostPort(r.Host)
	if host == "" {
		host = r.Host
	}
	targetContainer := strings.TrimSuffix(host, "."+tools.Config.RootDomain)
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

		cookieValue, err := repo.UserRepo.GetAssociatedCookieValueAndDeleteSecret(urlSecret)
		if err != nil {
			logger.Error("Failed to get associated cookie value or delete secret: %v", err)
			http.Error(w, "failed to get associated cookie value or delete secret", http.StatusInternalServerError)
			return
		}

		cookie.Name = tools.CookieName // TODO Can be removed? Duplication?
		cookie.Value = cookieValue
		http.SetCookie(w, cookie)

		// TODO I dont have a user here. -> so the logic should be: does the secret exist? If so, delete it and use the cookie from the database here.

		// TODO Also add a mechanism that makes secrets older than 5 minutes invalid.
		// TODO Is there a ways to search for unused methods and variables? Maybe let CI fail if such are detected.

		redirectURL := *r.URL
		redirectURL.RawQuery = ""
		http.Redirect(w, r, redirectURL.String(), http.StatusFound) // TODO write a test for that redirect.
		return
	}

	// TODO Update cookie expiration time if cookie is valid or a valid secret was provided.

	// TODO Use the auth info to allow access to the app.
	_, err = repo.GetAuthentication(w, r)
	if err != nil {
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
