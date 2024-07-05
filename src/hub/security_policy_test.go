//go:build acceptance

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func applySecurityPolicy(operation Operation, req *http.Request) *http.Request {
	policyToApply := policy.getPolicyFor(operation)
	if policyToApply.IsCredentialsRequired {
		creds := LoginCredentials{
			Username: Hub.Username,
			Password: Hub.Password,
		}
		payloadBytes, _ := json.Marshal(creds)
		payloadReader := bytes.NewReader(payloadBytes)
		req.Body = io.NopCloser(payloadReader)
	}
	if policyToApply.IsOriginRequired {
		req.Header.Set("Origin", Hub.Origin)
	}
	if policyToApply.IsCookieRequired {
		req.Header.Set(Hub.Cookie.Name, Hub.Cookie.Value)
	}
	return req
}

var policy = getPolicy()

type Operation int

// Define constants for DaysOfWeek
const (
	FindApps Operation = iota
	DownloadApp
	Register
	ChangeOrigin
	ChangePassword
	Login
	DeleteUser
	CreateApp
	DeleteApp
	UploadTag
	DeleteTag
)

func getPolicy() SecurityPolicyCollection {
	return SecurityPolicyCollection{
		policies: []*SecurityPolicy{
			{false, false, false, []Operation{FindApps, DownloadApp}},
			{true, false, false, []Operation{Register, ChangeOrigin, ChangePassword}},
			{true, true, false, []Operation{Login}},
			{false, true, true, []Operation{DeleteUser, CreateApp, DeleteApp, UploadTag, DeleteTag}},
		},
	}
}

type SecurityPolicyCollection struct {
	policies []*SecurityPolicy
}

func (s *SecurityPolicyCollection) getPolicyFor(operation Operation) *SecurityPolicy {
	for _, policy := range s.policies {
		for _, currentOperation := range policy.OperationsAffected {
			if currentOperation == operation {
				return policy
			}
		}
	}
	panicMessage := fmt.Sprintf("policy for operation '%v' does not exist", operation)
	panic(panicMessage)
}

type SecurityPolicy struct {
	IsCredentialsRequired bool
	IsOriginRequired      bool
	IsCookieRequired      bool
	OperationsAffected    []Operation
}
