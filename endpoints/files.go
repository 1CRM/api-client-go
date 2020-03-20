package endpoints

import (
	"fmt"
	"io"

	api "github.com/1CRM/api-client-go"
)

// Files is an interface providing file-related oprations
type Files interface {
	Download(model string, id string, options ...api.RequestOption) (io.ReadCloser, error)
	Metadata(model string, id string, options ...api.RequestOption) (*FileMetadata, error)
	Upload(name string, content io.Reader, options ...api.RequestOption) (string, error)
}

// FileMetadata represents metadata of Document, DocumentRevision or a Note attachment
type FileMetadata struct {
	Name     string `json:"name"`
	Size     int    `json:"size"`
	MIMEType string `json:"mime_type"`
	Modified int    `json:"modified"`
}

type filesImpl struct {
	client *Client
}

type uploadResult struct {
	ID string `json:"id"`
}

// Download implements Files.Download
func (files *filesImpl) Download(model string, id string, options ...api.RequestOption) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("files/download/%s/%s", model, id)
	res, err := files.client.Get(endpoint, options...)
	if err != nil {
		return nil, err
	}
	return res.HTTPResponse.Body, nil
}

// Metadata implements Files.Metadata
func (files *filesImpl) Metadata(model string, id string, options ...api.RequestOption) (*FileMetadata, error) {
	endpoint := fmt.Sprintf("files/info/%s/%s", model, id)
	res, err := files.client.Get(endpoint, options...)
	if err != nil {
		return nil, err
	}
	var meta FileMetadata
	if err := res.ParseJSON(&meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// Upload implements Files.Upload
func (files *filesImpl) Upload(name string, content io.Reader, options ...api.RequestOption) (string, error) {
	endpoint := "files/upload"
	res, err := files.client.Post(endpoint,
		append(options,
			api.WithHeader("X-OneCRM-Filename", name),
			api.WithBody(content),
		)...,
	)
	if err != nil {
		return "", err
	}
	var result uploadResult
	if err = res.ParseJSON(&result); err != nil {
		return "", err
	}
	return result.ID, nil
}
