package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestValidateName(t *testing.T) {
	assert.True(t, validate("validusername", Name))
	assert.True(t, validate("user123", Name))
	assert.False(t, validate("InvalidUsername", Name))          // Contains uppercase
	assert.False(t, validate("user!@#", Name))                  // Contains special characters
	assert.False(t, validate("us", Name))                       // Too short
	assert.False(t, validate("thisusernameiswaytoolong", Name)) // Too long
}

func TestValidateTag(t *testing.T) {
	assert.True(t, validate("valid.tagname", Tag))
	assert.True(t, validate("tag123", Tag))
	assert.True(t, validate("tag.name123", Tag))
	assert.False(t, validate("invalid.tagname!", Tag))             // Contains special characters other than dot
	assert.False(t, validate("ta", Tag))                           // Too short
	assert.False(t, validate("this.tagname.is.way.too.long", Tag)) // Too long
}

func TestValidatePassword(t *testing.T) {
	assert.True(t, validate("validpassword!", Password))
	assert.True(t, validate("valid_pass123", Password))
	assert.False(t, validate("InvalidPassword", Password))           // Contains uppercase
	assert.True(t, validate("valid!@#", Password))                   // Contains special characters
	assert.False(t, validate("vp", Password))                        // Too short
	assert.False(t, validate("thispasswordiswaytoolong!", Password)) // Too long
}

func TestValidateOrigin(t *testing.T) {
	assert.True(t, validateOrigin("http://example.com"))
	assert.True(t, validateOrigin("https://example.com"))
	assert.True(t, validateOrigin("http://example.com:8080"))
	assert.False(t, validateOrigin("http://example.com/path"))      // Contains path
	assert.False(t, validateOrigin("ftp://example.com"))            // Invalid scheme
	assert.False(t, validateOrigin("https://example.com:99999"))    // Invalid port
	assert.False(t, validateOrigin("https://example.com?query=1"))  // Contains query
	assert.False(t, validateOrigin("https://example.com#fragment")) // Contains fragment
	assert.False(t, validateOrigin("http://"))                      // Missing hostname
	assert.False(t, validateOrigin("http://:8080"))                 // Missing hostname
}
