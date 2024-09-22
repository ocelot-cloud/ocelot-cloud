package tools

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestTheTestProfileHostParams(t *testing.T) {
	hostParams, err := getHostParams(TEST, "")
	assert.Nil(t, err)
	assert.Equal(t, "http", hostParams.Scheme)
	assert.Equal(t, "localhost", hostParams.Domain)
	assert.Equal(t, "8080", hostParams.Port)
}

func TestTheProdProfileHostParams(t *testing.T) {
	hostParams, err := getHostParams(PROD, "https://my-domain.com")
	assert.Nil(t, err)
	assert.Equal(t, "https", hostParams.Scheme)
	assert.Equal(t, "my-domain.com", hostParams.Domain)
	assert.Equal(t, "443", hostParams.Port)

	hostParams, err = getHostParams(PROD, "http://my-domain.com")
	assert.Nil(t, err)
	assert.Equal(t, "http", hostParams.Scheme)
	assert.Equal(t, "my-domain.com", hostParams.Domain)
	assert.Equal(t, "80", hostParams.Port)

	hostParams, err = getHostParams(PROD, "http://my-domain.com:3000")
	assert.Nil(t, err)
	assert.Equal(t, "http", hostParams.Scheme)
	assert.Equal(t, "my-domain.com", hostParams.Domain)
	assert.Equal(t, "3000", hostParams.Port)
}

func TestIfGetHostParamsFailsCorrectly(t *testing.T) {
	_, err := getHostParams(PROD, "")
	assert.NotNil(t, err)
	_, err = getHostParams(PROD, "htt://my-domain.com")
	assert.NotNil(t, err)
}
