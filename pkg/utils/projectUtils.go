package utils

import (
	"encoding/json"

	"github.com/keptn/go-utils/pkg/models"
)

// ProjectHandler handles projects
type ProjectHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
}

// NewProjectHandler returns a new ProjectHandler
func NewProjectHandler(baseURL string) *ProjectHandler {
	return &ProjectHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
	}
}

// NewAuthenticatedProjectHandler returns a new ProjectHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedProjectHandler(baseURL string, authToken string, authHeader string) *ProjectHandler {
	return &ProjectHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
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

// CreateProject creates a new project
func (p *ProjectHandler) CreateProject(project models.Project) (*models.Error, error) {

	bodyStr, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}
	return post("http://"+p.getBaseURL()+"/v1/project", bodyStr, p)
}
