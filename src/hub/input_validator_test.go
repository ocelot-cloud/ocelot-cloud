package main

import (
	"fmt"
	"github.com/ocelot-cloud/shared/assert"
	"testing"
)

func TestValidateName(t *testing.T) {
	assert.Nil(t, validate("validusername", User))
	assert.Nil(t, validate("user123", User))
	assert.NotNil(t, validate("InvalidUsername", User))          // Contains uppercase
	assert.NotNil(t, validate("user!@#", User))                  // Contains special characters
	assert.NotNil(t, validate("us", User))                       // Too short
	assert.NotNil(t, validate("thisusernameiswaytoolong", User)) // Too long
}

func TestValidateTag(t *testing.T) {
	assert.Nil(t, validate("valid.tagname", Tag))
	assert.Nil(t, validate("tag123", Tag))
	assert.Nil(t, validate("tag.name123", Tag))
	assert.NotNil(t, validate("invalid.tagname!", Tag))             // Contains special characters other than dot
	assert.NotNil(t, validate("ta", Tag))                           // Too short
	assert.NotNil(t, validate("this.tagname.is.way.too.long", Tag)) // Too long
}

func TestValidatePassword(t *testing.T) {
	assert.Nil(t, validate("validpassword!", Password))
	assert.Nil(t, validate("valid_pass123", Password))
	assert.Nil(t, validate("InvalidPassword", Password)) // Contains uppercase
	assert.Nil(t, validate("valid!@#", Password))        // Contains special characters
	assert.NotNil(t, validate("1234567", Password))      // Too short
	assert.Nil(t, validate("12345678", Password))
	assert.NotNil(t, validate("thispasswordiswaytoolong_xxxxx!", Password)) // Too long
}

func TestValidateOrigin(t *testing.T) {
	assert.Nil(t, validate("http://example.com", Origin))
	assert.Nil(t, validate("https://example.com", Origin))
	assert.Nil(t, validate("http://example.com:8080", Origin))
	assert.NotNil(t, validate("ftp://example.com", Origin))
	assert.NotNil(t, validate("http://example.com/path", Origin))      // Contains path
	assert.NotNil(t, validate("ftp://example.com", Origin))            // Invalid scheme
	assert.NotNil(t, validate("https://example.com:99999", Origin))    // Invalid port
	assert.NotNil(t, validate("https://example.com?query=1", Origin))  // Contains query
	assert.NotNil(t, validate("https://example.com#fragment", Origin)) // Contains fragment
	assert.NotNil(t, validate("http://", Origin))                      // Missing hostname
	assert.NotNil(t, validate("http://:8080", Origin))                 // Missing hostname
}

func TestValidateCookie(t *testing.T) {
	sixtyOneHexDecimalLetters := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde"

	assert.NotNil(t, validate(sixtyOneHexDecimalLetters, Cookie))
	assert.Nil(t, validate(sixtyOneHexDecimalLetters+"f", Cookie))
	assert.NotNil(t, validate(sixtyOneHexDecimalLetters+"ff", Cookie))
	assert.NotNil(t, validate(sixtyOneHexDecimalLetters+"g", Cookie))
	assert.NotNil(t, validate("", Cookie))
}

func TestValidateEmail(t *testing.T) {
	assert.Nil(t, validate("admin@admin.com", Email))
	assert.NotNil(t, validate("@admin.com", Email))
	assert.NotNil(t, validate("admin@.com", Email))
	assert.NotNil(t, validate("admin@admin.", Email))
	assert.NotNil(t, validate("adminadmin.com", Email))
	assert.NotNil(t, validate("admin@admincom", Email))

	thirtyCharacters := "abcdefghijklmnopqrstuvwxyz1234"
	validEmail := fmt.Sprintf("%s@%s.de", thirtyCharacters, thirtyCharacters)
	assert.Nil(t, validate(validEmail, Email))
	tooLongEmail := fmt.Sprintf("%s@%s.com", thirtyCharacters, thirtyCharacters)
	assert.NotNil(t, validate(tooLongEmail, Email))
}
