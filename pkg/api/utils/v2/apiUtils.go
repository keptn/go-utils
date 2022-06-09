package v2

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const v1EventPath = "/v1/event"
const v1MetadataPath = "/v1/metadata"

// APISendEventOptions are options for APIInterface.SendEvent().
type APISendEventOptions struct{}

// APITriggerEvaluationOptions are options for APIInterface.TriggerEvaluation().
type APITriggerEvaluationOptions struct{}

// APICreateProjectOptions are options for APIInterface.CreateProject().
type APICreateProjectOptions struct{}

// APIUpdateProjectOptions are options for APIInterface.UpdateProject().
type APIUpdateProjectOptions struct{}

// APIDeleteProjectOptions are options for APIInterface.DeleteProject().
type APIDeleteProjectOptions struct{}

// APICreateServiceOptions are options for APIInterface.CreateService().
type APICreateServiceOptions struct{}

// APIDeleteServiceOptions are options for APIInterface.DeleteService().
type APIDeleteServiceOptions struct{}

// APIGetMetadataOptions are options for APIInterface.GetMetadata().
type APIGetMetadataOptions struct{}

type APIInterface interface {
	// SendEvent sends an event to Keptn.
	SendEvent(ctx context.Context, event models.KeptnContextExtendedCE, opts APISendEventOptions) (*models.EventContext, *models.Error)

	// TriggerEvaluation triggers a new evaluation.
	TriggerEvaluation(ctx context.Context, project string, stage string, service string, evaluation models.Evaluation, opts APITriggerEvaluationOptions) (*models.EventContext, *models.Error)

	// CreateProject creates a new project.
	CreateProject(ctx context.Context, project models.CreateProject, opts APICreateProjectOptions) (string, *models.Error)

	// UpdateProject updates a project.
	UpdateProject(ctx context.Context, project models.CreateProject, opts APIUpdateProjectOptions) (string, *models.Error)

	// DeleteProject deletes a project.
	DeleteProject(ctx context.Context, project models.Project, opts APIDeleteProjectOptions) (*models.DeleteProjectResponse, *models.Error)

	// CreateService creates a new service.
	CreateService(ctx context.Context, project string, service models.CreateService, opts APICreateServiceOptions) (string, *models.Error)

	// DeleteService deletes a service.
	DeleteService(ctx context.Context, project string, service string, opts APIDeleteServiceOptions) (*models.DeleteServiceResponse, *models.Error)

	// GetMetadata retrieves Keptn metadata information.
	GetMetadata(ctx context.Context, opts APIGetMetadataOptions) (*models.Metadata, *models.Error)
}

type APIHandler struct {
	baseURL    string
	authToken  string
	authHeader string
	httpClient *http.Client
	scheme     string
}

