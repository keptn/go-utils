package v2

import (
	"bytes"
	"context"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

const pathToResource = "/resource"
const pathToService = "/service"
const pathToStage = "/stage"
const configurationServiceBaseURL = "configuration-service"

var ResourceNotFoundError = errors.New("Resource not found")

type ResourcesV1Interface interface {
	// CreateResources creates a resource for the specified entity.
	CreateResources(project string, stage string, service string, resources []*models.Resource) (*models.EventContext, *models.Error)

	// CreateResourcesWithContext creates a resource for the specified entity.
	CreateResourcesWithContext(ctx context.Context, project string, stage string, service string, resources []*models.Resource) (*models.EventContext, *models.Error)

	// CreateProjectResources creates multiple project resources.
	CreateProjectResources(project string, resources []*models.Resource) (string, error)

	// CreateProjectResourcesWithContext creates multiple project resources.
	CreateProjectResourcesWithContext(ctx context.Context, project string, resources []*models.Resource) (string, error)

	// GetProjectResource retrieves a project resource from the configuration service.
	// Deprecated: use GetResource instead.
	GetProjectResource(project string, resourceURI string) (*models.Resource, error)

	// GetProjectResourceWithContext retrieves a project resource from the configuration service.
	// Deprecated: use GetResourceWithContext instead.
	GetProjectResourceWithContext(ctx context.Context, project string, resourceURI string) (*models.Resource, error)

	// UpdateProjectResource updates a project resource.
	// Deprecated: use UpdateResource instead.
	UpdateProjectResource(project string, resource *models.Resource) (string, error)

	// UpdateProjectResourceWithContext updates a project resource.
	// Deprecated: use UpdateResourceWithContext instead.
	UpdateProjectResourceWithContext(ctx context.Context, project string, resource *models.Resource) (string, error)

	// DeleteProjectResource deletes a project resource.
	// Deprecated: use DeleteResource instead.
	DeleteProjectResource(project string, resourceURI string) error

	// DeleteProjectResourceWithContext deletes a project resource.
	// Deprecated: use DeleteResourceWithContext instead.
	DeleteProjectResourceWithContext(ctx context.Context, project string, resourceURI string) error

	// UpdateProjectResources updates multiple project resources.
	UpdateProjectResources(project string, resources []*models.Resource) (string, error)

	// UpdateProjectResourcesWithContext updates multiple project resources.
	UpdateProjectResourcesWithContext(ctx context.Context, project string, resources []*models.Resource) (string, error)

	// CreateStageResources creates a stage resource.
	// Deprecated: use CreateResource instead.
	CreateStageResources(project string, stage string, resources []*models.Resource) (string, error)

	// CreateStageResourcesWithContext creates a stage resource.
	// Deprecated: use CreateResourceWithContext instead.
	CreateStageResourcesWithContext(ctx context.Context, project string, stage string, resources []*models.Resource) (string, error)

	// GetStageResource retrieves a stage resource from the configuration service.
	// Deprecated: use GetResource instead.
	GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error)

	// GetStageResourceContext retrieves a stage resource from the configuration service.
	// Deprecated: use GetResourceWithContext instead.
	GetStageResourceWithContext(ctx context.Context, project string, stage string, resourceURI string) (*models.Resource, error)

	// UpdateStageResource updates a stage resource.
	// Deprecated: use UpdateResource instead.
	UpdateStageResource(project string, stage string, resource *models.Resource) (string, error)

	// UpdateStageResourceWithContext updates a stage resource.
	// Deprecated: use UpdateResourceWithContext instead.
	UpdateStageResourceWithContext(ctx context.Context, project string, stage string, resource *models.Resource) (string, error)

	// UpdateStageResources updates multiple stage resources.
	// Deprecated: use UpdateResource instead.
	UpdateStageResources(project string, stage string, resources []*models.Resource) (string, error)

	// UpdateStageResourcesWithContext updates multiple stage resources.
	// Deprecated: use UpdateResourceWithContext instead.
	UpdateStageResourcesWithContext(ctx context.Context, project string, stage string, resources []*models.Resource) (string, error)

	// DeleteStageResource deletes a stage resource.
	// Deprecated: use DeleteResource instead.
	DeleteStageResource(project string, stage string, resourceURI string) error

	// DeleteStageResourceWithContext deletes a stage resource.
	// Deprecated: use DeleteResourceWithContext instead.
	DeleteStageResourceWithContext(ctx context.Context, project string, stage string, resourceURI string) error

	// CreateServiceResources creates a service resource.
	// Deprecated: use CreateResource instead.
	CreateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error)

	// CreateServiceResourcesWithContext creates a service resource.
	// Deprecated: use CreateResourceWithContext instead.
	CreateServiceResourcesWithContext(ctx context.Context, project string, stage string, service string, resources []*models.Resource) (string, error)

	// GetServiceResource retrieves a service resource from the configuration service.
	// Deprecated: use GetResource instead.
	GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error)

	// GetServiceResourceWithContext retrieves a service resource from the configuration service.
	// Deprecated: use GetResourceWithContext instead.
	GetServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resourceURI string) (*models.Resource, error)

	// UpdateServiceResource updates a service resource.
	// Deprecated: use UpdateResource instead.
	UpdateServiceResource(project string, stage string, service string, resource *models.Resource) (string, error)

	// UpdateServiceResourceWithContext updates a service resource.
	// Deprecated: use UpdateResourceWithContext instead.
	UpdateServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resource *models.Resource) (string, error)

	// UpdateServiceResources updates multiple service resources.
	UpdateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error)

	// UpdateServiceResourcesWithContext updates multiple service resources.
	UpdateServiceResourcesWithContext(ctx context.Context, project string, stage string, service string, resources []*models.Resource) (string, error)

	// DeleteServiceResource deletes a service resource.
	// Deprecated: use DeleteResource instead.
	DeleteServiceResource(project string, stage string, service string, resourceURI string) error

	// DeleteServiceResourceWithContext deletes a service resource.
	// Deprecated: use DeleteResourceWithContext instead.
	DeleteServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resourceURI string) error

	// GetAllStageResources returns a list of all resources.
	GetAllStageResources(project string, stage string) ([]*models.Resource, error)

	// GetAllStageResourcesWithContext returns a list of all resources.
	GetAllStageResourcesWithContext(ctx context.Context, project string, stage string) ([]*models.Resource, error)

	// GetAllServiceResources returns a list of all resources.
	GetAllServiceResources(project string, stage string, service string) ([]*models.Resource, error)

	// GetAllServiceResourcesWithContext returns a list of all resources.
	GetAllServiceResourcesWithContext(ctx context.Context, project string, stage string, service string) ([]*models.Resource, error)
}

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

