package httputils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Downloader struct {
	Timeout time.Duration
}

type DownloaderOption func(*Downloader)

// WithTimeout sets a timeout on the underlying HTTP client
func WithTimeout(timeout time.Duration) DownloaderOption {
	return func(d *Downloader) {
		d.Timeout = timeout
	}
}

// NewDownloader creates a new Downloader to download files from an URL
func NewDownloader(opts ...DownloaderOption) *Downloader {
	d := &Downloader{}

	for _, opt := range opts {
		opt(d)
	}
	return d
}

// DownloadFromURL downloads a file from an URL and returns the content
// as a slice of bytes
func (d Downloader) DownloadFromURL(URL string) ([]byte, error) {
	if !IsValidURL(URL) {
		return nil, fmt.Errorf("%s is not a valid URL", URL)
	}
	c := http.Client{
		Timeout: d.Timeout,
	}
	resp, err := c.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// DownloadFromURL downloads a file from an URL and returns the content
// as a slice of bytes
func DownloadFromURL(URL string) ([]byte, error) {
	return NewDownloader().DownloadFromURL(URL)
}

// IsValidURL checks whether the given string is a valid URL or not
func IsValidURL(strURL string) bool {
	_, err := url.ParseRequestURI(strURL)
	if err != nil {
		return false
	}
	u, err := url.Parse(strURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
