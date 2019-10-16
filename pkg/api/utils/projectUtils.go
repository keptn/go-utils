package utils

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
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
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
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
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
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
func (p *ProjectHandler) CreateProject(project models.Project) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(project)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(p.Scheme+"://"+p.getBaseURL()+"/v1/project", bodyStr, p)
}

// DeleteProject deletes a project
func (p *ProjectHandler) DeleteProject(project models.Project) (*models.EventContext, *models.Error) {
	return delete(p.Scheme+"://"+p.getBaseURL()+"/v1/project/"+*project.Name, p)
}

// GetProject returns a project
func (p *ProjectHandler) GetProject(project models.Project) (*models.Project, *models.Error) {
	return getProject(p.Scheme+"://"+p.getBaseURL()+"/v1/project/"+*project.Name, p)
}

func getProject(uri string, api APIService) (*models.Project, *models.Error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", uri, nil)
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
