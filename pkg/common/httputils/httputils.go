package httputils

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// DownloadFromURL downloads a file from the given url and returns
// its content as a slice of bytes
func DownloadFromURL(URL string) ([]byte, error) {
	c := http.Client{
		Timeout: 5 * time.Second,
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
