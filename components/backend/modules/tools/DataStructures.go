package tools

type ResponsePayloadDto struct {
	Name    string `json:"name"`
	State   string `json:"state"`
	UrlPath string `json:"urlPath"`
}

type StackInfo struct {
	Name string `json:"name"`
}
