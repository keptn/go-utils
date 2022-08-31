package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const pathToResource = "/resource"
const pathToService = "/service"
const pathToStage = "/stage"
const configurationServiceBaseURL = "resource-service"

var ResourceNotFoundError = v2.ResourceNotFoundError

type ResourcesV1Interface interface {
	// CreateResources creates a resource for the specified entity.
	CreateResources(project string, stage string, service string, resources []*models.Resource) (*models.EventContext, *models.Error)

	// CreateProjectResources creates multiple project resources.
	CreateProjectResources(project string, resources []*models.Resource) (string, error)

	// GetProjectResource retrieves a project resource from the configuration service.
	// Deprecated: use GetResource instead.
	GetProjectResource(project string, resourceURI string) (*models.Resource, error)

	// UpdateProjectResource updates a project resource.
	// Deprecated: use UpdateResource instead.
	UpdateProjectResource(project string, resource *models.Resource) (string, error)

	// DeleteProjectResource deletes a project resource.
	// Deprecated: use DeleteResource instead.
	DeleteProjectResource(project string, resourceURI string) error

	// UpdateProjectResources updates multiple project resources.
	UpdateProjectResources(project string, resources []*models.Resource) (string, error)

	// CreateStageResources creates a stage resource.
	// Deprecated: use CreateResource instead.
	CreateStageResources(project string, stage string, resources []*models.Resource) (string, error)

	// GetStageResource retrieves a stage resource from the configuration service.
	// Deprecated: use GetResource instead.
	GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error)

	// UpdateStageResource updates a stage resource.
	// Deprecated: use UpdateResource instead.
	UpdateStageResource(project string, stage string, resource *models.Resource) (string, error)

	// UpdateStageResources updates multiple stage resources.
	// Deprecated: use UpdateResource instead.
	UpdateStageResources(project string, stage string, resources []*models.Resource) (string, error)

	// DeleteStageResource deletes a stage resource.
	// Deprecated: use DeleteResource instead.
	DeleteStageResource(project string, stage string, resourceURI string) error

	// CreateServiceResources creates a service resource.
	// Deprecated: use CreateResource instead.
	CreateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error)

	// GetServiceResource retrieves a service resource from the configuration service.
	// Deprecated: use GetResource instead.
	GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error)

	// UpdateServiceResource updates a service resource.
	// Deprecated: use UpdateResource instead.
	UpdateServiceResource(project string, stage string, service string, resource *models.Resource) (string, error)

	// UpdateServiceResources updates multiple service resources.
	UpdateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error)

	// DeleteServiceResource deletes a service resource.
	// Deprecated: use DeleteResource instead.
	DeleteServiceResource(project string, stage string, service string, resourceURI string) error

	// GetAllStageResources returns a list of all resources.
	GetAllStageResources(project string, stage string) ([]*models.Resource, error)

	// GetAllServiceResources returns a list of all resources.
	GetAllServiceResources(project string, stage string, service string) ([]*models.Resource, error)
}

