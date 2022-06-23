package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

type ServicesV1Interface interface {
	// CreateServiceInStage creates a new service.
	CreateServiceInStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)

	// DeleteServiceFromStage deletes a service from a stage.
	DeleteServiceFromStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)

	// GetService gets a service.
	GetService(project, stage, service string) (*models.Service, error)

	// GetAllServices returns a list of all services.
	GetAllServices(project string, stage string) ([]*models.Service, error)
}

// ServiceHandler handles services
type ServiceHandler struct {
	serviceHandler *v2.ServiceHandler
	BaseURL        string
	AuthToken      string
	AuthHeader     string
	HTTPClient     *http.Client
	Scheme         string
}

// NewServiceHandler returns a new ServiceHandler which sends all requests directly to the configuration-service
func NewServiceHandler(baseURL string) *ServiceHandler {
	return NewServiceHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewServiceHandlerWithHTTPClient returns a new ServiceHandler which sends all requests directly to the configuration-service using the specified http.Client
func NewServiceHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *ServiceHandler {
	return &ServiceHandler{
		BaseURL:        httputils.TrimHTTPScheme(baseURL),
		HTTPClient:     httpClient,
		Scheme:         "http",
		serviceHandler: v2.NewServiceHandlerWithHTTPClient(baseURL, httpClient),
	}
}

// NewAuthenticatedServiceHandler returns a new ServiceHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
// Deprecated: use APISet instead
func NewAuthenticatedServiceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ServiceHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedServiceHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedServiceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ServiceHandler {
	v2ServiceHandler := v2.NewAuthenticatedServiceHandler(baseURL, authToken, authHeader, httpClient, scheme)

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return &ServiceHandler{
		BaseURL:        httputils.TrimHTTPScheme(baseURL),
		AuthHeader:     authHeader,
		AuthToken:      authToken,
		HTTPClient:     httpClient,
		Scheme:         scheme,
		serviceHandler: v2ServiceHandler,
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

// CreateServiceInStage creates a new service.
func (s *ServiceHandler) CreateServiceInStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error) {
	s.ensureHandlerIsSet()
	return s.serviceHandler.CreateServiceInStage(context.TODO(), project, stage, serviceName, v2.ServicesCreateServiceInStageOptions{})
}

// DeleteServiceFromStage deletes a service from a stage.
func (s *ServiceHandler) DeleteServiceFromStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error) {
	s.ensureHandlerIsSet()
	return s.serviceHandler.DeleteServiceFromStage(context.TODO(), project, stage, serviceName, v2.ServicesDeleteServiceFromStageOptions{})
}

// GetService gets a service.
func (s *ServiceHandler) GetService(project, stage, service string) (*models.Service, error) {
	s.ensureHandlerIsSet()
	return s.serviceHandler.GetService(context.TODO(), project, stage, service, v2.ServicesGetServiceOptions{})
}

// GetAllServices returns a list of all services.
func (s *ServiceHandler) GetAllServices(project string, stage string) ([]*models.Service, error) {
	s.ensureHandlerIsSet()
	return s.serviceHandler.GetAllServices(context.TODO(), project, stage, v2.ServicesGetAllServicesOptions{})
}

func (s *ServiceHandler) ensureHandlerIsSet() {
	if s.serviceHandler != nil {
		return
	}

	if s.AuthToken != "" {
		s.serviceHandler = v2.NewAuthenticatedServiceHandler(s.BaseURL, s.AuthToken, s.AuthHeader, s.HTTPClient, s.Scheme)
	} else {
		s.serviceHandler = v2.NewServiceHandlerWithHTTPClient(s.BaseURL, s.HTTPClient)
	}
}