// ResourceScope contains the necessary information to get a resource
type ResourceScope struct {
	project  string
	stage    string
	service  string
	resource string
}

// NewResourceScope returns an empty ResourceScope to fill in calling Project Stage Service or Resource functions
func NewResourceScope() *ResourceScope {
	return &ResourceScope{}
}

// Project sets the resource scope project value
func (s *ResourceScope) Project(project string) *ResourceScope {
	s.project = project
	return s
}

// Stage sets the resource scope stage value
func (s *ResourceScope) Stage(stage string) *ResourceScope {
	s.stage = stage
	return s
}

// Service sets the resource scope service value
func (s *ResourceScope) Service(service string) *ResourceScope {
	s.service = service
	return s
}

// Resource sets the resource scope resource
func (s *ResourceScope) Resource(resource string) *ResourceScope {
	s.resource = resource
	return s
}

// GetProjectPath returns a string to construct the url to path eg. /<api-version>/project/<project-name>
//or an empty string if the project is not set
func (s *ResourceScope) GetProjectPath() string {
	return buildPath(v1ProjectPath, s.project)
}

// GetStagePath returns a string to construct the url to a stage eg. /stage/<stage-name>
//or an empty string if the stage is unset
func (s *ResourceScope) GetStagePath() string {
	return buildPath(pathToStage, s.stage)
}

