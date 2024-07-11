package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestCreateFileInfo(t *testing.T) {
	result, err := createFileInfo("user_app_tag.tar.gz")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user", result.User)
	assert.Equal(t, "app", result.App)
	assert.Equal(t, "tag", result.Tag)

	result, err = createFileInfo("user_apptag.tar.gz")
	assert.NotNil(t, err)
	assert.Nil(t, result)

	result, err = createFileInfo("user_app_tag_extra.tar.gz")
	assert.NotNil(t, err)
	assert.Nil(t, result)

	result, err = createFileInfo("user_app_long_tag.tar.gz")
	assert.NotNil(t, err)
	assert.Nil(t, result)

	result, err = createFileInfo("user_app_tag.tar2.gz")
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
