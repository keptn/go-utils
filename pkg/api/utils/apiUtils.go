package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

// APIHandler handles projects
type APIHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewAuthenticatedAPIHandler returns a new APIHandler that authenticates at the api-service endpoint via the provided token
func NewAuthenticatedAPIHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *APIHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &APIHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (p *APIHandler) getBaseURL() string {
	return p.BaseURL
}

func (p *APIHandler) getAuthToken() string {
	return p.AuthToken
}

func (p *APIHandler) getAuthHeader() string {
	return p.AuthHeader
}

func (p *APIHandler) getHTTPClient() *http.Client {
	return p.HTTPClient
}

// CreateProject creates a new project
func (p *APIHandler) CreateProject(project models.CreateProject) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(project)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(p.Scheme+"://"+p.getBaseURL()+"/v1/project", bodyStr, p)
}

// DeleteProject deletes a project
func (p *APIHandler) DeleteProject(project models.Project) (*models.EventContext, *models.Error) {
	return delete(p.Scheme+"://"+p.getBaseURL()+"/v1/project/"+project.ProjectName, p)
}

// CreateService creates a new service
func (s *APIHandler) CreateService(project string, service models.CreateService) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(service)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(s.Scheme+"://"+s.getBaseURL()+"/v1/project/"+project+"/service", bodyStr, s)
}