// GetServicePath returns a string to construct the url to a service eg. /service/<service-name>
//or an empty string if the service is unset
func (s *ResourceScope) GetServicePath() string {
	return buildPath(pathToService, url.QueryEscape(s.service))
}

// GetResourcePath returns a string to construct the url to a resource eg. /resource/<escaped-resource-name>
//or /resource if the resource scope is empty
func (s *ResourceScope) GetResourcePath() string {
	path := pathToResource
	if s.resource != "" {
		path += "/" + url.QueryEscape(s.resource)
	}
	return path
}

func (r *ResourceHandler) buildResourceURI(scope ResourceScope) string {
	buildURI := r.Scheme + "://" + r.BaseURL + scope.GetProjectPath() + scope.GetStagePath() + scope.GetServicePath() + scope.GetResourcePath()
	return buildURI
}

// URIOption returns a function that modifies an url
type URIOption func(url string) string

// AppendQuery returns an option function that can modify an URI by appending a map of url query values
func AppendQuery(queryParams url.Values) URIOption {
	return func(buildURI string) string {
		if queryParams != nil {
			buildURI = buildURI + "?" + queryParams.Encode()
		}
		return buildURI
	}
}

func (r *ResourceHandler) applyOptions(buildURI string, options []URIOption) string {
	for _, option := range options {
		buildURI = option(buildURI)
	}
	return buildURI
}

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
	if !strings.HasSuffix(baseURL, configurationServiceBaseURL) {
		baseURL += "/" + configurationServiceBaseURL
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

// CreateResources creates a resource for the specified entity.
func (r *ResourceHandler) CreateResources(project string, stage string, service string, resources []*models.Resource) (*models.EventContext, *models.Error) {
	return r.CreateResourcesWithContext(context.TODO(), project, stage, service, resources)
}

// CreateResourcesWithContext creates a resource for the specified entity.
func (r *ResourceHandler) CreateResourcesWithContext(ctx context.Context, project string, stage string, service string, resources []*models.Resource) (*models.EventContext, *models.Error) {
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
		return postWithEventContext(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+service+pathToResource, requestStr, r)
	} else if project != "" && stage != "" && service == "" {
		return postWithEventContext(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource, requestStr, r)
	} else {
		return postWithEventContext(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+"/"+pathToResource, requestStr, r)
	}
}

// CreateProjectResources creates multiple project resources.
func (r *ResourceHandler) CreateProjectResources(project string, resources []*models.Resource) (string, error) {
	return r.CreateProjectResourcesWithContext(context.TODO(), project, resources)
}

// CreateProjectResourcesWithContext creates multiple project resources.
func (r *ResourceHandler) CreateProjectResourcesWithContext(ctx context.Context, project string, resources []*models.Resource) (string, error) {
	return r.createResources(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToResource, resources)
}

// GetProjectResource retrieves a project resource from the configuration service.
// Deprecated: use GetResource instead.
func (r *ResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return r.GetProjectResourceWithContext(context.TODO(), project, resourceURI)
}

// GetProjectResourceWithContext retrieves a project resource from the configuration service.
// Deprecated: use GetResourceWithContext instead.
func (r *ResourceHandler) GetProjectResourceWithContext(ctx context.Context, project string, resourceURI string) (*models.Resource, error) {
	buildURI := r.Scheme + "://" + r.BaseURL + v1ProjectPath + "/" + project + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.getResource(ctx, buildURI)
}

// UpdateProjectResource updates a project resource.
// Deprecated: use UpdateResource instead.
func (r *ResourceHandler) UpdateProjectResource(project string, resource *models.Resource) (string, error) {
	return r.UpdateProjectResourceWithContext(context.TODO(), project, resource)
}

