package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client is used to send requests to the API
type Client struct {
	auth Auth
	url  string
	ctx  context.Context
}

type requestOptions struct {
	body        io.Reader
	query       url.Values
	contentType string
	ctx         context.Context
}

// RequestOption is a function used to set options of a request to the API
type RequestOption func(*requestOptions) error

// NewClient returns a new API Client
func NewClient(ctx context.Context, url string, auth Auth) *Client {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Client{
		url:  strings.TrimRight(url, "/"),
		auth: auth,
		ctx:  ctx,
	}
}

// WithContentType is a RequestOption that sets the Content-Type header of the HTTP request
func WithContentType(contentType string) RequestOption {
	return func(opts *requestOptions) error {
		opts.contentType = contentType
		return nil
	}
}

// WithBody is a RequestOption that sets the body of the HTTP request
func WithBody(body io.Reader) RequestOption {
	return func(opts *requestOptions) error {
		if prev := opts.body; prev != nil && prev.(io.Closer) != nil && prev != body {
			prev.(io.Closer).Close()
		}
		opts.body = body
		return nil
	}
}

// WithQuery is a RequestOption that sets the query part of the HTTP request URL
func WithQuery(query url.Values) RequestOption {
	return func(opts *requestOptions) error {
		opts.query = query
		return nil
	}
}

// WithQueryValue is a RequestOption that sets one parameter of the the query
// part of the HTTP request URL
func WithQueryValue(key, value string, append bool) RequestOption {
	return func(opts *requestOptions) error {
		if opts.query == nil {
			opts.query = make(url.Values)
		}
		if append {
			opts.query.Add(key, value)
		} else {
			opts.query.Set(key, value)
		}
		return nil
	}
}

// WithJSONBody is a RequestOption that sets the request's body to a JSON string
func WithJSONBody(content interface{}) RequestOption {
	return func(opts *requestOptions) error {
		b, err := json.Marshal(content)
		if err != nil {
			return err
		}
		if prev := opts.body; prev != nil && prev.(io.Closer) != nil {
			prev.(io.Closer).Close()
		}
		opts.body = bytes.NewBuffer(b)
		return nil
	}
}

// WitchContext is a RequestOption that sets the request's context
func WitchContext(ctx context.Context) RequestOption {
	return func(opts *requestOptions) error {
		opts.ctx = ctx
		return nil
	}
}

// Request sends an HTTP request with arbitrary HTTP method
func (c *Client) Request(method string, endpoint string, options ...RequestOption) (*Response, error) {
	opts := requestOptions{
		contentType: "application/json",
		ctx:         c.ctx,
	}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}
	u, err := url.Parse(c.url + "/" + endpoint)
	if err != nil {
		return nil, err
	}
	if opts.query != nil {
		u.RawQuery = opts.query.Encode()
	}
	req, err := http.NewRequest(method, u.String(), opts.body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", opts.contentType)
	if c.auth != nil {
		err = c.auth.ApplyRequestOptions(req)
		if err != nil {
			return nil, err
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, &resError{code: res.StatusCode, err: string(b)}
	}
	return &Response{HttpResponse: res}, nil
}

// Get is convenience wrapper around Request() to send a GET request
func (c *Client) Get(endpoint string, options ...RequestOption) (*Response, error) {
	return c.Request("GET", endpoint, options...)
}

// Post is convenience wrapper around Request() to send a POST request
func (c *Client) Post(endpoint string, options ...RequestOption) (*Response, error) {
	return c.Request("POST", endpoint, options...)
}

// Patch is convenience wrapper around Request() to send a Patch request
func (c *Client) Patch(endpoint string, options ...RequestOption) (*Response, error) {
	return c.Request("PATCH", endpoint, options...)
}

// Put is convenience wrapper around Request() to send a PUT request
func (c *Client) Put(endpoint string, options ...RequestOption) (*Response, error) {
	return c.Request("PUT", endpoint, options...)
}
