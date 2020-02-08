package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	auth Auth
	url  string
}

type requestOptions struct {
	body        io.Reader
	query       url.Values
	contentType string
}

type RequestOption func(*requestOptions) error

func NewClient(url string, auth Auth) *Client {
	return &Client{
		url:  strings.TrimRight(url, "/"),
		auth: auth,
	}
}

func WithContentType(contentType string) RequestOption {
	return func(opts *requestOptions) error {
		opts.contentType = contentType
		return nil
	}
}

func WithBody(body io.Reader) RequestOption {
	return func(opts *requestOptions) error {
		if prev := opts.body; prev != nil && prev.(io.Closer) != nil && prev != body {
			prev.(io.Closer).Close()
		}
		opts.body = body
		return nil
	}
}

func WithQuery(query url.Values) RequestOption {
	return func(opts *requestOptions) error {
		opts.query = query
		return nil
	}
}

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

func WithJsonBody(content interface{}) RequestOption {
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

func (c *Client) Request(method string, endpoint string, options ...RequestOption) (*Response, error) {
	opts := requestOptions{
		contentType: "application/json",
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

func (c *Client) Get(endpoint string, options ...RequestOption) (*Response, error) {
	return c.Request("GET", endpoint, options...)
}

func (c *Client) Post(endpoint string, options ...RequestOption) (*Response, error) {
	return c.Request("POST", endpoint, options...)
}

func (c *Client) Put(endpoint string, options ...RequestOption) (*Response, error) {
	return c.Request("PUT", endpoint, options...)
}
