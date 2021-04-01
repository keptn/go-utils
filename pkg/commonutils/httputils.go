package commonutils

import (
	"io/ioutil"
	"net/http"
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