// UpdateProjectResourceWithContext updates a project resource.
// Deprecated: use UpdateResourceWithContext instead.
func (r *ResourceHandler) UpdateProjectResourceWithContext(ctx context.Context, project string, resource *models.Resource) (string, error) {
	return r.updateResource(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// DeleteProjectResource deletes a project resource.
// Deprecated: use DeleteResource instead.
func (r *ResourceHandler) DeleteProjectResource(project string, resourceURI string) error {
	return r.DeleteProjectResourceWithContext(context.TODO(), project, resourceURI)
}

// DeleteProjectResourceWithContext deletes a project resource.
// Deprecated: use DeleteResourceWithContext instead.
func (r *ResourceHandler) DeleteProjectResourceWithContext(ctx context.Context, project string, resourceURI string) error {
	return r.deleteResource(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToResource+"/"+url.QueryEscape(resourceURI))
}

// UpdateProjectResources updates multiple project resources.
func (r *ResourceHandler) UpdateProjectResources(project string, resources []*models.Resource) (string, error) {
	return r.UpdateProjectResourcesWithContext(context.TODO(), project, resources)
}

// UpdateProjectResourcesWithContext updates multiple project resources.
func (r *ResourceHandler) UpdateProjectResourcesWithContext(ctx context.Context, project string, resources []*models.Resource) (string, error) {
	return r.updateResources(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToResource, resources)
}

// CreateStageResources creates a stage resource.
// Deprecated: use CreateResource instead.
func (r *ResourceHandler) CreateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	return r.CreateStageResourcesWithContext(context.TODO(), project, stage, resources)
}

// CreateStageResourcesWithContext creates a stage resource.
// Deprecated: use CreateResourceWithContext instead.
func (r *ResourceHandler) CreateStageResourcesWithContext(ctx context.Context, project string, stage string, resources []*models.Resource) (string, error) {
	return r.createResources(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource, resources)
}

// GetStageResource retrieves a stage resource from the configuration service.
// Deprecated: use GetResource instead.
func (r *ResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return r.GetStageResourceWithContext(context.TODO(), project, stage, resourceURI)
}

// GetStageResourceContext retrieves a stage resource from the configuration service.
// Deprecated: use GetResourceWithContext instead.
func (r *ResourceHandler) GetStageResourceWithContext(ctx context.Context, project string, stage string, resourceURI string) (*models.Resource, error) {
	buildURI := r.Scheme + "://" + r.BaseURL + v1ProjectPath + "/" + project + pathToStage + "/" + stage + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.getResource(ctx, buildURI)
}

// UpdateStageResource updates a stage resource.
// Deprecated: use UpdateResource instead.
func (r *ResourceHandler) UpdateStageResource(project string, stage string, resource *models.Resource) (string, error) {
	return r.UpdateStageResourceWithContext(context.TODO(), project, stage, resource)
}

// UpdateStageResourceWithContext updates a stage resource.
// Deprecated: use UpdateResourceWithContext instead.
func (r *ResourceHandler) UpdateStageResourceWithContext(ctx context.Context, project string, stage string, resource *models.Resource) (string, error) {
	return r.updateResource(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateStageResources updates multiple stage resources.
// Deprecated: use UpdateResource instead.
func (r *ResourceHandler) UpdateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	return r.UpdateStageResourcesWithContext(context.TODO(), project, stage, resources)
}

// UpdateStageResourcesWithContext updates multiple stage resources.
// Deprecated: use UpdateResourceWithContext instead.
func (r *ResourceHandler) UpdateStageResourcesWithContext(ctx context.Context, project string, stage string, resources []*models.Resource) (string, error) {
	return r.updateResources(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource, resources)
}

// DeleteStageResource deletes a stage resource.
// Deprecated: use DeleteResource instead.
func (r *ResourceHandler) DeleteStageResource(project string, stage string, resourceURI string) error {
	return r.DeleteStageResourceWithContext(context.TODO(), project, stage, resourceURI)
}

// DeleteStageResourceWithContext deletes a stage resource.
// Deprecated: use DeleteResourceWithContext instead.
func (r *ResourceHandler) DeleteStageResourceWithContext(ctx context.Context, project string, stage string, resourceURI string) error {
	return r.deleteResource(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource+"/"+url.QueryEscape(resourceURI))
}

// CreateServiceResources creates a service resource.
// Deprecated: use CreateResource instead.
func (r *ResourceHandler) CreateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.CreateServiceResourcesWithContext(context.TODO(), project, stage, service, resources)
}

// CreateServiceResourcesWithContext creates a service resource.
// Deprecated: use CreateResourceWithContext instead.
func (r *ResourceHandler) CreateServiceResourcesWithContext(ctx context.Context, project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.createResources(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+service+pathToResource, resources)
}

// GetServiceResource retrieves a service resource from the configuration service.
// Deprecated: use GetResource instead.
func (r *ResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return r.GetServiceResourceWithContext(context.TODO(), project, stage, service, resourceURI)
}

// GetServiceResourceWithContext retrieves a service resource from the configuration service.
// Deprecated: use GetResourceWithContext instead.
func (r *ResourceHandler) GetServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	buildURI := r.Scheme + "://" + r.BaseURL + v1ProjectPath + "/" + project + pathToStage + "/" + stage + pathToService + "/" + url.QueryEscape(service) + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.getResource(ctx, buildURI)
}

