package tools

type AppInfo struct {
	Name    string `json:"name"`
	State   string `json:"state"`
	UrlPath string `json:"urlPath"`
}

// TODO I think I should get rid of this. Replace it by utils.SingleString
type StackInfo struct {
	Name string `json:"name"`
}
