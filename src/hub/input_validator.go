package main

import (
	"net/url"
	"regexp"
	"strconv"
)

func validateName(input string) bool {
	pattern := `^[a-z0-9]{3,20}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}

func validateTag(input string) bool {
	pattern := `^[a-z0-9.]{3,20}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}

func validatePasswords(input string) bool {
	pattern := `^[a-z0-9!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':",.<>\/?\\|` + "`" + `~]{3,20}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
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