// UpdateServiceResource updates a service resource.
// Deprecated: use UpdateResource instead.
func (r *ResourceHandler) UpdateServiceResource(project string, stage string, service string, resource *models.Resource) (string, error) {
	return r.UpdateServiceResourceWithContext(context.TODO(), project, stage, service, resource)
}

// UpdateServiceResourceWithContext updates a service resource.
// Deprecated: use UpdateResourceWithContext instead.
func (r *ResourceHandler) UpdateServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resource *models.Resource) (string, error) {
	return r.updateResource(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+url.QueryEscape(service)+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateServiceResources updates multiple service resources.
func (r *ResourceHandler) UpdateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.UpdateServiceResourcesWithContext(context.TODO(), project, stage, service, resources)
}

// UpdateServiceResourcesWithContext updates multiple service resources.
func (r *ResourceHandler) UpdateServiceResourcesWithContext(ctx context.Context, project string, stage string, service string, resources []*models.Resource) (string, error) {
	return r.updateResources(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+url.QueryEscape(service)+pathToResource, resources)
}

// DeleteServiceResource deletes a service resource.
// Deprecated: use DeleteResource instead.
func (r *ResourceHandler) DeleteServiceResource(project string, stage string, service string, resourceURI string) error {
	return r.DeleteServiceResourceWithContext(context.TODO(), project, stage, service, resourceURI)
}

// DeleteServiceResourceWithContext deletes a service resource.
// Deprecated: use DeleteResourceWithContext instead.
func (r *ResourceHandler) DeleteServiceResourceWithContext(ctx context.Context, project string, stage string, service string, resourceURI string) error {
	return r.deleteResource(ctx, r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+url.QueryEscape(service)+pathToResource+"/"+url.QueryEscape(resourceURI))
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

	resourceStr, err := resReq.ToJSON()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(resourceStr))
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

func (r *ResourceHandler) updateResource(ctx context.Context, uri string, resource *models.Resource) (string, error) {
	return r.writeResource(ctx, uri, "PUT", resource)
}

func (r *ResourceHandler) writeResource(ctx context.Context, uri string, method string, resource *models.Resource) (string, error) {

	copiedResource := &models.Resource{ResourceURI: resource.ResourceURI, ResourceContent: b64.StdEncoding.EncodeToString([]byte(resource.ResourceContent))}

	resourceStr, err := copiedResource.ToJSON()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(resourceStr))
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

//GetResource returns a resource from the defined ResourceScope after applying all URI change configured in the options.
func (r *ResourceHandler) GetResource(scope ResourceScope, options ...URIOption) (*models.Resource, error) {
	return r.GetResourceWithContext(context.TODO(), scope, options...)
}

//GetResourceWithContext returns a resource from the defined ResourceScope after applying all URI change configured in the options.
func (r *ResourceHandler) GetResourceWithContext(ctx context.Context, scope ResourceScope, options ...URIOption) (*models.Resource, error) {
	buildURI := r.buildResourceURI(scope)
	return r.getResource(ctx, r.applyOptions(buildURI, options))
}

//DeleteResource delete a resource from the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) DeleteResource(scope ResourceScope, options ...URIOption) error {
	return r.DeleteResourceWithContext(context.TODO(), scope, options...)
}

