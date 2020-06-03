package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
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
	scheme := "http"
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
		scheme = "http"
	}
	return &ProjectHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     scheme,
	}
}

const configurationServiceBaseUrl = "configuration-service"

// NewAuthenticatedProjectHandler returns a new ProjectHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedProjectHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ProjectHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, configurationServiceBaseUrl) {
		baseURL += "/" + configurationServiceBaseUrl
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
func (p *ProjectHandler) CreateProject(project models.CreateProject) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(project)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(p.Scheme+"://"+p.getBaseURL()+"/v1/project", bodyStr, p)
}

func (p *ProjectHandler) CreateConfigurationServiceProject(project models.Project) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(project)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(p.Scheme+"://"+p.getBaseURL()+"/v1/project", bodyStr, p)
}

// DeleteProject deletes a project
func (p *ProjectHandler) DeleteProject(project models.Project) (*models.EventContext, *models.Error) {
	return delete(p.Scheme+"://"+p.getBaseURL()+"/v1/project/"+project.ProjectName, p)
}

// GetProject returns a project
func (p *ProjectHandler) GetProject(project models.Project) (*models.Project, *models.Error) {
	return getProject(p.Scheme+"://"+p.getBaseURL()+"/v1/project/"+project.ProjectName, p)
}

// GetProjects returns a project
func (p *ProjectHandler) GetAllProjects() ([]*models.Project, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	projects := []*models.Project{}

	nextPageKey := ""

	for {
		url, err := url.Parse(p.Scheme + "://" + p.getBaseURL() + "/v1/project/")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String(), nil)
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
			return nil, errors.New("Response Error Code: " + string(respErr.Code) + " Message: " + *respErr.Message)
		}
	}

	return projects, nil
}

func getProject(uri string, api APIService) (*models.Project, *models.Error) {

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
