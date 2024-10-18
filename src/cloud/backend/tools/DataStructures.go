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

type AppInfoNew struct {
	App         RepoApp
	Port        string
	Path        string
	IsAvailable bool
}

// TODO Put ID's first in the structs.
type RepoApp struct {
	Maintainer      string
	Name            string
	AppId           int
	ActiveTagName   string
	ActiveTagId     int
	ShouldBeRunning bool // TODO Implement in database, set when starting/stopping app
}

// TODO Mabye put in shared module and reuse in hub when refactoring its API?
type SingleInt struct {
	Value int `json:"value"`
}

// TODO Duplication with hub, put to shared module
type App struct {
	Maintainer string `json:"user"`
	Name       string `json:"name"`
	Id         int    `json:"id"`
}

// TODO Duplication with hub, put to shared module
type Tag struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
