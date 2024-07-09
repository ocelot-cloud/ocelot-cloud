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
	TagFile
	Email
)

var (
	namePattern     = regexp.MustCompile(`^[a-z0-9]{3,20}$`)
	tagPattern      = regexp.MustCompile(`^[a-z0-9.]{3,20}$`)
	passwordPattern = regexp.MustCompile(`^[a-z0-9!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':",.<>\/?\\|` + "`" + `~]{3,20}$`)
	emailPattern    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
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
	case TagFile:
		return validateTagFile(input)
	case Email:
		re = emailPattern
	default:
		return false
	}

	return re.MatchString(input)
}

func validateTagFile(input string) bool {
	fileInfo, err := createFileInfo(input)
	if err != nil {
		return false
	}
	return validate(fileInfo.User, User) && validate(fileInfo.App, App) && validate(fileInfo.Tag, Tag)
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
