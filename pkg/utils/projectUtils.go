package utils

import (
	"encoding/json"
	"net/http"

	"github.com/keptn/go-utils/pkg/models"
)

// ProjectHandler handles projects
type ProjectHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewProjectHandler returns a new ProjectHandler
func NewProjectHandler(baseURL string) *ProjectHandler {
	return &ProjectHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{},
		Scheme:     "http",
	}
}

// NewAuthenticatedProjectHandler returns a new ProjectHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedProjectHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ProjectHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
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
func (p *ProjectHandler) CreateProject(project models.Project) (*models.Error, error) {

	bodyStr, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}
	return post(p.Scheme+p.getBaseURL()+"/v1/project", bodyStr, p)
}