// NewAPIHandler returns a new APIHandler
func NewAPIHandler(baseURL string) *APIHandler {
	return NewAPIHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewAPIHandlerWithHTTPClient returns a new APIHandler using the specified http.Client
func NewAPIHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *APIHandler {
	return createAPIHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedAPIHandler returns a new APIHandler that authenticates at the api-service endpoint via the provided token
func NewAuthenticatedAPIHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *APIHandler {
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return createAPIHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAPIHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *APIHandler {
	return &APIHandler{
		baseURL:    httputils.TrimHTTPScheme(baseURL),
		authHeader: authHeader,
		authToken:  authToken,
		httpClient: httpClient,
		scheme:     scheme,
	}
}

func (a *APIHandler) getBaseURL() string {
	return a.baseURL
}

func (a *APIHandler) getAuthToken() string {
	return a.authToken
}

func (a *APIHandler) getAuthHeader() string {
	return a.authHeader
}

func (a *APIHandler) getHTTPClient() *http.Client {
	return a.httpClient
}

// SendEvent sends an event to Keptn.
func (a *APIHandler) SendEvent(ctx context.Context, event models.KeptnContextExtendedCE, opts APISendEventOptions) (*models.EventContext, *models.Error) {
	bodyStr, err := event.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	baseURL := a.getBaseURL()
	if strings.HasSuffix(baseURL, "/"+shipyardControllerBaseURL) {
		baseURL = strings.TrimSuffix(a.getBaseURL(), "/"+shipyardControllerBaseURL)
		baseURL += "/api"
	}

	return postWithEventContext(ctx, a.scheme+"://"+baseURL+v1EventPath, bodyStr, a)
}

// TriggerEvaluation triggers a new evaluation.
func (a *APIHandler) TriggerEvaluation(ctx context.Context, project, stage, service string, evaluation models.Evaluation, opts APITriggerEvaluationOptions) (*models.EventContext, *models.Error) {
	bodyStr, err := evaluation.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, a.scheme+"://"+a.getBaseURL()+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+service+"/evaluation", bodyStr, a)
}

// CreateProject creates a new project.
func (a *APIHandler) CreateProject(ctx context.Context, project models.CreateProject, opts APICreateProjectOptions) (string, *models.Error) {

	bodyStr, err := project.ToJSON()
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return post(ctx, a.scheme+"://"+a.getBaseURL()+v1ProjectPath, bodyStr, a)
}

// UpdateProject updates a project.
func (a *APIHandler) UpdateProject(ctx context.Context, project models.CreateProject, opts APIUpdateProjectOptions) (string, *models.Error) {
	bodyStr, err := project.ToJSON()
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return put(ctx, a.scheme+"://"+a.getBaseURL()+v1ProjectPath, bodyStr, a)
}

// DeleteProject deletes a project.
func (a *APIHandler) DeleteProject(ctx context.Context, project models.Project, opts APIDeleteProjectOptions) (*models.DeleteProjectResponse, *models.Error) {
	resp, err := delete(ctx, a.scheme+"://"+a.getBaseURL()+v1ProjectPath+"/"+project.ProjectName, a)
	if err != nil {
		return nil, err
	}

	deletePrjResponse := &models.DeleteProjectResponse{}
	if err2 := deletePrjResponse.FromJSON([]byte(resp)); err2 != nil {
		msg := "Could not decode DeleteProjectResponse: " + err2.Error()
		return nil, &models.Error{
			Message: &msg,
		}
	}
	return deletePrjResponse, nil
}

// CreateService creates a new service.
func (a *APIHandler) CreateService(ctx context.Context, project string, service models.CreateService, opts APICreateServiceOptions) (string, *models.Error) {
	bodyStr, err := service.ToJSON()
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return post(ctx, a.scheme+"://"+a.getBaseURL()+v1ProjectPath+"/"+project+pathToService, bodyStr, a)
}

// DeleteService deletes a service.
func (a *APIHandler) DeleteService(ctx context.Context, project, service string, opts APIDeleteServiceOptions) (*models.DeleteServiceResponse, *models.Error) {
	resp, err := delete(ctx, a.scheme+"://"+a.getBaseURL()+v1ProjectPath+"/"+project+pathToService+"/"+service, a)

	if err != nil {
		return nil, err
	}

	deleteSvcResponse := &models.DeleteServiceResponse{}
	if err2 := deleteSvcResponse.FromJSON([]byte(resp)); err2 != nil {
		msg := "Could not decode DeleteServiceResponse: " + err2.Error()
		return nil, &models.Error{
			Message: &msg,
		}
	}
	return deleteSvcResponse, nil
}

// GetMetadata retrieves Keptn metadata information.
func (a *APIHandler) GetMetadata(ctx context.Context, opts APIGetMetadataOptions) (*models.Metadata, *models.Error) {
	baseURL := a.getBaseURL()
	if strings.HasSuffix(baseURL, "/"+shipyardControllerBaseURL) {
		baseURL = strings.TrimSuffix(a.getBaseURL(), "/"+shipyardControllerBaseURL)
		baseURL += "/api"
	}

	body, mErr := getAndExpectSuccess(ctx, a.scheme+"://"+baseURL+v1MetadataPath, nil)
	if mErr != nil {
		return nil, mErr

	}

	respMetadata := &models.Metadata{}
	if err := respMetadata.FromJSON(body); err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	return respMetadata, nil
}
