package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	HttpResponse *http.Response
}

func (res *Response) String() (string, error) {
	body := res.HttpResponse.Body
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (res *Response) ParseJSON(out interface{}) error {
	body := res.HttpResponse.Body
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
