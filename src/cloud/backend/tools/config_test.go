package tools

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

// TODO I think ci-runner is not yet executing these test.

func TestTheTestProfileHostParams(t *testing.T) {
	hostParams := getHostParams(TEST, "")
	assert.Equal(t, "http", hostParams.Scheme)
	assert.Equal(t, "localhost", hostParams.Domain)
	assert.Equal(t, "8080", hostParams.Port)
}

func TestTheProdProfileHostParams(t *testing.T) {
	hostParams := getHostParams(PROD, "https://my-domain.com")
	assert.Equal(t, "https", hostParams.Scheme)
	assert.Equal(t, "my-domain.com", hostParams.Domain)
	assert.Equal(t, "443", hostParams.Port)

	hostParams = getHostParams(PROD, "http://my-domain.com")
	assert.Equal(t, "http", hostParams.Scheme)
	assert.Equal(t, "my-domain.com", hostParams.Domain)
	assert.Equal(t, "80", hostParams.Port)

	hostParams = getHostParams(PROD, "http://my-domain.com:3000")
	assert.Equal(t, "http", hostParams.Scheme)
	assert.Equal(t, "my-domain.com", hostParams.Domain)
	assert.Equal(t, "3000", hostParams.Port)
}
