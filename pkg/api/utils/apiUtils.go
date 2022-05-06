package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
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

// APIV2SendEventOptions are options for APIV2Interface.SendEvent().
type APIV2SendEventOptions struct{}

// APIV2TriggerEvaluationOptions are options for APIV2Interface.TriggerEvaluation().
type APIV2TriggerEvaluationOptions struct{}

// APIV2CreateProjectOptions are options for APIV2Interface.CreateProject().
type APIV2CreateProjectOptions struct{}

// APIV2UpdateProjectOptions are options for APIV2Interface.UpdateProject().
type APIV2UpdateProjectOptions struct{}

// APIV2DeleteProjectOptions are options for APIV2Interface.DeleteProject().
type APIV2DeleteProjectOptions struct{}

// APIV2CreateServiceOptions are options for APIV2Interface.CreateService().
type APIV2CreateServiceOptions struct{}

// APIV2DeleteServiceOptions are options for APIV2Interface.DeleteService().
type APIV2DeleteServiceOptions struct{}

// APIV2GetMetadataOptions are options for APIV2Interface.GetMetadata().
type APIV2GetMetadataOptions struct{}

type APIV2Interface interface {
	// SendEvent sends an event to Keptn.
	SendEvent(ctx context.Context, event models.KeptnContextExtendedCE, opts APIV2SendEventOptions) (*models.EventContext, *models.Error)

	// TriggerEvaluation triggers a new evaluation.
	TriggerEvaluation(ctx context.Context, project string, stage string, service string, evaluation models.Evaluation, opts APIV2TriggerEvaluationOptions) (*models.EventContext, *models.Error)

	// CreateProject creates a new project.
	CreateProject(ctx context.Context, project models.CreateProject, opts APIV2CreateProjectOptions) (string, *models.Error)

	// UpdateProject updates a project.
	UpdateProject(ctx context.Context, project models.CreateProject, opts APIV2UpdateProjectOptions) (string, *models.Error)

	// DeleteProject deletes a project.
	DeleteProject(ctx context.Context, project models.Project, opts APIV2DeleteProjectOptions) (*models.DeleteProjectResponse, *models.Error)

	// CreateService creates a new service.
	CreateService(ctx context.Context, project string, service models.CreateService, opts APIV2CreateServiceOptions) (string, *models.Error)

	// DeleteService deletes a service.
	DeleteService(ctx context.Context, project string, service string, opts APIV2DeleteServiceOptions) (*models.DeleteServiceResponse, *models.Error)

	// GetMetadata retrieves Keptn metadata information.
	GetMetadata(ctx context.Context, opts APIV2GetMetadataOptions) (*models.Metadata, *models.Error)
}

// APIHandler handles projects
type APIHandler struct {
	APIV2Handler
}

// APIV2Handler handles projects
type APIV2Handler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
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
	return &APIHandler{
		APIV2Handler: *createAuthenticatedAPIV2Handler(baseURL, authToken, authHeader, httpClient, scheme),
	}
}

func createAuthenticatedAPIV2Handler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *APIV2Handler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return &APIV2Handler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (a *APIV2Handler) getBaseURL() string {
	return a.BaseURL
}

func (a *APIV2Handler) getAuthToken() string {
	return a.AuthToken
}

func (a *APIV2Handler) getAuthHeader() string {
	return a.AuthHeader
}

func (a *APIV2Handler) getHTTPClient() *http.Client {
	return a.HTTPClient
}

// SendEvent sends an event to Keptn.
func (a *APIHandler) SendEvent(event models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
	return a.APIV2Handler.SendEvent(context.TODO(), event, APIV2SendEventOptions{})
}

// SendEvent sends an event to Keptn.
func (a *APIV2Handler) SendEvent(ctx context.Context, event models.KeptnContextExtendedCE, opts APIV2SendEventOptions) (*models.EventContext, *models.Error) {
	bodyStr, err := event.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	baseURL := a.getBaseURL()
	if strings.HasSuffix(baseURL, "/"+shipyardControllerBaseURL) {
		baseURL = strings.TrimSuffix(a.getBaseURL(), "/"+shipyardControllerBaseURL)
		baseURL += "/api"
	}

	return postWithEventContext(ctx, a.Scheme+"://"+baseURL+v1EventPath, bodyStr, a)
}

// TriggerEvaluation triggers a new evaluation.
func (a *APIHandler) TriggerEvaluation(project, stage, service string, evaluation models.Evaluation) (*models.EventContext, *models.Error) {
	return a.APIV2Handler.TriggerEvaluation(context.TODO(), project, stage, service, evaluation, APIV2TriggerEvaluationOptions{})
}

