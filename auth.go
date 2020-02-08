package client

import "net/http"

// Auth is an interface that provides authentication info
type Auth interface {
	ApplyRequestOptions(*http.Request) error
}
