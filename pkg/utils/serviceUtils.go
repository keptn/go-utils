package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/keptn/go-utils/pkg/models"
)

// ServiceHandler handles services
type ServiceHandler struct {
	BaseURL string
}

// NewServiceHandler returns a new ServiceHandler
func NewServiceHandler(baseURL string) *ServiceHandler {
	return &ServiceHandler{
		BaseURL: baseURL,
	}
}

// GetService returns a list of services.
func (r *ServiceHandler) GetService(project string, stage string, pageSize int, nextPageKey string) (*models.Services, error) {
	url, err := url.Parse("http://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/service")
	if err != nil {
		return nil, err
	}
	q := url.Query()
	q.Set("pageSize", strconv.Itoa(pageSize))
	if nextPageKey != "" {
		q.Set("nextPageKey", nextPageKey)
	}
	req, err := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var services models.Services
	err = json.Unmarshal(body, &services)
	return &services, err
}
