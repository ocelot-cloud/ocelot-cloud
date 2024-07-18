package main

type TagInfo struct {
	User string `json:"user"`
	App  string `json:"app"`
	Tag  string `json:"tag"`
}

type AppAndTag struct {
	App string `json:"app"`
	Tag string `json:"tag"`
}

type TagUpload struct {
	App     string `json:"app"`
	Tag     string `json:"tag"`
	Content []byte `json:"content"`
}

type UserAndApp struct {
	User string `json:"user"`
	App  string `json:"app"`
}

type SingleString struct {
	Value string `json:"name"`
}

type ChangePasswordForm struct {
	User        string `json:"user"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangeOriginForm struct {
	User      string `json:"user"`
	Password  string `json:"password"`
	NewOrigin string `json:"new_origin"`
}

type RegistrationForm struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Origin   string `json:"host"`
}

type LoginCredentials struct {
	User     string `json:"username"`
	Password string `json:"password"`
}
