package api

import (
	"bytes"
	"context"
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
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     "http",
	}
}

// NewAuthenticatedResourceHandler returns a new ResourceHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
func NewAuthenticatedResourceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ResourceHandler {
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
//
// Deprecated: Use CreateResourcesWithContext instead
func (r *ResourceHandler) CreateResources(project string, stage string, service string, resources []*models.Resource) (*models.EventContext, *models.Error) {
	return r.CreateResourcesWithContext(context.Background(), project, stage, service, resources)
}

// CreateResourcesWithContext creates a resource for the specified entity
func (r *ResourceHandler) CreateResourcesWithContext(
	ctx context.Context, project string, stage string, service string, resources []*models.Resource) (*models.EventContext, *models.Error) {

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
		return postWithEventContext(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/resource", requestStr, r)
	} else if project != "" && stage != "" && service == "" {
		return postWithEventContext(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", requestStr, r)
	} else {
		return postWithEventContext(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource", requestStr, r)
	}
}

// CreateProjectResources creates multiple project resources
//
// Deprecated: Use CreateProjectResourcesWithContext instead
func (r *ResourceHandler) CreateProjectResources(project string, resources []*models.Resource) (string, error) {
	return r.CreateProjectResourcesWithContext(context.Background(), project, resources)
}

// CreateProjectResourcesWithContext creates multiple project resources
func (r *ResourceHandler) CreateProjectResourcesWithContext(ctx context.Context, project string, resources []*models.Resource) (string, error) {
	return r.createResources(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource", resources)
}

// GetProjectResource retrieves a project resource from the configuration service
//
// Deprecated: Use GetProjectResourceWithContext instead
func (r *ResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return r.GetProjectResourceWithContext(context.Background(), project, resourceURI)
}

// GetProjectResourceWithContext retrieves a project resource from the configuration service
func (r *ResourceHandler) GetProjectResourceWithContext(ctx context.Context, project string, resourceURI string) (*models.Resource, error) {
	return r.getResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource/"+url.QueryEscape(resourceURI))
}

// UpdateProjectResource updates a project resource
//
// Deprecated: Use UpdateProjectResourceWithContext instead
func (r *ResourceHandler) UpdateProjectResource(project string, resource *models.Resource) (string, error) {
	return r.UpdateProjectResourceWithContext(context.Background(), project, resource)
}

// UpdateProjectResourceWithContext updates a project resource
func (r *ResourceHandler) UpdateProjectResourceWithContext(ctx context.Context, project string, resource *models.Resource) (string, error) {
	return r.updateResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// DeleteProjectResource deletes a project resource
//
// Deprecated: Use DeleteProjectResourceWithContext instead
func (r *ResourceHandler) DeleteProjectResource(project string, resourceURI string) error {
	return r.DeleteProjectResourceWithContext(context.Background(), project, resourceURI)
}

// DeleteProjectResourceWithContext deletes a project resource
func (r *ResourceHandler) DeleteProjectResourceWithContext(ctx context.Context, project string, resourceURI string) error {
	return r.deleteResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource/"+url.QueryEscape(resourceURI))
}

// UpdateProjectResources updates multiple project resources
//
// Deprecated: Use UpdateProjectResourcesWithContext instead
func (r *ResourceHandler) UpdateProjectResources(project string, resources []*models.Resource) (string, error) {
	return r.UpdateProjectResourcesWithContext(context.Background(), project, resources)
}

// UpdateProjectResourcesWithContext updates multiple project resources
func (r *ResourceHandler) UpdateProjectResourcesWithContext(ctx context.Context, project string, resources []*models.Resource) (string, error) {
	return r.updateResources(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource", resources)
}

// CreateStageResources creates a stage resource
//
// Deprecated: Use CreateStageResourcesWithContext instead
func (r *ResourceHandler) CreateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	return r.CreateStageResourcesWithContext(context.Background(), project, stage, resources)
}

// CreateStageResourcesWithContext creates a stage resource
func (r *ResourceHandler) CreateStageResourcesWithContext(ctx context.Context, project string, stage string, resources []*models.Resource) (string, error) {
	return r.createResources(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", resources)
}

// GetStageResource retrieves a stage resource from the configuration service
//
// Deprecated: Use GetStageResourceWithContext instead
func (r *ResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return r.GetStageResourceWithContext(context.Background(), project, stage, resourceURI)
}

// GetStageResourceWithContext retrieves a stage resource from the configuration service
func (r *ResourceHandler) GetStageResourceWithContext(ctx context.Context, project string, stage string, resourceURI string) (*models.Resource, error) {
	return r.getResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource/"+url.QueryEscape(resourceURI))
}

// UpdateStageResource updates a stage resource
//
// Deprecated: Use UpdateStageResourceWithContext instead
func (r *ResourceHandler) UpdateStageResource(project string, stage string, resource *models.Resource) (string, error) {
	return r.UpdateStageResourceWithContext(context.Background(), project, stage, resource)
}

// UpdateStageResourceWithContext updates a stage resource
func (r *ResourceHandler) UpdateStageResourceWithContext(ctx context.Context, project string, stage string, resource *models.Resource) (string, error) {
	return r.updateResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateStageResources updates multiple stage resources
//
// Deprecated: Use UpdateStageResourcesWithContext instead
func (r *ResourceHandler) UpdateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	return r.UpdateStageResourcesWithContext(context.Background(), project, stage, resources)
}

// UpdateStageResourcesWithContext updates multiple stage resources
func (r *ResourceHandler) UpdateStageResourcesWithContext(ctx context.Context, project string, stage string, resources []*models.Resource) (string, error) {
	return r.updateResources(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", resources)
}

// DeleteStageResource deletes a stage resource
//
// Deprecated: Use DeleteStageResourceWithContext instead
func (r *ResourceHandler) DeleteStageResource(project string, stage string, resourceURI string) error {
	return r.DeleteStageResourceWithContext(context.Background(), project, stage, resourceURI)
}

// DeleteStageResourceWithContext deletes a stage resource
func (r *ResourceHandler) DeleteStageResourceWithContext(ctx context.Context, project string, stage string, resourceURI string) error {
	return r.deleteResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource/"+url.QueryEscape(resourceURI))
}

// CreateServiceResources creates a service resource
//
// Deprecated: Use CreateServiceResourcesWithContext instead
func (r *ResourceHandler) CreateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.CreateServiceResourcesWithContext(context.Background(), project, stage, service, resources)
}

// CreateServiceResourcesWithContext creates a service resource
func (r *ResourceHandler) CreateServiceResourcesWithContext(ctx context.Context, project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.createResources(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/resource", resources)
}

// GetServiceResource retrieves a service resource from the configuration service
//
// Deprecated: Use GetServiceResourceWithContext instead
func (r *ResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return r.GetServiceResourceWithContext(context.Background(), project, stage, service, resourceURI)
}

// GetServiceResourceWithContext retrieves a service resource from the configuration service
func (r *ResourceHandler) GetServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return r.getResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+url.QueryEscape(service)+"/resource/"+url.QueryEscape(resourceURI))
}

// UpdateServiceResource updates a service resource
//
// Deprecated: Use UpdateServiceResourceWithContext instead
func (r *ResourceHandler) UpdateServiceResource(project string, stage string, service string, resource *models.Resource) (string, error) {
	return r.UpdateServiceResourceWithContext(context.Background(), project, stage, service, resource)
}

// UpdateServiceResourceWithContext updates a service resource
func (r *ResourceHandler) UpdateServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resource *models.Resource) (string, error) {
	return r.updateResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+url.QueryEscape(service)+"/resource/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateServiceResources updates multiple service resources
//
// Deprecated: Use UpdateServiceResourcesWithContext instead
func (r *ResourceHandler) UpdateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.UpdateServiceResourcesWithContext(context.Background(), project, stage, service, resources)
}

// UpdateServiceResourcesWithContext updates multiple service resources
func (r *ResourceHandler) UpdateServiceResourcesWithContext(ctx context.Context, project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.updateResources(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+url.QueryEscape(service)+"/resource", resources)
}

// DeleteServiceResource deletes a service resource
//
// Deprecated: Use DeleteServiceResourceWithContext instead
func (r *ResourceHandler) DeleteServiceResource(project string, stage string, service string, resourceURI string) error {
	return r.DeleteServiceResourceWithContext(context.Background(), project, stage, service, resourceURI)
}

// DeleteServiceResourceWithContext deletes a service resource
func (r *ResourceHandler) DeleteServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resourceURI string) error {
	return r.deleteResource(ctx, r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+url.QueryEscape(service)+"/resource/"+url.QueryEscape(resourceURI))
}

func (r *ResourceHandler) createResources(ctx context.Context, uri string, resources []*models.Resource) (string, error) {
	return r.writeResources(ctx, uri, "POST", resources)
}

func (r *ResourceHandler) updateResources(ctx context.Context, uri string, resources []*models.Resource) (string, error) {
	return r.writeResources(ctx, uri, "PUT", resources)
}

func (r *ResourceHandler) writeResources(ctx context.Context, uri string, method string, resources []*models.Resource) (string, error) {

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
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(resourceStr))
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

func (r *ResourceHandler) updateResource(ctx context.Context, uri string, resource *models.Resource) (string, error) {
	return r.writeResource(ctx, uri, "PUT", resource)
}

func (r *ResourceHandler) writeResource(ctx context.Context, uri string, method string, resource *models.Resource) (string, error) {

	copiedResource := &models.Resource{ResourceURI: resource.ResourceURI, ResourceContent: b64.StdEncoding.EncodeToString([]byte(resource.ResourceContent))}

	resourceStr, err := json.Marshal(copiedResource)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(resourceStr))
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

func (r *ResourceHandler) getResource(ctx context.Context, uri string) (*models.Resource, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
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

func (r *ResourceHandler) deleteResource(ctx context.Context, uri string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequestWithContext(ctx, "DELETE", uri, nil)
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
//
// Deprecated: Use GetAllStageResourcesWithContext instead
func (r *ResourceHandler) GetAllStageResources(project string, stage string) ([]*models.Resource, error) {
	return r.GetAllStageResourcesWithContext(context.Background(), project, stage)
}

// GetAllStageResourcesWithContext returns a list of all resources.
func (r *ResourceHandler) GetAllStageResourcesWithContext(ctx context.Context, project string, stage string) ([]*models.Resource, error) {
	url, err := url.Parse(r.Scheme + "://" + r.getBaseURL() + "/v1/project/" + project + "/stage/" + stage + "/resource")
	if err != nil {
		return nil, err
	}
	return r.getAllResources(ctx, url)
}

// GetAllServiceResources returns a list of all resources.
//
// Deprecated: Use GetAllServiceResourcesWithContext instead
func (r *ResourceHandler) GetAllServiceResources(project string, stage string, service string) ([]*models.Resource, error) {
	return r.GetAllServiceResourcesWithContext(context.Background(), project, stage, service)
}

// GetAllServiceResourcesWithContext returns a list of all resources.
func (r *ResourceHandler) GetAllServiceResourcesWithContext(ctx context.Context, project string, stage string, service string) ([]*models.Resource, error) {
	url, err := url.Parse(r.Scheme + "://" + r.getBaseURL() + "/v1/project/" + project + "/stage/" + stage +
		"/service/" + service + "/resource/")
	if err != nil {
		return nil, err
	}
	return r.getAllResources(ctx, url)
}

func (r *ResourceHandler) getAllResources(ctx context.Context, u *url.URL) ([]*models.Resource, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resources := []*models.Resource{}

	nextPageKey := ""

	for {
		if nextPageKey != "" {
			q := u.Query()
			q.Set("nextPageKey", nextPageKey)
			u.RawQuery = q.Encode()
		}
		req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
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
