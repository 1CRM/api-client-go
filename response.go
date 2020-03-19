package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Response is a response returned from requests to the API
type Response struct {
	HTTPResponse *http.Response
}

func (res *Response) String() (string, error) {
	body := res.HTTPResponse.Body
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ParseJSON parses response body into supplied argument
func (res *Response) ParseJSON(out interface{}) error {
	body := res.HTTPResponse.Body
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
