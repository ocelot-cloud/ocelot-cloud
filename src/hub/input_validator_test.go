package main

import (
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestValidateName(t *testing.T) {
	assert.True(t, validate("validusername", User))
	assert.True(t, validate("user123", User))
	assert.False(t, validate("InvalidUsername", User))          // Contains uppercase
	assert.False(t, validate("user!@#", User))                  // Contains special characters
	assert.False(t, validate("us", User))                       // Too short
	assert.False(t, validate("thisusernameiswaytoolong", User)) // Too long
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
	assert.False(t, validate("InvalidPassword", Password))                 // Contains uppercase
	assert.True(t, validate("valid!@#", Password))                         // Contains special characters
	assert.False(t, validate("vp", Password))                              // Too short
	assert.False(t, validate("thispasswordiswaytoolong_xxxxx!", Password)) // Too long
}

func TestValidateOrigin(t *testing.T) {
	assert.True(t, validate("http://example.com", Origin))
	assert.True(t, validate("https://example.com", Origin))
	assert.True(t, validate("http://example.com:8080", Origin))
	assert.False(t, validate("ftp://example.com", Origin))
	assert.False(t, validate("http://example.com/path", Origin))      // Contains path
	assert.False(t, validate("ftp://example.com", Origin))            // Invalid scheme
	assert.False(t, validate("https://example.com:99999", Origin))    // Invalid port
	assert.False(t, validate("https://example.com?query=1", Origin))  // Contains query
	assert.False(t, validate("https://example.com#fragment", Origin)) // Contains fragment
	assert.False(t, validate("http://", Origin))                      // Missing hostname
	assert.False(t, validate("http://:8080", Origin))                 // Missing hostname
}

func TestValidateCookie(t *testing.T) {
	sixtyOneHexDecimalLetters := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde"

	assert.False(t, validate(sixtyOneHexDecimalLetters, Cookie))
	assert.True(t, validate(sixtyOneHexDecimalLetters+"f", Cookie))
	assert.False(t, validate(sixtyOneHexDecimalLetters+"ff", Cookie))
	assert.False(t, validate(sixtyOneHexDecimalLetters+"g", Cookie))
	assert.False(t, validate("", Cookie))
}
