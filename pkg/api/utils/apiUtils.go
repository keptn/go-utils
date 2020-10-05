package api

import (
	"encoding/json"
	"io/ioutil"
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
func NewAuthenticatedAPIHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string, clientCertPath, clientKeyPath, rootCertPath string) *APIHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	httpClient.Transport = getClientTransport(clientCertPath, clientKeyPath, rootCertPath)

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

// SendEvent sends an event to Keptn
func (e *APIHandler) SendEvent(event models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(event)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(e.Scheme+"://"+e.getBaseURL()+"/v1/event", bodyStr, e)
}

// TriggerEvaluation triggers a new evaluation
func (e *APIHandler) TriggerEvaluation(project, stage, service string, evaluation models.Evaluation) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(evaluation)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(e.Scheme+"://"+e.getBaseURL()+"/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/evaluation", bodyStr, e)
}

// GetEvent returns an event specified by keptnContext and eventType
//
// Deprecated: this function is deprecated and should be replaced with the GetEvents function
func (e *APIHandler) GetEvent(keptnContext string, eventType string) (*models.KeptnContextExtendedCE, *models.Error) {
	return getEvent(e.Scheme+"://"+e.getBaseURL()+"/v1/event?keptnContext="+keptnContext+"&type="+eventType+"&pageSize=10", e)
}

func getEvent(uri string, api APIService) (*models.KeptnContextExtendedCE, *models.Error) {

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
			var cloudEvent models.KeptnContextExtendedCE
			err = json.Unmarshal(body, &cloudEvent)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}

			return &cloudEvent, nil
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

// DeleteProject deletes a project
func (p *APIHandler) DeleteService(project, service string) (*models.EventContext, *models.Error) {
	return delete(p.Scheme+"://"+p.getBaseURL()+"/v1/project/"+project+"/service/"+service, p)
}

// GetMetadata retrieve keptn MetaData information
func (s *APIHandler) GetMetadata() (*models.Metadata, *models.Error) {
	//return get(s.Scheme+"://"+s.getBaseURL()+"/v1/metadata", nil, s)

	req, err := http.NewRequest("GET", s.Scheme+"://"+s.getBaseURL()+"/v1/metadata", nil)
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, s)

	resp, err := s.getHTTPClient().Do(req)
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
			var respMetadata models.Metadata
			err = json.Unmarshal(body, &respMetadata)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}

			return &respMetadata, nil
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
