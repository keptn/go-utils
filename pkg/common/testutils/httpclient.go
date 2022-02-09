package testutils

import "net/http"

// HTTPClientMock is mock implementation of http client usable for testing
type HTTPClientMock struct {
	DoFunc func(r *http.Request) (*http.Response, error)
}

func (h HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	return h.DoFunc(r)
}
