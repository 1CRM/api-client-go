package endpoints

import (
	"context"

	api "github.com/1CRM/api-client-go"
)

// Client is a wrapper around api.Client that adds new methods
type Client struct {
	api.Client
}

// NewClient creates a new Client by wrapping api.Client
func NewClient(ctx context.Context, url string, auth api.Auth) *Client {
	cl := api.NewClient(ctx, url, auth)
	return &Client{
		*cl,
	}
}

// Files retuns an instance of Files interface that can be used to perform file
// uploads and downloads
func (cl *Client) Files() Files {
	return &filesImpl{
		client: cl,
	}
}
