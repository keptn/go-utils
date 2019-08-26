package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/keptn/go-utils/pkg/models"
)

// ConfigService represents the interface for accessing the configuration service
type ConfigService interface {
	getBaseURL() string
	getAuthToken() string
	getAuthHeader() string
}

func post(uri string, data []byte, c ConfigService) (*models.Error, error) {

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, c)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respErr models.Error
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return nil, err
	}

	return &respErr, nil
}

func addAuthHeader(req *http.Request, c ConfigService) {
	if c.getAuthHeader() != "" && c.getAuthToken() != "" {
		req.Header.Set(c.getAuthHeader(), c.getAuthToken())
	}
}
