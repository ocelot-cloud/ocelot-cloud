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

type UserAndApp struct {
	User string `json:"user"`
	App  string `json:"app"`
}

type TagInfo struct {
	User string `json:"user"`
	App  string `json:"app"`
	Tag  string `json:"tag"`
}
