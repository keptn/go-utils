package v2

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

// ServicesCreateServiceInStageOptions are options for ServicesInterface.CreateServiceInStage().
type ServicesCreateServiceInStageOptions struct{}

// ServicesDeleteServiceFromStageOptions are options for ServicesInterface.DeleteServiceFromStage().
type ServicesDeleteServiceFromStageOptions struct{}

// ServicesGetServiceOptions are options for ServicesInterface.GetService().
type ServicesGetServiceOptions struct{}

// ServicesGetAllServicesOptions are options for ServicesInterface.GetAllServices().
type ServicesGetAllServicesOptions struct{}

type ServicesInterface interface {

	// CreateServiceInStage creates a new service.
	CreateServiceInStage(ctx context.Context, project string, stage string, serviceName string, opts ServicesCreateServiceInStageOptions) (*models.EventContext, *models.Error)

	// DeleteServiceFromStage deletes a service from a stage.
	DeleteServiceFromStage(ctx context.Context, project string, stage string, serviceName string, opts ServicesDeleteServiceFromStageOptions) (*models.EventContext, *models.Error)

	// GetService gets a service.
	GetService(ctx context.Context, project, stage, service string, opts ServicesGetServiceOptions) (*models.Service, error)

	// GetAllServices returns a list of all services.
	GetAllServices(ctx context.Context, project string, stage string, opts ServicesGetAllServicesOptions) ([]*models.Service, error)
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
func (s *ServiceHandler) CreateServiceInStage(ctx context.Context, project string, stage string, serviceName string, opts ServicesCreateServiceInStageOptions) (*models.EventContext, *models.Error) {
	service := models.Service{ServiceName: serviceName}
	body, err := service.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, s.Scheme+"://"+s.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService, body, s)
}

// DeleteServiceFromStage deletes a service from a stage.
func (s *ServiceHandler) DeleteServiceFromStage(ctx context.Context, project string, stage string, serviceName string, opts ServicesDeleteServiceFromStageOptions) (*models.EventContext, *models.Error) {
	return deleteWithEventContext(ctx, s.Scheme+"://"+s.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+serviceName, s)
}

// GetService gets a service.
func (s *ServiceHandler) GetService(ctx context.Context, project, stage, service string, opts ServicesGetServiceOptions) (*models.Service, error) {
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
func (s *ServiceHandler) GetAllServices(ctx context.Context, project string, stage string, opts ServicesGetAllServicesOptions) ([]*models.Service, error) {

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
