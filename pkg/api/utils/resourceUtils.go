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

const pathToResource = "/resource"
const pathToService = "/service"
const pathToStage = "/stage"
const pathToProject = v1ProjectPath + "/"

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

type option func(url string) string

func appendQuery(queryParams url.Values) option {
	return func(buildURI string) string {
		if queryParams != nil {
			buildURI = buildURI + "?" + queryParams.Encode()
		}
		return buildURI
	}
}

func (r *ResourceHandler) applyOptions(buildURI string, options []option) string {
	for _, option := range options {
		buildURI = option(buildURI)
	}
	return buildURI
}

const configurationServiceBaseUrl = "configuration-service"

var ResourceNotFoundError = errors.New("Resource not found")

// ToJSON converts object to JSON string
func (r *resourceRequest) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

func (r *resourceRequest) FromJSON(b []byte) error {
	var res resourceRequest
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*r = res
	return nil
}

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
	requestStr, err := resReq.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	if project != "" && stage != "" && service != "" {
		return postWithEventContext(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToStage+"/"+stage+pathToService+"/"+service+pathToResource, requestStr, r)
	} else if project != "" && stage != "" && service == "" {
		return postWithEventContext(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToStage+"/"+stage+pathToResource, requestStr, r)
	} else {
		return postWithEventContext(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToResource, requestStr, r)
	}
}

// CreateProjectResources creates multiple project resources
func (r *ResourceHandler) CreateProjectResources(project string, resources []*models.Resource) (string, error) {
	return r.createResources(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToResource, resources)
}

// GetProjectResource retrieves a project resource from the configuration service
func (r *ResourceHandler) GetProjectResource(project string, resourceURI string, options ...option) (*models.Resource, error) {
	buildURI := r.Scheme + "://" + r.BaseURL + pathToProject + project + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.getResource(r.applyOptions(buildURI, options))
}

// UpdateProjectResource updates a project resource
func (r *ResourceHandler) UpdateProjectResource(project string, resource *models.Resource) (string, error) {
	return r.updateResource(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// DeleteProjectResource deletes a project resource
func (r *ResourceHandler) DeleteProjectResource(project string, resourceURI string) error {
	return r.deleteResource(r.Scheme + "://" + r.BaseURL + pathToProject + project + pathToResource + "/" + url.QueryEscape(resourceURI))
}

// UpdateProjectResources updates multiple project resources
func (r *ResourceHandler) UpdateProjectResources(project string, resources []*models.Resource) (string, error) {
	return r.updateResources(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToResource, resources)
}

// CreateStageResources creates a stage resource
func (r *ResourceHandler) CreateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	return r.createResources(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToStage+"/"+stage+pathToResource, resources)
}

// GetStageResource retrieves a stage resource from the configuration service
func (r *ResourceHandler) GetStageResource(project string, stage string, resourceURI string, options ...option) (*models.Resource, error) {
	buildURI := r.Scheme + "://" + r.BaseURL + pathToProject + project + pathToStage + "/" + stage + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.getResource(r.applyOptions(buildURI, options))
}

// UpdateStageResource updates a stage resource
func (r *ResourceHandler) UpdateStageResource(project string, stage string, resource *models.Resource) (string, error) {
	return r.updateResource(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToStage+"/"+stage+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateStageResources updates multiple stage resources
func (r *ResourceHandler) UpdateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	return r.updateResources(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToStage+"/"+stage+pathToResource, resources)
}

// DeleteStageResource deletes a stage resource
func (r *ResourceHandler) DeleteStageResource(project string, stage string, resourceURI string) error {
	return r.deleteResource(r.Scheme + "://" + r.BaseURL + pathToProject + project + pathToStage + "/" + stage + pathToResource + "/" + url.QueryEscape(resourceURI))
}

// CreateServiceResources creates a service resource
func (r *ResourceHandler) CreateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.createResources(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToStage+"/"+stage+pathToService+"/"+service+pathToResource, resources)
}

// GetServiceResource retrieves a service resource from the configuration service
func (r *ResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string, options ...option) (*models.Resource, error) {
	buildURI := r.Scheme + "://" + r.BaseURL + pathToProject + project + pathToStage + "/" + stage + pathToService + "/" + url.QueryEscape(service) + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.getResource(r.applyOptions(buildURI, options))
}

// UpdateServiceResource updates a service resource
func (r *ResourceHandler) UpdateServiceResource(project string, stage string, service string, resource *models.Resource) (string, error) {
	return r.updateResource(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToStage+"/"+stage+pathToService+"/"+url.QueryEscape(service)+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateServiceResources updates multiple service resources
func (r *ResourceHandler) UpdateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.updateResources(r.Scheme+"://"+r.BaseURL+pathToProject+project+pathToStage+"/"+stage+pathToService+"/"+url.QueryEscape(service)+pathToResource, resources)
}

// DeleteServiceResource deletes a service resource
func (r *ResourceHandler) DeleteServiceResource(project string, stage string, service string, resourceURI string) error {
	return r.deleteResource(r.Scheme + "://" + r.BaseURL + pathToProject + project + pathToStage + "/" + stage + pathToService + "/" + url.QueryEscape(service) + pathToResource + "/" + url.QueryEscape(resourceURI))
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

	resourceStr, err := resReq.ToJSON()
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

	version := &models.Version{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", errors.New(string(body))
	}

	if err = version.FromJSON(body); err != nil {
		return "", err
	}

	return version.Version, nil
}

func (r *ResourceHandler) updateResource(uri string, resource *models.Resource) (string, error) {
	return r.writeResource(uri, "PUT", resource)
}

func (r *ResourceHandler) writeResource(uri string, method string, resource *models.Resource) (string, error) {

	copiedResource := &models.Resource{ResourceURI: resource.ResourceURI, ResourceContent: b64.StdEncoding.EncodeToString([]byte(resource.ResourceContent))}

	resourceStr, err := copiedResource.ToJSON()
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

	version := &models.Version{}
	if err = version.FromJSON(body); err != nil {
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

	resource := &models.Resource{}
	if err = resource.FromJSON(body); err != nil {
		return nil, err
	}

	// decode resource content
	decodedStr, err := b64.StdEncoding.DecodeString(resource.ResourceContent)
	if err != nil {
		return nil, err
	}
	resource.ResourceContent = string(decodedStr)

	return resource, nil
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
	myURL, err := url.Parse(r.Scheme + "://" + r.getBaseURL() + pathToProject + project + pathToStage + "/" + stage + pathToResource)
	if err != nil {
		return nil, err
	}
	return r.getAllResources(myURL)
}

// GetAllServiceResources returns a list of all resources.
func (r *ResourceHandler) GetAllServiceResources(project string, stage string, service string) ([]*models.Resource, error) {
	myURL, err := url.Parse(r.Scheme + "://" + r.getBaseURL() + pathToProject + project + pathToStage + "/" + stage +
		pathToService + "/" + service + pathToResource + "/")
	if err != nil {
		return nil, err
	}
	return r.getAllResources(myURL)
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
			received := &models.Resources{}
			if err = received.FromJSON(body); err != nil {
				return nil, err
			}
			resources = append(resources, received.Resources...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey

		} else {
			respErr := &models.Error{}
			if err = respErr.FromJSON(body); err != nil {
				return nil, err
			}
			return nil, errors.New(*respErr.Message)
		}
	}

	return resources, nil
}
