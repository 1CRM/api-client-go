package client

import (
	"context"
	"net/url"
	"os"
)

// AuthFlow contains parameters of OAuth2 flow
type AuthFlow struct {
	clientID     string
	clientSecret string
	redirectURI  string
	username     string
	password     string
	scope        string
	ownerType    string
	state        string
	url          string
}

// AuthFlowOption provdes parameters for AuthFlow
type AuthFlowOption func(*AuthFlow)

// ResourceOwnerRequest is used to initiate Resource Owner flow
type ResourceOwnerRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	Scope        string `json:"scope`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

// ClientCredentialsRequest is used to initiate Client Credentials flow
type ClientCredentialsRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	Scope        string `json:"scope`
	ClientSecret string `json:"client_secret"`
}

// AuthCodeRequest is used to initiate Authorization Code flow
type AuthCodeRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	Scope        string `json:"scope`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

// WithClientID is used to set Client ID for Oauth2 flow
func WithClientID(id string) func(*AuthFlow) {
	return func(p *AuthFlow) {
		p.clientID = id
	}
}

// WithClientSecret is used to set Client Secret for Oauth2 flow
func WithClientSecret(secret string) func(*AuthFlow) {
	return func(p *AuthFlow) {
		p.clientSecret = secret
	}
}

// WithRedirectURI is used to set redirect URI for Oauth2 flow
func WithRedirectURI(uri string) func(*AuthFlow) {
	return func(p *AuthFlow) {
		p.redirectURI = uri
	}
}

// WithUserName is used to set User Name for Oauth2 flow
func WithUserName(username string) func(*AuthFlow) {
	return func(p *AuthFlow) {
		p.username = username
	}
}

// WithPassword is used to set password for Oauth2 flow
func WithPassword(password string) func(*AuthFlow) {
	return func(p *AuthFlow) {
		p.password = password
	}
}

// WithScope is used to set Scope for Oauth2 flow
func WithScope(scope string) func(*AuthFlow) {
	return func(p *AuthFlow) {
		p.scope = scope
	}
}

// WithState is used to set State for Oauth2 flow
func WithState(state string) func(*AuthFlow) {
	return func(p *AuthFlow) {
		p.state = state
	}
}

// WithOwnerType is used to set Owner Type for Oauth2 flow
func WithOwnerType(ownerType string) func(*AuthFlow) {
	return func(p *AuthFlow) {
		p.ownerType = ownerType
	}
}

// NewAuthFlow returns a AuthFlow
func NewAuthFlow(url string, params ...AuthFlowOption) *AuthFlow {
	flow := AuthFlow{
		scope:     "profile",
		ownerType: "user",
		url:       url,
	}
	flow.clientID, _ = os.LookupEnv("ONECRM_CLIENT_ID")
	flow.clientSecret, _ = os.LookupEnv("ONECRM_CLIENT_SECRET")
	flow.redirectURI, _ = os.LookupEnv("ONECRM_REDIRECT_URI")
	flow.username, _ = os.LookupEnv("ONECRM_USERNAME")
	flow.password, _ = os.LookupEnv("ONECRM_PASSWORD")
	for _, p := range params {
		p(&flow)
	}
	return &flow
}

// InitAuthCode is used to initiate Authorization Code flow.
// It returns an URL the user is to be redirected to in order to complete the flow.
func (flow *AuthFlow) InitAuthCode() (string, error) {
	u, err := url.Parse(flow.url + "/auth/" + flow.ownerType + "/authorize")
	if err != nil {
		return "", err
	}
	values := make(url.Values)
	values.Set("response_type", "code")
	values.Set("client_id", flow.clientID)
	values.Set("redirect_uri", flow.redirectURI)
	values.Set("state", flow.state)
	u.RawQuery = values.Encode()
	return u.String(), nil
}

// FinalizeAuthCode completes the Authorization Code flow and returns an access token
func (flow *AuthFlow) FinalizeAuthCode(code string, ctx context.Context) (*OAuth2AccessToken, error) {
	body := AuthCodeRequest{
		GrantType:    "authorization_code",
		ClientID:     flow.clientID,
		ClientSecret: flow.clientSecret,
		Scope:        flow.scope,
		RedirectURI:  flow.redirectURI,
		Code:         code,
	}
	c := NewClient(flow.url, nil, ctx)
	res, err := c.Post(
		"auth/"+flow.ownerType+"/access_token",
		WithJsonBody(body),
	)
	if err != nil {
		return nil, err
	}
	var token OAuth2AccessToken
	if err = res.ParseJSON(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

// InitResourceOwner returns an access token for resource owner (user or contact)
func (flow *AuthFlow) InitResourceOwner(ctx context.Context) (*OAuth2AccessToken, error) {
	body := ResourceOwnerRequest{
		GrantType:    "password",
		ClientID:     flow.clientID,
		ClientSecret: flow.clientSecret,
		Scope:        flow.scope,
		Username:     flow.username,
		Password:     flow.password,
	}
	c := NewClient(flow.url, nil, ctx)
	res, err := c.Post(
		"auth/"+flow.ownerType+"/access_token",
		WithJsonBody(body),
	)
	if err != nil {
		return nil, err
	}
	var token OAuth2AccessToken
	if err = res.ParseJSON(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

// InitClientCredentials returns an access token using the client credentials
func (flow *AuthFlow) InitClientCredentials(ctx context.Context) (*OAuth2AccessToken, error) {
	body := ClientCredentialsRequest{
		GrantType:    "client_credentials",
		ClientID:     flow.clientID,
		ClientSecret: flow.clientSecret,
		Scope:        flow.scope,
	}
	c := NewClient(flow.url, nil, ctx)
	res, err := c.Post(
		"auth/"+flow.ownerType+"/access_token",
		WithJsonBody(body),
	)
	if err != nil {
		return nil, err
	}
	var token OAuth2AccessToken
	if err = res.ParseJSON(&token); err != nil {
		return nil, err
	}
	return &token, nil
}
