package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const v1ProjectPath = "/v1/project"

type ProjectsV1Interface interface {
	// CreateProject creates a new project.
	CreateProject(project models.Project) (*models.EventContext, *models.Error)

	// DeleteProject deletes a project.
	DeleteProject(project models.Project) (*models.EventContext, *models.Error)

	// GetProject returns a project.
	GetProject(project models.Project) (*models.Project, *models.Error)

	// GetAllProjects returns all projects.
	GetAllProjects() ([]*models.Project, error)

	// UpdateConfigurationServiceProject updates a configuration service project.
	UpdateConfigurationServiceProject(project models.Project) (*models.EventContext, *models.Error)
}

// ProjectHandler handles projects
type ProjectHandler struct {
	projectHandler *v2.ProjectHandler
	BaseURL        string
	AuthToken      string
	AuthHeader     string
	HTTPClient     *http.Client
	Scheme         string
}

// NewProjectHandler returns a new ProjectHandler which sends all requests directly to the configuration-service
func NewProjectHandler(baseURL string) *ProjectHandler {
	return NewProjectHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewProjectHandlerWithHTTPClient returns a new ProjectHandler which sends all requests directly to the configuration-service using the specified http.Client
func NewProjectHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *ProjectHandler {
	return &ProjectHandler{
		BaseURL:        httputils.TrimHTTPScheme(baseURL),
		HTTPClient:     httpClient,
		Scheme:         "http",
		projectHandler: v2.NewProjectHandlerWithHTTPClient(baseURL, httpClient),
	}
}

// NewAuthenticatedProjectHandler returns a new ProjectHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
// Deprecated: use APISet instead
func NewAuthenticatedProjectHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ProjectHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedProjectHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedProjectHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ProjectHandler {
	v2ProjectHandler := v2.NewAuthenticatedProjectHandler(baseURL, authToken, authHeader, httpClient, scheme)

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return &ProjectHandler{
		BaseURL:        httputils.TrimHTTPScheme(baseURL),
		AuthHeader:     authHeader,
		AuthToken:      authToken,
		HTTPClient:     httpClient,
		Scheme:         scheme,
		projectHandler: v2ProjectHandler,
	}
}

func (p *ProjectHandler) getBaseURL() string {
	return p.BaseURL
}

func (p *ProjectHandler) getAuthToken() string {
	return p.AuthToken
}

func (p *ProjectHandler) getAuthHeader() string {
	return p.AuthHeader
}

func (p *ProjectHandler) getHTTPClient() *http.Client {
	return p.HTTPClient
}

// CreateProject creates a new project.
func (p *ProjectHandler) CreateProject(project models.Project) (*models.EventContext, *models.Error) {
	return p.projectHandler.CreateProject(context.TODO(), project, v2.ProjectsCreateProjectOptions{})
}

// DeleteProject deletes a project.
func (p *ProjectHandler) DeleteProject(project models.Project) (*models.EventContext, *models.Error) {
	return p.projectHandler.DeleteProject(context.TODO(), project, v2.ProjectsDeleteProjectOptions{})
}

// GetProject returns a project.
func (p *ProjectHandler) GetProject(project models.Project) (*models.Project, *models.Error) {
	return p.projectHandler.GetProject(context.TODO(), project, v2.ProjectsGetProjectOptions{})
}

// GetAllProjects returns all projects.
func (p *ProjectHandler) GetAllProjects() ([]*models.Project, error) {
	return p.projectHandler.GetAllProjects(context.TODO(), v2.ProjectsGetAllProjectsOptions{})
}

// UpdateConfigurationServiceProject updates a configuration service project.
func (p *ProjectHandler) UpdateConfigurationServiceProject(project models.Project) (*models.EventContext, *models.Error) {
	return p.projectHandler.UpdateConfigurationServiceProject(context.TODO(), project, v2.ProjectsUpdateConfigurationServiceProjectOptions{})
}