// ResourceHandler handles resources
type ResourceHandler struct {
	resourceHandler *v2.ResourceHandler
	BaseURL         string
	AuthToken       string
	AuthHeader      string
	HTTPClient      *http.Client
	Scheme          string
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

// GetProject returns the project name
func (s *ResourceScope) GetProject() string {
	return s.project
}

// GetStage returns the stage name
func (s *ResourceScope) GetStage() string {
	return s.stage
}

// GetService returns the service name
func (s *ResourceScope) GetService() string {
	return s.service
}

// GetResource returns the resource name
func (s *ResourceScope) GetResource() string {
	return s.resource
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

// NewResourceHandler returns a new ResourceHandler which sends all requests directly to the resource-service
func NewResourceHandler(baseURL string) *ResourceHandler {
	return NewResourceHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewResourceHandlerWithHTTPClient returns a new ResourceHandler which sends all requests directly to the resource-service using the specified http.Client
func NewResourceHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *ResourceHandler {
	return &ResourceHandler{
		BaseURL:         httputils.TrimHTTPScheme(baseURL),
		HTTPClient:      httpClient,
		Scheme:          "http",
		resourceHandler: v2.NewResourceHandlerWithHTTPClient(baseURL, httpClient),
	}
}

// NewAuthenticatedResourceHandler returns a new ResourceHandler that authenticates at the api via the provided token
// and sends all requests directly to the resource-service
// Deprecated: use APISet instead
func NewAuthenticatedResourceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ResourceHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedResourceHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedResourceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ResourceHandler {
	v2ResourceHandler := v2.NewAuthenticatedResourceHandler(baseURL, authToken, authHeader, httpClient, scheme)

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, configurationServiceBaseURL) {
		baseURL += "/" + configurationServiceBaseURL
	}

	return &ResourceHandler{
		BaseURL:         httputils.TrimHTTPScheme(baseURL),
		AuthHeader:      authHeader,
		AuthToken:       authToken,
		HTTPClient:      httpClient,
		Scheme:          scheme,
		resourceHandler: v2ResourceHandler,
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
	r.ensureHandlerIsSet()
	return r.resourceHandler.CreateResources(context.TODO(), project, stage, service, resources, v2.ResourcesCreateResourcesOptions{})
}

// CreateProjectResources creates multiple project resources.
func (r *ResourceHandler) CreateProjectResources(project string, resources []*models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.CreateProjectResources(context.TODO(), project, resources, v2.ResourcesCreateProjectResourcesOptions{})
}

// GetProjectResource retrieves a project resource from the configuration service.
// Deprecated: use GetResource instead.
func (r *ResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	r.ensureHandlerIsSet()
	buildURI := r.Scheme + "://" + r.BaseURL + v1ProjectPath + "/" + project + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.resourceHandler.GetResourceByURI(context.TODO(), buildURI)
}

// UpdateProjectResource updates a project resource.
// Deprecated: use UpdateResource instead.
func (r *ResourceHandler) UpdateProjectResource(project string, resource *models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.UpdateResourceByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// DeleteProjectResource deletes a project resource.
// Deprecated: use DeleteResource instead.
func (r *ResourceHandler) DeleteProjectResource(project string, resourceURI string) error {
	r.ensureHandlerIsSet()
	return r.resourceHandler.DeleteResourceByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToResource+"/"+url.QueryEscape(resourceURI))
}

// UpdateProjectResources updates multiple project resources.
func (r *ResourceHandler) UpdateProjectResources(project string, resources []*models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.UpdateProjectResources(context.TODO(), project, resources, v2.ResourcesUpdateProjectResourcesOptions{})
}

// CreateStageResources creates a stage resource.
// Deprecated: use CreateResource instead.
func (r *ResourceHandler) CreateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.CreateResourcesByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource, resources)
}

// GetStageResource retrieves a stage resource from the configuration service.
// Deprecated: use GetResource instead.
func (r *ResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	r.ensureHandlerIsSet()
	buildURI := r.Scheme + "://" + r.BaseURL + v1ProjectPath + "/" + project + pathToStage + "/" + stage + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.resourceHandler.GetResourceByURI(context.TODO(), buildURI)
}

// UpdateStageResource updates a stage resource.
// Deprecated: use UpdateResource instead.
func (r *ResourceHandler) UpdateStageResource(project string, stage string, resource *models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.UpdateResourceByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateStageResources updates multiple stage resources.
// Deprecated: use UpdateResource instead.
func (r *ResourceHandler) UpdateStageResources(project string, stage string, resources []*models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.UpdateResourcesByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource, resources)
}

// DeleteStageResource deletes a stage resource.
// Deprecated: use DeleteResource instead.
func (r *ResourceHandler) DeleteStageResource(project string, stage string, resourceURI string) error {
	r.ensureHandlerIsSet()
	return r.resourceHandler.DeleteResourceByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource+"/"+url.QueryEscape(resourceURI))
}

// CreateServiceResources creates a service resource.
// Deprecated: use CreateResource instead.
func (r *ResourceHandler) CreateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.CreateResourcesByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+service+pathToResource, resources)
}

// GetServiceResource retrieves a service resource from the configuration service.
// Deprecated: use GetResource instead.
func (r *ResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	r.ensureHandlerIsSet()
	buildURI := r.Scheme + "://" + r.BaseURL + v1ProjectPath + "/" + project + pathToStage + "/" + stage + pathToService + "/" + url.QueryEscape(service) + pathToResource + "/" + url.QueryEscape(resourceURI)
	return r.resourceHandler.GetResourceByURI(context.TODO(), buildURI)
}

