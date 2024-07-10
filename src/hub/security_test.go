package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (h *HubClient) applySecurityPolicy(operation Operation, req *http.Request) *http.Request {
	policyToApply := securityPolicies.getPolicyFor(operation)
	if policyToApply.IsCredentialsRequired {
		creds := LoginCredentials{
			User:     h.User,
			Password: h.Password,
		}
		payloadBytes, _ := json.Marshal(creds)
		payloadReader := bytes.NewReader(payloadBytes)
		req.Body = io.NopCloser(payloadReader)
	}
	if policyToApply.IsOriginRequired {
		req.Header.Set("Origin", h.Origin)
	}
	if policyToApply.IsCookieRequired {
		req.Header.Set(h.Cookie.Name, h.Cookie.Value)
	}
	return req
}

var securityPolicies = getPolicy()

type Operation int

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
	GetTags
	WipeData
)

func getPolicy() SecurityPolicyCollection {
	return SecurityPolicyCollection{
		policies: []*SecurityPolicy{
			{false, false, false, []Operation{FindApps, DownloadApp, GetTags}},
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
