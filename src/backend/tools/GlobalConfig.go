package tools

const NginxDefault = "nginx-default"
const NginxDefault2 = "nginx-default2"
const NginxCustomPath = "nginx-custom-path"
const NginxSlowStart = "nginx-slow-start"
const NginxDownloading = "nginx-download"

type GlobalConfig struct {
	AreCrossOriginRequestsAllowed bool
	AreMocksEnabled               bool
	IsGuiEnabled                  bool
	UseDummyStacks                bool
	UseRealDatabase               bool
	CreateDefaultAdminUser        bool
	OpenDataWipeEndpoint          bool
	UseRealHubClient              bool
	HttpScheme                    string
	RootDomain                    string
	PubliclyAvailablePort         string
	BackendExecutablePort         string
	Profile                       BackendProfile
}
