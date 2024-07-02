package tools

const NginxDefault = "nginx-default"
const NginxDefault2 = "nginx-default2"
const NginxCustomPath = "nginx-custom-path"
const NginxSlowStart = "nginx-slow-start"
const NginxDownloading = "nginx-download"

type GlobalConfig struct {
	// AreCrossOriginRequestsAllowed controls whether the server will accept cross-origin requests.
	// Setting this to true relaxes the CORS policy by allowing specific cross-origin requests.
	AreCrossOriginRequestsAllowed    bool
	AreMocksEnabled                  bool
	IsGuiEnabled                     bool
	IsSecurityEnabled                bool
	BackendMode                      BackendComponentMode
	WaitForSecurityBeforeOpeningPort bool
	UseDummyStacks                   bool
	Scheme                           string // "http" or "https"
	RootDomain                       string // e.g. "localhost"
	Port                             string // e.g. "8082"
}
