package api

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

// ResourceHandler handles resources
type ResourceHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

type resourceRequest struct {
	Resources []*models.Resource `json:"resources"`
}

const configurationServiceBaseUrl = "configuration-service"

var ResourceNotFoundError = errors.New("Resource not found")

// NewResourceHandler returns a new ResourceHandler which sends all requests directly to the configuration-service
func NewResourceHandler(baseURL string) *ResourceHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &ResourceHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))},
		Scheme:     "http",
	}
}

// NewAuthenticatedResourceHandler returns a new ResourceHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
// Deprecated: use APISet instead
func NewAuthenticatedResourceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ResourceHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedResourceHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedResourceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ResourceHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, configurationServiceBaseUrl) {
		baseURL += "/" + configurationServiceBaseUrl
	}
	return &ResourceHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (r *ResourceHandler) getBaseURL() string {
	return r.BaseURL
}

func (r *ResourceHandler) getAuthToken() string {
	return r.AuthToken
}

func (r *ResourceHandler) getAuthHeader() string {
	return r.AuthHeader
}

func (r *ResourceHandler) getHTTPClient() *http.Client {
	return r.HTTPClient
}

// CreateResources creates a resource for the specified entity
func (r *ResourceHandler) CreateResources(project string, stage string, service string, resources []*models.Resource) (*models.EventContext, *models.Error) {

	copiedResources := make([]*models.Resource, len(resources), len(resources))
	for i, val := range resources {
		resourceContent := b64.StdEncoding.EncodeToString([]byte(val.ResourceContent))
		copiedResources[i] = &models.Resource{ResourceURI: val.ResourceURI, ResourceContent: resourceContent}
	}

	resReq := &resourceRequest{
		Resources: copiedResources,
	}
	requestStr, err := json.Marshal(resReq)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	if project != "" && stage != "" && service != "" {
		return postWithEventContext(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/resource", requestStr, r)
	} else if project != "" && stage != "" && service == "" {
		return postWithEventContext(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", requestStr, r)
	} else {
		return postWithEventContext(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource", requestStr, r)
	}
}

// CreateProjectResources creates multiple project resources
func (r *ResourceHandler) CreateProjectResources(project string, resources []*models.Resource) (string, error) {
	return r.createResources(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource", resources)
}

// GetProjectResource retrieves a project resource from the configuration service
func (r *ResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return r.getResource(r.Scheme + "://" + r.BaseURL + "/v1/project/" + project + "/resource/" + url.QueryEscape(resourceURI))
}

// UpdateProjectResource updates a project resource
func (r *ResourceHandler) UpdateProjectResource(project string, resource *models.Resource) (string, error) {
	return r.updateResource(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// DeleteProjectResource deletes a project resource
func (r *ResourceHandler) DeleteProjectResource(project string, resourceURI string) error {
	return r.deleteResource(r.Scheme + "://" + r.BaseURL + "/v1/project/" + project + "/resource/" + url.QueryEscape(resourceURI))
}

// UpdateProjectResources updates multiple project resources
func (r *ResourceHandler) UpdateProjectResources(project string, resources []*models.Resource) (string, error) {
	return r.updateResources(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource", resources)
}

// CreateStageResources creates a stage resource
func (r *ResourceHandler) CreateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	return r.createResources(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", resources)
}

// GetStageResource retrieves a stage resource from the configuration service
func (r *ResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return r.getResource(r.Scheme + "://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/resource/" + url.QueryEscape(resourceURI))
}

// UpdateStageResource updates a stage resource
func (r *ResourceHandler) UpdateStageResource(project string, stage string, resource *models.Resource) (string, error) {
	return r.updateResource(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateStageResources updates multiple stage resources
func (r *ResourceHandler) UpdateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	return r.updateResources(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", resources)
}

// DeleteStageResource deletes a stage resource
func (r *ResourceHandler) DeleteStageResource(project string, stage string, resourceURI string) error {
	return r.deleteResource(r.Scheme + "://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/resource/" + url.QueryEscape(resourceURI))
}

// CreateServiceResources creates a service resource
func (r *ResourceHandler) CreateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.createResources(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/resource", resources)
}

// GetServiceResource retrieves a service resource from the configuration service
func (r *ResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return r.getResource(r.Scheme + "://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/service/" + url.QueryEscape(service) + "/resource/" + url.QueryEscape(resourceURI))
}

// UpdateServiceResource updates a service resource
func (r *ResourceHandler) UpdateServiceResource(project string, stage string, service string, resource *models.Resource) (string, error) {
	return r.updateResource(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+url.QueryEscape(service)+"/resource/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateServiceResources updates multiple service resources
func (r *ResourceHandler) UpdateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.updateResources(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+url.QueryEscape(service)+"/resource", resources)
}

// DeleteServiceResource deletes a service resource
func (r *ResourceHandler) DeleteServiceResource(project string, stage string, service string, resourceURI string) error {
	return r.deleteResource(r.Scheme + "://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/service/" + url.QueryEscape(service) + "/resource/" + url.QueryEscape(resourceURI))
}

func (r *ResourceHandler) createResources(uri string, resources []*models.Resource) (string, error) {
	return r.writeResources(uri, "POST", resources)
}

func (r *ResourceHandler) updateResources(uri string, resources []*models.Resource) (string, error) {
	return r.writeResources(uri, "PUT", resources)
}

func (r *ResourceHandler) writeResources(uri string, method string, resources []*models.Resource) (string, error) {

	copiedResources := make([]*models.Resource, len(resources), len(resources))
	for i, val := range resources {
		copiedResources[i] = &models.Resource{ResourceURI: val.ResourceURI, ResourceContent: b64.StdEncoding.EncodeToString([]byte(val.ResourceContent))}
	}
	resReq := &resourceRequest{
		Resources: copiedResources,
	}

	resourceStr, err := json.Marshal(resReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(resourceStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, r)

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var version models.Version
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", errors.New(string(body))
	}

	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", err
	}

	return version.Version, nil
}

func (r *ResourceHandler) updateResource(uri string, resource *models.Resource) (string, error) {
	return r.writeResource(uri, "PUT", resource)
}

func (r *ResourceHandler) writeResource(uri string, method string, resource *models.Resource) (string, error) {

	copiedResource := &models.Resource{ResourceURI: resource.ResourceURI, ResourceContent: b64.StdEncoding.EncodeToString([]byte(resource.ResourceContent))}

	resourceStr, err := json.Marshal(copiedResource)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(resourceStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, r)

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", errors.New(string(body))
	}

	var version models.Version
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", err
	}

	return version.Version, nil
}

func (r *ResourceHandler) getResource(uri string) (*models.Resource, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, r)

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		// need to handle this case differently (e.g. https://github.com/keptn/keptn/issues/1480)
		return nil, ResourceNotFoundError
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, errors.New(string(body))
	}

	var resource models.Resource
	err = json.Unmarshal(body, &resource)
	if err != nil {
		return nil, err
	}

	// decode resource content
	decodedStr, err := b64.StdEncoding.DecodeString(resource.ResourceContent)
	if err != nil {
		return nil, err
	}
	resource.ResourceContent = string(decodedStr)

	return &resource, nil
}

func (r *ResourceHandler) deleteResource(uri string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, r)

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetAllStageResources returns a list of all resources.
func (r *ResourceHandler) GetAllStageResources(project string, stage string) ([]*models.Resource, error) {
	url, err := url.Parse(r.Scheme + "://" + r.getBaseURL() + "/v1/project/" + project + "/stage/" + stage + "/resource")
	if err != nil {
		return nil, err
	}
	return r.getAllResources(url)
}

// GetAllServiceResources returns a list of all resources.
func (r *ResourceHandler) GetAllServiceResources(project string, stage string, service string) ([]*models.Resource, error) {
	url, err := url.Parse(r.Scheme + "://" + r.getBaseURL() + "/v1/project/" + project + "/stage/" + stage +
		"/service/" + service + "/resource/")
	if err != nil {
		return nil, err
	}
	return r.getAllResources(url)
}

func (r *ResourceHandler) getAllResources(u *url.URL) ([]*models.Resource, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resources := []*models.Resource{}

	nextPageKey := ""

	for {
		if nextPageKey != "" {
			q := u.Query()
			q.Set("nextPageKey", nextPageKey)
			u.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, r)

		resp, err := r.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 200 {
			var received models.Resources
			err = json.Unmarshal(body, &received)
			if err != nil {
				return nil, err
			}
			resources = append(resources, received.Resources...)

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

	return resources, nil
}
