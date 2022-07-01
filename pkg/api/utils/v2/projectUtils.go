package v2

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/common/httputils"

	"github.com/keptn/go-utils/pkg/api/models"
)

const v1ProjectPath = "/v1/project"

// ProjectsCreateProjectOptions are options for ProjectsInterface.CreateProject().
type ProjectsCreateProjectOptions struct{}

// ProjectsDeleteProjectOptions are options for ProjectsInterface.DeleteProject().
type ProjectsDeleteProjectOptions struct{}

// ProjectsGetProjectOptions are options for ProjectsInterface.GetProject().
type ProjectsGetProjectOptions struct{}

// ProjectsGetAllProjectsOptions are options for ProjectsInterface.GetAllProjects().
type ProjectsGetAllProjectsOptions struct{}

// ProjectsUpdateConfigurationServiceProjectOptions are options for ProjectsInterface.UpdateConfigurationServiceProject().
type ProjectsUpdateConfigurationServiceProjectOptions struct{}

type ProjectsInterface interface {
	// CreateProject creates a new project.
	CreateProject(ctx context.Context, project models.Project, opts ProjectsCreateProjectOptions) (*models.EventContext, *models.Error)

	// DeleteProject deletes a project.
	DeleteProject(ctx context.Context, project models.Project, opts ProjectsDeleteProjectOptions) (*models.EventContext, *models.Error)

	// GetProject returns a project.
	GetProject(ctx context.Context, project models.Project, opts ProjectsGetProjectOptions) (*models.Project, *models.Error)

	// GetAllProjects returns all projects.
	GetAllProjects(ctx context.Context, opts ProjectsGetAllProjectsOptions) ([]*models.Project, error)

	// UpdateConfigurationServiceProject updates a configuration service project.
	UpdateConfigurationServiceProject(ctx context.Context, project models.Project, opts ProjectsUpdateConfigurationServiceProjectOptions) (*models.EventContext, *models.Error)
}

// ProjectHandler handles projects
type ProjectHandler struct {
	baseURL    string
	authToken  string
	authHeader string
	httpClient *http.Client
	scheme     string
}

// NewProjectHandler returns a new ProjectHandler which sends all requests directly to the resource-service
func NewProjectHandler(baseURL string) *ProjectHandler {
	return NewProjectHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewProjectHandlerWithHTTPClient returns a new ProjectHandler which sends all requests directly to the resource-service using the specified http.Client
func NewProjectHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *ProjectHandler {
	return createProjectHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedProjectHandler returns a new ProjectHandler that authenticates at the api via the provided token
// and sends all requests directly to the resource-service
func NewAuthenticatedProjectHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ProjectHandler {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return createProjectHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createProjectHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ProjectHandler {
	return &ProjectHandler{
		baseURL:    httputils.TrimHTTPScheme(baseURL),
		authHeader: authHeader,
		authToken:  authToken,
		httpClient: httpClient,
		scheme:     scheme,
	}
}

func (p *ProjectHandler) getBaseURL() string {
	return p.baseURL
}

func (p *ProjectHandler) getAuthToken() string {
	return p.authToken
}

func (p *ProjectHandler) getAuthHeader() string {
	return p.authHeader
}

func (p *ProjectHandler) getHTTPClient() *http.Client {
	return p.httpClient
}

// CreateProject creates a new project.
func (p *ProjectHandler) CreateProject(ctx context.Context, project models.Project, opts ProjectsCreateProjectOptions) (*models.EventContext, *models.Error) {
	bodyStr, err := project.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, p.scheme+"://"+p.getBaseURL()+v1ProjectPath, bodyStr, p)
}

// DeleteProject deletes a project.
func (p *ProjectHandler) DeleteProject(ctx context.Context, project models.Project, opts ProjectsDeleteProjectOptions) (*models.EventContext, *models.Error) {
	return deleteWithEventContext(ctx, p.scheme+"://"+p.getBaseURL()+v1ProjectPath+"/"+project.ProjectName, p)
}

// GetProject returns a project.
func (p *ProjectHandler) GetProject(ctx context.Context, project models.Project, opts ProjectsGetProjectOptions) (*models.Project, *models.Error) {
	body, mErr := getAndExpectSuccess(ctx, p.scheme+"://"+p.getBaseURL()+v1ProjectPath+"/"+project.ProjectName, p)
	if mErr != nil {
		return nil, mErr
	}

	respProject := &models.Project{}
	if err := respProject.FromJSON(body); err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	return respProject, nil
}

// GetAllProjects returns all projects.
func (p *ProjectHandler) GetAllProjects(ctx context.Context, opts ProjectsGetAllProjectsOptions) ([]*models.Project, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	projects := []*models.Project{}

	nextPageKey := ""

	for {
		url, err := url.Parse(p.scheme + "://" + p.getBaseURL() + v1ProjectPath)
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}

		body, mErr := getAndExpectOK(ctx, url.String(), p)
		if mErr != nil {
			return nil, mErr.ToError()
		}

		received := &models.Projects{}
		if err = received.FromJSON(body); err != nil {
			return nil, err
		}
		projects = append(projects, received.Projects...)

		if received.NextPageKey == "" || received.NextPageKey == "0" {
			break
		}
		nextPageKey = received.NextPageKey
	}

	return projects, nil
}

// UpdateConfigurationServiceProject updates a configuration service project.
func (p *ProjectHandler) UpdateConfigurationServiceProject(ctx context.Context, project models.Project, opts ProjectsUpdateConfigurationServiceProjectOptions) (*models.EventContext, *models.Error) {
	bodyStr, err := project.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return putWithEventContext(ctx, p.scheme+"://"+p.getBaseURL()+v1ProjectPath+"/"+project.ProjectName, bodyStr, p)
}
