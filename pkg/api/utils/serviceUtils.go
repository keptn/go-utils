package api

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

type ServicesV1Interface interface {
	// CreateServiceInStage creates a new service.
	CreateServiceInStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)

	// CreateServiceInStageWithContext creates a new service.
	CreateServiceInStageWithContext(ctx context.Context, project string, stage string, serviceName string) (*models.EventContext, *models.Error)

	// DeleteServiceFromStage deletes a service from a stage.
	DeleteServiceFromStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)

	// DeleteServiceFromStageWithContext deletes a service from a stage.
	DeleteServiceFromStageWithContext(ctx context.Context, project string, stage string, serviceName string) (*models.EventContext, *models.Error)

	// GetService gets a service.
	GetService(project, stage, service string) (*models.Service, error)

	// GetServiceWithContext gets a service.
	GetServiceWithContext(ctx context.Context, project, stage, service string) (*models.Service, error)

	// GetAllServices returns a list of all services.
	GetAllServices(project string, stage string) ([]*models.Service, error)

	// GetAllServicesWithContext returns a list of all services.
	GetAllServicesWithContext(ctx context.Context, project string, stage string) ([]*models.Service, error)
}

// ServiceHandler handles services
type ServiceHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewServiceHandler returns a new ServiceHandler which sends all requests directly to the configuration-service
func NewServiceHandler(baseURL string) *ServiceHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &ServiceHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))},
		Scheme:     "http",
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
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")

	baseURL = strings.TrimRight(baseURL, "/")

	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

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

// CreateServiceInStage creates a new service.
func (s *ServiceHandler) CreateServiceInStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error) {
	return s.CreateServiceInStageWithContext(context.TODO(), project, stage, serviceName)
}

// CreateServiceInStageWithContext creates a new service.
func (s *ServiceHandler) CreateServiceInStageWithContext(ctx context.Context, project string, stage string, serviceName string) (*models.EventContext, *models.Error) {
	service := models.Service{ServiceName: serviceName}
	body, err := service.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, s.Scheme+"://"+s.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService, body, s)
}

// DeleteServiceFromStage deletes a service from a stage.
func (s *ServiceHandler) DeleteServiceFromStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error) {
	return s.DeleteServiceFromStageWithContext(context.TODO(), project, stage, serviceName)
}

// DeleteServiceFromStageWithContext deletes a service from a stage.
func (s *ServiceHandler) DeleteServiceFromStageWithContext(ctx context.Context, project string, stage string, serviceName string) (*models.EventContext, *models.Error) {
	return deleteWithEventContext(ctx, s.Scheme+"://"+s.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+serviceName, s)
}

// GetService gets a service.
func (s *ServiceHandler) GetService(project, stage, service string) (*models.Service, error) {
	return s.GetServiceWithContext(context.TODO(), project, stage, service)
}

// GetServiceWithContext gets a service.
func (s *ServiceHandler) GetServiceWithContext(ctx context.Context, project, stage, service string) (*models.Service, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + v1ProjectPath + "/" + project + pathToStage + "/" + stage + pathToService + "/" + service)
	if err != nil {
		return nil, err
	}

	body, mErr := getAndExpectOK(ctx, url.String(), s)
	if mErr != nil {
		return nil, mErr.ToError()
	}

	received := &models.Service{}
	if err = received.FromJSON(body); err != nil {
		return nil, err
	}
	return received, nil
}

// GetAllServices returns a list of all services.
func (s *ServiceHandler) GetAllServices(project string, stage string) ([]*models.Service, error) {
	return s.GetAllServicesWithContext(context.TODO(), project, stage)
}

// GetAllServicesWithContext returns a list of all services.
func (s *ServiceHandler) GetAllServicesWithContext(ctx context.Context, project string, stage string) ([]*models.Service, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	services := []*models.Service{}

	nextPageKey := ""

	for {
		url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + v1ProjectPath + "/" + project + pathToStage + "/" + stage + pathToService)
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}

		body, mErr := getAndExpectOK(ctx, url.String(), s)
		if mErr != nil {
			return nil, mErr.ToError()
		}

		received := &models.Services{}
		if err = received.FromJSON(body); err != nil {
			return nil, err
		}
		services = append(services, received.Services...)

		if received.NextPageKey == "" || received.NextPageKey == "0" {
			break
		}
		nextPageKey = received.NextPageKey
	}

	return services, nil
}
