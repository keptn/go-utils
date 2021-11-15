package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/common/httputils"

	"github.com/keptn/go-utils/pkg/api/models"
)

const v1ProjectPath = "/v1/project"

// ProjectHandler handles projects
type ProjectHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewProjectHandler returns a new ProjectHandler which sends all requests directly to the configuration-service
func NewProjectHandler(baseURL string) *ProjectHandler {
	baseURL = httputils.TrimHTTPScheme(baseURL)
	return &ProjectHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getInstrumentedClientTransport()},
		Scheme:     "http",
	}
}

// NewAuthenticatedProjectHandler returns a new ProjectHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
func NewAuthenticatedProjectHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ProjectHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getInstrumentedClientTransport()

	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	baseURL = strings.TrimRight(baseURL, "/")

	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return &ProjectHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
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

// CreateProject creates a new project
//
// Deprecated: Use CreateProjectWithContext instead
func (p *ProjectHandler) CreateProject(project models.Project) (*models.EventContext, *models.Error) {
	return p.CreateProjectWithContext(context.Background(), project)
}

// CreateProjectWithContext creates a new project
func (p *ProjectHandler) CreateProjectWithContext(ctx context.Context, project models.Project) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(project)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, p.Scheme+"://"+p.getBaseURL()+v1ProjectPath, bodyStr, p)
}

// DeleteProject deletes a project
//
// Deprecated: Use DeleteProjectWithContext instead
func (p *ProjectHandler) DeleteProject(project models.Project) (*models.EventContext, *models.Error) {
	return p.DeleteProjectWithContext(context.Background(), project)
}

// DeleteProjectWithContext deletes a project
func (p *ProjectHandler) DeleteProjectWithContext(ctx context.Context, project models.Project) (*models.EventContext, *models.Error) {
	return deleteWithEventContext(ctx, p.Scheme+"://"+p.getBaseURL()+v1ProjectPath+"/"+project.ProjectName, p)
}

// GetProject returns a project
//
// Deprecated: Use GetProjectWithContext instead
func (p *ProjectHandler) GetProject(project models.Project) (*models.Project, *models.Error) {
	return p.GetProjectWithContext(context.Background(), project)
}

// GetProjectWithContext returns a project
func (p *ProjectHandler) GetProjectWithContext(ctx context.Context, project models.Project) (*models.Project, *models.Error) {
	return getProject(ctx, p.Scheme+"://"+p.getBaseURL()+v1ProjectPath+"/"+project.ProjectName, p)
}

// GetAllProjects returns all projects
//
// Deprecated: Use GetAllProjectsWithContext instead
func (p *ProjectHandler) GetAllProjects() ([]*models.Project, error) {
	return p.GetAllProjectsWithContext(context.Background())
}

// GetAllProjectsWithContext returns all projects
func (p *ProjectHandler) GetAllProjectsWithContext(ctx context.Context) ([]*models.Project, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	projects := []*models.Project{}

	nextPageKey := ""

	for {
		url, err := url.Parse(p.Scheme + "://" + p.getBaseURL() + v1ProjectPath)
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, p)

		resp, err := p.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 200 {

			var received models.Projects
			err = json.Unmarshal(body, &received)
			if err != nil {
				return nil, err
			}
			projects = append(projects, received.Projects...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey
		} else {
			var respErr models.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(*respErr.Message)
		}
	}

	return projects, nil
}

func getProject(ctx context.Context, uri string, api APIService) (*models.Project, *models.Error) {
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, api)

	resp, err := api.getHTTPClient().Do(req)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		if len(body) > 0 {
			var respProject models.Project
			err = json.Unmarshal(body, &respProject)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}

			return &respProject, nil
		}

		return nil, nil
	}

	var respErr models.Error
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	return nil, &respErr
}

// UpdateConfigurationServiceProject updates the project configuration
//
// Deprecated: Use UpdateConfigurationServiceProjectWithContext instead
func (p *ProjectHandler) UpdateConfigurationServiceProject(project models.Project) (*models.EventContext, *models.Error) {
	return p.UpdateConfigurationServiceProjectWithContext(context.Background(), project)
}

// UpdateConfigurationServiceProjectWithContext updates the project configuration
func (p *ProjectHandler) UpdateConfigurationServiceProjectWithContext(ctx context.Context, project models.Project) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(project)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return putWithEventContext(ctx, p.Scheme+"://"+p.getBaseURL()+v1ProjectPath+"/"+project.ProjectName, bodyStr, p)
}
