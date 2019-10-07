package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

// ServiceHandler handles services
type ServiceHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewServiceHandler returns a new ServiceHandler
func NewServiceHandler(baseURL string) *ServiceHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &ServiceHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{},
		Scheme:     "http",
	}
}

// NewAuthenticatedServiceHandler returns a new ServiceHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedServiceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ServiceHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &ServiceHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (s *ServiceHandler) getBaseURL() string {
	return s.BaseURL
}

func (s *ServiceHandler) getAuthToken() string {
	return s.AuthToken
}

func (s *ServiceHandler) getAuthHeader() string {
	return s.AuthHeader
}

func (s *ServiceHandler) getHTTPClient() *http.Client {
	return s.HTTPClient
}

// CreateService creates a new service
func (s *ServiceHandler) CreateService(project string, stage string, serviceName string) (*models.Error, error) {

	service := models.Service{ServiceName: serviceName}
	body, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}
	return post(s.Scheme+"://"+s.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service", body, s)
}
