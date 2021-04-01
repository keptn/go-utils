package commonutils

import (
	"io/ioutil"
	"net/http"
	"time"
)

func DownloadFromURL(url string) ([]byte, error) {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Get(url)
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
