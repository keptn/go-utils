package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const v1EventPath = "/v1/event"
const v1MetadataPath = "/v1/metadata"

type APIV1Interface interface {
	// SendEvent sends an event to Keptn.
	SendEvent(event models.KeptnContextExtendedCE) (*models.EventContext, *models.Error)

	// TriggerEvaluation triggers a new evaluation.
	TriggerEvaluation(project string, stage string, service string, evaluation models.Evaluation) (*models.EventContext, *models.Error)

	// CreateProject creates a new project.
	CreateProject(project models.CreateProject) (string, *models.Error)

	// UpdateProject updates a project.
	UpdateProject(project models.CreateProject) (string, *models.Error)

	// DeleteProject deletes a project.
	DeleteProject(project models.Project) (*models.DeleteProjectResponse, *models.Error)

	// CreateService creates a new service.
	CreateService(project string, service models.CreateService) (string, *models.Error)

	// DeleteService deletes a service.
	DeleteService(project string, service string) (*models.DeleteServiceResponse, *models.Error)

	// GetMetadata retrieves Keptn metadata information.
	GetMetadata() (*models.Metadata, *models.Error)
}

// APIHandler handles projects
type APIHandler struct {
	apiHandler *v2.APIHandler
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewAPIHandler returns a new APIHandler
func NewAPIHandler(baseURL string) *APIHandler {
	return NewAPIHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewAPIHandlerWithHTTPClient returns a new APIHandler that uses the specified http.Client
func NewAPIHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *APIHandler {
	return &APIHandler{
		BaseURL:    httputils.TrimHTTPScheme(baseURL),
		HTTPClient: httpClient,
		Scheme:     "http",
		apiHandler: v2.NewAPIHandlerWithHTTPClient(baseURL, httpClient),
	}
}

// NewAuthenticatedAPIHandler returns a new APIHandler that authenticates at the api-service endpoint via the provided token
// Deprecated: use APISet instead
func NewAuthenticatedAPIHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *APIHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedAPIHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedAPIHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *APIHandler {
	v2APIHandler := v2.NewAuthenticatedAPIHandler(baseURL, authToken, authHeader, httpClient, scheme)

	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return &APIHandler{
		BaseURL:    httputils.TrimHTTPScheme(baseURL),
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
		apiHandler: v2APIHandler,
	}
}

func (a *APIHandler) getBaseURL() string {
	return a.BaseURL
}

func (a *APIHandler) getAuthToken() string {
	return a.AuthToken
}

func (a *APIHandler) getAuthHeader() string {
	return a.AuthHeader
}

func (a *APIHandler) getHTTPClient() *http.Client {
	return a.HTTPClient
}

// SendEvent sends an event to Keptn.
func (a *APIHandler) SendEvent(event models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
	return a.apiHandler.SendEvent(context.TODO(), event, v2.APISendEventOptions{})
}

// TriggerEvaluation triggers a new evaluation.
func (a *APIHandler) TriggerEvaluation(project, stage, service string, evaluation models.Evaluation) (*models.EventContext, *models.Error) {
	return a.apiHandler.TriggerEvaluation(context.TODO(), project, stage, service, evaluation, v2.APITriggerEvaluationOptions{})
}

// CreateProject creates a new project.
func (a *APIHandler) CreateProject(project models.CreateProject) (string, *models.Error) {
	return a.apiHandler.CreateProject(context.TODO(), project, v2.APICreateProjectOptions{})
}

// UpdateProject updates a project.
func (a *APIHandler) UpdateProject(project models.CreateProject) (string, *models.Error) {
	return a.apiHandler.UpdateProject(context.TODO(), project, v2.APIUpdateProjectOptions{})
}

// DeleteProject deletes a project.
func (a *APIHandler) DeleteProject(project models.Project) (*models.DeleteProjectResponse, *models.Error) {
	return a.apiHandler.DeleteProject(context.TODO(), project, v2.APIDeleteProjectOptions{})
}

// CreateService creates a new service.
func (a *APIHandler) CreateService(project string, service models.CreateService) (string, *models.Error) {
	return a.apiHandler.CreateService(context.TODO(), project, service, v2.APICreateServiceOptions{})
}

// DeleteService deletes a service.
func (a *APIHandler) DeleteService(project, service string) (*models.DeleteServiceResponse, *models.Error) {
	return a.apiHandler.DeleteService(context.TODO(), project, service, v2.APIDeleteServiceOptions{})
}

// GetMetadata retrieves Keptn metadata information.
func (a *APIHandler) GetMetadata() (*models.Metadata, *models.Error) {
	return a.apiHandler.GetMetadata(context.TODO(), v2.APIGetMetadataOptions{})
}
