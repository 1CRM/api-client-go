package client

import (
	"net/http"
)

// OAuth2AccessToken is a token used for OAuth 2.0
type OAuth2AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// ApplyRequestOptions implements Auth.ApplyRequestOptions
func (auth *OAuth2AccessToken) ApplyRequestOptions(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	return nil
}
