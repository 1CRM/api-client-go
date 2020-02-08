package client

import (
	"net/http"
)

// BasicAuth provides Basic authentication
type BasicAuth struct {
	username string
	password string
}

// NewBasicAuth returns a new BasicAuth
func NewBasicAuth(username, password string) *BasicAuth {
	auth := BasicAuth{
		username: username,
		password: password,
	}
	return &auth
}

// ApplyRequestOptions implements Auth.ApplyRequestOptions
func (auth *BasicAuth) ApplyRequestOptions(req *http.Request) error {
	req.SetBasicAuth(auth.username, auth.password)
	return nil
}
