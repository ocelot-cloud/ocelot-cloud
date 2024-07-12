package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestCreateFileInfo(t *testing.T) {
	result, err := createAppAndTag("app_tag.tar.gz")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "app", result.App)
	assert.Equal(t, "tag", result.Tag)

	result, err = createAppAndTag("apptag.tar.gz")
	assert.NotNil(t, err)
	assert.Nil(t, result)

	result, err = createAppAndTag("app_tag_extra.tar.gz")
	assert.NotNil(t, err)
	assert.Nil(t, result)

	result, err = createAppAndTag("app_long_tag.tar.gz")
	assert.NotNil(t, err)
	assert.Nil(t, result)

	result, err = createAppAndTag("app_tag.tar2.gz")
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
