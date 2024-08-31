package apps

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// TODO Make sure to remove the ocelot cookie before proxying a request to the service behind, so that it can't read/steal it.
func ProxyRequestToTheDockerContainer(w http.ResponseWriter, r *http.Request) {
	logger.Trace("Proxying request with target host %s", r.Host)
	targetContainer := strings.TrimSuffix(r.Host, "."+config.RootDomain)
	targetPort := stackConfigService.getAppConfig(targetContainer).Port
	targetURL, err := url.Parse("http://" + targetContainer + ":" + targetPort)
	if err != nil {
		logger.Error("error when parsing URL, %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// the path of original request is preserved
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.URL.Host = targetContainer
	r.URL.Scheme = "http"
	r.Header.Set("X-Forwarded-Host", r.Host)
	proxy.ServeHTTP(w, r)
}
