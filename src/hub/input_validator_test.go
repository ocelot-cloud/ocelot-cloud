package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestValidateName(t *testing.T) {
	assert.True(t, isValidName("validusername"))
	assert.True(t, isValidName("user123"))
	assert.False(t, isValidName("InvalidUsername"))          // Contains uppercase
	assert.False(t, isValidName("user!@#"))                  // Contains special characters
	assert.False(t, isValidName("us"))                       // Too short
	assert.False(t, isValidName("thisusernameiswaytoolong")) // Too long
}

func TestValidateTag(t *testing.T) {
	assert.True(t, validateTag("valid.tagname"))
	assert.True(t, validateTag("tag123"))
	assert.True(t, validateTag("tag.name123"))
	assert.False(t, validateTag("invalid.tagname!"))             // Contains special characters other than dot
	assert.False(t, validateTag("ta"))                           // Too short
	assert.False(t, validateTag("this.tagname.is.way.too.long")) // Too long
}

func TestValidatePasswords(t *testing.T) {
	assert.True(t, validatePasswords("validpassword!"))
	assert.True(t, validatePasswords("valid_pass123"))
	assert.False(t, validatePasswords("InvalidPassword"))           // Contains uppercase
	assert.True(t, validatePasswords("valid!@#"))                   // Contains special characters
	assert.False(t, validatePasswords("vp"))                        // Too short
	assert.False(t, validatePasswords("thispasswordiswaytoolong!")) // Too long
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
