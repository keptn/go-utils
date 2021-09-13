package httputils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
		Timeout:   d.Timeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	// TODO: NewRequestWithContext in order to get proper traces
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

// TrimScheme takes a string of an URL and removes the leading scheme (http or https)
func TrimHTTPScheme(strURL string) string {
	trimmedURL := strURL
	if strings.HasPrefix(strURL, "https://") {
		trimmedURL = strings.TrimPrefix(strURL, "https://")
	} else if strings.HasPrefix(strURL, "http://") {
		trimmedURL = strings.TrimPrefix(strURL, "http://")
	}
	return trimmedURL
}
