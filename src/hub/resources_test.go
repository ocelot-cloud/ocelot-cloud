package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"github.com/ocelot-cloud/shared/hub"
	"testing"
)

type Operation int

const (
	FindApps Operation = iota
	DownloadTag
	Register
	ChangePassword
	Login
	CreateApp
	DeleteApp
	UploadTag
	DeleteTag
	GetTags
	CheckAuth
)

func GetHubAndLogin(t *testing.T) *hub.HubClient {
	hub := hub.GetHub()
	assert.Nil(t, hub.RegisterUser())
	err := hub.Login()
	assert.Nil(t, err)
	return hub
}