// UpdateServiceResource updates a service resource.
// Deprecated: use UpdateResource instead.
func (r *ResourceHandler) UpdateServiceResource(project string, stage string, service string, resource *models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.UpdateResourceByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+url.QueryEscape(service)+pathToResource+"/"+url.QueryEscape(*resource.ResourceURI), resource)
}

// UpdateServiceResources updates multiple service resources.
func (r *ResourceHandler) UpdateServiceResources(project string, stage string, service string, resources []*models.Resource) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.UpdateServiceResources(context.TODO(), project, stage, service, resources, v2.ResourcesUpdateServiceResourcesOptions{})
}

// DeleteServiceResource deletes a service resource.
// Deprecated: use DeleteResource instead.
func (r *ResourceHandler) DeleteServiceResource(project string, stage string, service string, resourceURI string) error {
	r.ensureHandlerIsSet()
	return r.resourceHandler.DeleteResourceByURI(context.TODO(), r.Scheme+"://"+r.BaseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+url.QueryEscape(service)+pathToResource+"/"+url.QueryEscape(resourceURI))
}

//GetResource returns a resource from the defined ResourceScope after applying all URI change configured in the options.
func (r *ResourceHandler) GetResource(scope ResourceScope, options ...URIOption) (*models.Resource, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.GetResource(context.TODO(), toV2ResourceScope(scope), v2.ResourcesGetResourceOptions{URIOptions: toV2URIOptions(options)})
}

//DeleteResource delete a resource from the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) DeleteResource(scope ResourceScope, options ...URIOption) error {
	r.ensureHandlerIsSet()
	return r.resourceHandler.DeleteResource(context.TODO(), toV2ResourceScope(scope), v2.ResourcesDeleteResourceOptions{URIOptions: toV2URIOptions(options)})
}

//UpdateResource updates a resource from the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) UpdateResource(resource *models.Resource, scope ResourceScope, options ...URIOption) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.UpdateResource(context.TODO(), resource, toV2ResourceScope(scope), v2.ResourcesUpdateResourceOptions{URIOptions: toV2URIOptions(options)})
}

//CreateResource creates one or more resources at the URI defined by ResourceScope and modified by the URIOption.
func (r *ResourceHandler) CreateResource(resource []*models.Resource, scope ResourceScope, options ...URIOption) (string, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.CreateResource(context.TODO(), resource, toV2ResourceScope(scope), v2.ResourcesCreateResourceOptions{URIOptions: toV2URIOptions(options)})
}

// GetAllStageResources returns a list of all resources.
func (r *ResourceHandler) GetAllStageResources(project string, stage string) ([]*models.Resource, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.GetAllStageResources(context.TODO(), project, stage, v2.ResourcesGetAllStageResourcesOptions{})
}

// GetAllServiceResources returns a list of all resources.
func (r *ResourceHandler) GetAllServiceResources(project string, stage string, service string) ([]*models.Resource, error) {
	r.ensureHandlerIsSet()
	return r.resourceHandler.GetAllServiceResources(context.TODO(), project, stage, service, v2.ResourcesGetAllServiceResourcesOptions{})
}

func buildPath(base, name string) string {
	path := ""
	if name != "" {
		path = base + "/" + name
	}
	return path
}

func toV2URIOptions(uriOptions []URIOption) []v2.URIOption {
	var v2URIOptions []v2.URIOption
	for _, v := range uriOptions {
		v2URIOptions = append(v2URIOptions, v2.URIOption(v))
	}
	return v2URIOptions
}

func toV2ResourceScope(scope ResourceScope) v2.ResourceScope {
	return *(v2.NewResourceScope().Project(scope.project).Stage(scope.stage).Service(scope.service).Resource(scope.resource))
}

func (r *ResourceHandler) ensureHandlerIsSet() {
	if r.resourceHandler != nil {
		return
	}

	if r.AuthToken != "" {
		r.resourceHandler = v2.NewAuthenticatedResourceHandler(r.BaseURL, r.AuthToken, r.AuthHeader, r.HTTPClient, r.Scheme)
	} else {
		r.resourceHandler = v2.NewResourceHandlerWithHTTPClient(r.BaseURL, r.HTTPClient)
	}
}
