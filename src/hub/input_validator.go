package main

import (
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
	namePattern     = regexp.MustCompile(`^[a-z0-9]{3,20}$`)
	tagPattern      = regexp.MustCompile(`^[a-z0-9.]{3,20}$`)
	passwordPattern = regexp.MustCompile(`^[a-zA-Z0-9!@#$%&_,.?]{3,30}$`)
	emailPattern    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	cookiePattern   = regexp.MustCompile(`^[a-f0-9]{64}$`)
)

func validate(input string, validationType ValidationType) bool {
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
		return false
	}

	result := re.MatchString(input)
	if result == false {
		if validationType == Password || validationType == Cookie {
			Logger.Warn("input validation failed for validation type: %s", getValidationTypeString(validationType))
		} else {
			Logger.Warn("input validation failed for validation type '%s' with input '%s'", getValidationTypeString(validationType), input)
		}
	}
	return result
}

func validateOrigin(input string) bool {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	host := parsedURL.Hostname()
	if len(host) == 0 {
		return false
	}

	if parsedURL.Port() != "" {
		port, err := strconv.Atoi(parsedURL.Port())
		if err != nil || port < 1 || port > 65535 {
			return false
		}
	}

	if parsedURL.Path != "" || parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return false
	}

	return true
}