//DeleteResourceWithContext delete a resource from the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) DeleteResourceWithContext(ctx context.Context, scope ResourceScope, options ...URIOption) error {
	buildURI := r.buildResourceURI(scope)
	return r.deleteResource(ctx, r.applyOptions(buildURI, options))
}

//UpdateResource updates a resource from the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) UpdateResource(resource *models.Resource, scope ResourceScope, options ...URIOption) (string, error) {
	return r.UpdateResourceWithContext(context.TODO(), resource, scope, options...)
}

//UpdateResourceWithContext updates a resource from the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) UpdateResourceWithContext(ctx context.Context, resource *models.Resource, scope ResourceScope, options ...URIOption) (string, error) {
	buildURI := r.buildResourceURI(scope)
	return r.updateResource(ctx, r.applyOptions(buildURI, options), resource)
}

//CreateResource creates one or more resources at the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) CreateResource(resource []*models.Resource, scope ResourceScope, options ...URIOption) (string, error) {
	return r.CreateResourceWithContext(context.TODO(), resource, scope, options...)
}

//CreateResourceWithContext creates one or more resources at the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) CreateResourceWithContext(ctx context.Context, resource []*models.Resource, scope ResourceScope, options ...URIOption) (string, error) {
	buildURI := r.buildResourceURI(scope)
	return r.createResources(ctx, r.applyOptions(buildURI, options), resource)
}

func (r *ResourceHandler) getResource(ctx context.Context, uri string) (*models.Resource, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	body, statusCode, status, mErr := get(ctx, uri, r)
	if mErr != nil {
		return nil, mErr.ToError()
	}

	if statusCode == 404 {
		// need to handle this case differently (e.g. https://github.com/keptn/keptn/issues/1480)
		return nil, ResourceNotFoundError
	}
	if !(statusCode >= 200 && statusCode < 300) {
		if len(body) > 0 {
			return nil, handleErrStatusCode(statusCode, body).ToError()
		}

		return nil, buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", statusCode, status)).ToError()
	}

	resource := &models.Resource{}
	if err := resource.FromJSON(body); err != nil {
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

func (r *ResourceHandler) deleteResource(ctx context.Context, uri string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequestWithContext(ctx, "DELETE", uri, nil)
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
	return r.GetAllStageResourcesWithContext(context.TODO(), project, stage)
}

// GetAllStageResourcesWithContext returns a list of all resources.
func (r *ResourceHandler) GetAllStageResourcesWithContext(ctx context.Context, project string, stage string) ([]*models.Resource, error) {
	myURL, err := url.Parse(r.Scheme + "://" + r.getBaseURL() + v1ProjectPath + "/" + project + pathToStage + "/" + stage + pathToResource)
	if err != nil {
		return nil, err
	}
	return r.getAllResources(ctx, myURL)
}

// GetAllServiceResources returns a list of all resources.
func (r *ResourceHandler) GetAllServiceResources(project string, stage string, service string) ([]*models.Resource, error) {
	return r.GetAllServiceResourcesWithContext(context.TODO(), project, stage, service)
}

// GetAllServiceResourcesWithContext returns a list of all resources.
func (r *ResourceHandler) GetAllServiceResourcesWithContext(ctx context.Context, project string, stage string, service string) ([]*models.Resource, error) {
	myURL, err := url.Parse(r.Scheme + "://" + r.getBaseURL() + v1ProjectPath + "/" + project + pathToStage + "/" + stage +
		pathToService + "/" + service + pathToResource + "/")
	if err != nil {
		return nil, err
	}
	return r.getAllResources(ctx, myURL)
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

		body, mErr := getAndExpectOK(ctx, u.String(), r)
		if mErr != nil {
			return nil, mErr.ToError()
		}

		received := &models.Resources{}
		if err := received.FromJSON(body); err != nil {
			return nil, err
		}

		resources = append(resources, received.Resources...)

		if received.NextPageKey == "" || received.NextPageKey == "0" {
			break
		}
		nextPageKey = received.NextPageKey
	}

	return resources, nil
}

func buildPath(base, name string) string {
	path := ""
	if name != "" {
		path = base + "/" + name
	}
	return path
}
