package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

type ValidationType int

const (
	User ValidationType = iota
	App
	Tag
	Password
	Origin
	Email
	Cookie
)

var validationTypeStrings = []string{"user", "app", "tag", "password", "origin", "email", "cookie"}

func getValidationTypeString(validationType ValidationType) string {
	return validationTypeStrings[validationType]
}

var (
	namePattern = regexp.MustCompile(`^[a-z0-9]{3,20}$`)
	tagPattern  = regexp.MustCompile(`^[a-z0-9.]{3,20}$`)
	// TODO Should be 8 at minimum
	passwordPattern = regexp.MustCompile(`^[a-zA-Z0-9!@#$%&_,.?]{3,30}$`)
	emailPattern    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	cookiePattern   = regexp.MustCompile(`^[a-f0-9]{64}$`)
)

func validate(input string, validationType ValidationType) error {
	var re *regexp.Regexp

	switch validationType {
	case User:
		re = namePattern
	case App:
		re = namePattern
	case Tag:
		re = tagPattern
	case Password:
		re = passwordPattern
	case Origin:
		return validateOrigin(input)
	case Email:
		re = emailPattern
	case Cookie:
		re = cookiePattern
	default:
		return fmt.Errorf("invalid validation type with index: %d", validationType)
	}

	if validationType == Email && len(input) > 64 {
		return fmt.Errorf("maximum email length of 64 characters is exceeded")
	}

	result := re.MatchString(input)
	if result == false {
		if validationType == Password || validationType == Cookie {
			Logger.Warn("input validation failed for validation type: %s", getValidationTypeString(validationType))
		} else {
			Logger.Warn("input validation failed for validation type '%s' with input '%s'", getValidationTypeString(validationType), input)
		}
		return fmt.Errorf("invalid signs or length of field: %s", getValidationTypeString(validationType))
	} else {
		return nil
	}
}

func validateOrigin(input string) error {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", input)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: %s", input)
	}

	host := parsedURL.Hostname()
	if len(host) == 0 {
		return fmt.Errorf("invalid domain: %s", input)
	}

	if parsedURL.Port() != "" {
		port, err := strconv.Atoi(parsedURL.Port())
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("invalid port: %s", input)
		}
	}

	if parsedURL.Path != "" || parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return fmt.Errorf("invalid URL - contains path, query or fragments: %s", input)
	}

	return nil
}