// TriggerEvaluation triggers a new evaluation.
func (a *APIV2Handler) TriggerEvaluation(ctx context.Context, project, stage, service string, evaluation models.Evaluation, opts APIV2TriggerEvaluationOptions) (*models.EventContext, *models.Error) {
	bodyStr, err := evaluation.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, a.Scheme+"://"+a.getBaseURL()+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+service+"/evaluation", bodyStr, a)
}

// CreateProject creates a new project.
func (a *APIHandler) CreateProject(project models.CreateProject) (string, *models.Error) {
	return a.APIV2Handler.CreateProject(context.TODO(), project, APIV2CreateProjectOptions{})
}

// CreateProject creates a new project.
func (a *APIV2Handler) CreateProject(ctx context.Context, project models.CreateProject, opts APIV2CreateProjectOptions) (string, *models.Error) {

	bodyStr, err := project.ToJSON()
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return post(ctx, a.Scheme+"://"+a.getBaseURL()+v1ProjectPath, bodyStr, a)
}

// UpdateProject updates a project.
func (a *APIHandler) UpdateProject(project models.CreateProject) (string, *models.Error) {
	return a.APIV2Handler.UpdateProject(context.TODO(), project, APIV2UpdateProjectOptions{})
}

// UpdateProject updates a project.
func (a *APIV2Handler) UpdateProject(ctx context.Context, project models.CreateProject, opts APIV2UpdateProjectOptions) (string, *models.Error) {
	bodyStr, err := project.ToJSON()
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return put(ctx, a.Scheme+"://"+a.getBaseURL()+v1ProjectPath, bodyStr, a)
}

// DeleteProject deletes a project.
func (a *APIHandler) DeleteProject(project models.Project) (*models.DeleteProjectResponse, *models.Error) {
	return a.APIV2Handler.DeleteProject(context.TODO(), project, APIV2DeleteProjectOptions{})
}

// DeleteProject deletes a project.
func (a *APIV2Handler) DeleteProject(ctx context.Context, project models.Project, opts APIV2DeleteProjectOptions) (*models.DeleteProjectResponse, *models.Error) {
	resp, err := delete(ctx, a.Scheme+"://"+a.getBaseURL()+v1ProjectPath+"/"+project.ProjectName, a)
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
func (a *APIHandler) CreateService(project string, service models.CreateService) (string, *models.Error) {
	return a.APIV2Handler.CreateService(context.TODO(), project, service, APIV2CreateServiceOptions{})
}

// CreateService creates a new service.
func (a *APIV2Handler) CreateService(ctx context.Context, project string, service models.CreateService, opts APIV2CreateServiceOptions) (string, *models.Error) {
	bodyStr, err := service.ToJSON()
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return post(ctx, a.Scheme+"://"+a.getBaseURL()+v1ProjectPath+"/"+project+pathToService, bodyStr, a)
}

// DeleteService deletes a service.
func (a *APIHandler) DeleteService(project, service string) (*models.DeleteServiceResponse, *models.Error) {
	return a.APIV2Handler.DeleteService(context.TODO(), project, service, APIV2DeleteServiceOptions{})
}

// DeleteService deletes a service.
func (a *APIV2Handler) DeleteService(ctx context.Context, project, service string, opts APIV2DeleteServiceOptions) (*models.DeleteServiceResponse, *models.Error) {
	resp, err := delete(ctx, a.Scheme+"://"+a.getBaseURL()+v1ProjectPath+"/"+project+pathToService+"/"+service, a)
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
func (a *APIHandler) GetMetadata() (*models.Metadata, *models.Error) {
	return a.APIV2Handler.GetMetadata(context.TODO(), APIV2GetMetadataOptions{})
}

// GetMetadata retrieves Keptn metadata information.
func (a *APIV2Handler) GetMetadata(ctx context.Context, opts APIV2GetMetadataOptions) (*models.Metadata, *models.Error) {
	baseURL := a.getBaseURL()
	if strings.HasSuffix(baseURL, "/"+shipyardControllerBaseURL) {
		baseURL = strings.TrimSuffix(a.getBaseURL(), "/"+shipyardControllerBaseURL)
		baseURL += "/api"
	}

	body, mErr := getAndExpectSuccess(ctx, a.Scheme+"://"+baseURL+v1MetadataPath, nil)
	if mErr != nil {
		return nil, mErr

	}

	respMetadata := &models.Metadata{}
	if err := respMetadata.FromJSON(body); err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	return respMetadata, nil
}
