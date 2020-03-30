package api

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
	scheme := "https"
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
		scheme = "http"
	}
	return &ServiceHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     scheme,
	}
}

// NewAuthenticatedServiceHandler returns a new ServiceHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedServiceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ServiceHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

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
func (s *ServiceHandler) CreateService(project string, service models.CreateService) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(service)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(s.Scheme+"://"+s.getBaseURL()+"/v1/project/"+project+"/service", bodyStr, s)
}
