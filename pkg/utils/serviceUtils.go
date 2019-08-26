package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/keptn/go-utils/pkg/models"
)

// ServiceHandler handles services
type ServiceHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
}

// NewServiceHandler returns a new ServiceHandler
func NewServiceHandler(baseURL string) *ServiceHandler {
	return &ServiceHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
	}
}

// NewAuthenticatedServiceHandler returns a new ServiceHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedServiceHandler(baseURL string, authToken string, authHeader string) *ServiceHandler {
	return &ServiceHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
	}
}

// GetAllServices returns a list of all services.
func (r *ServiceHandler) GetAllServices(project string, stage string) ([]*models.Service, error) {

	services := []*models.Service{}

	nextPageKey := ""

	for {
		url, err := url.Parse("http://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/service")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		if r.AuthHeader != "" && r.AuthToken != "" {
			req.Header.Set(r.AuthHeader, r.AuthToken)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var received models.Services
		err = json.Unmarshal(body, &services)
		if err != nil {
			return nil, err
		}
		services = append(services, received.Services...)

		if received.NextPageKey == "" {
			break
		}
		nextPageKey = received.NextPageKey
	}

	return services, nil
}
