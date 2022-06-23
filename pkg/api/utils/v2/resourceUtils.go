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
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const pathToResource = "/resource"
const pathToService = "/service"
const pathToStage = "/stage"
const configurationServiceBaseURL = "configuration-service"

var ResourceNotFoundError = errors.New("Resource not found")

// ResourcesCreateResourcesOptions are options for ResourcesInterface.CreateResources().
type ResourcesCreateResourcesOptions struct{}

// ResourcesCreateProjectResourcesOptions are options for ResourcesInterface.CreateProjectResources().
type ResourcesCreateProjectResourcesOptions struct{}

// ResourcesUpdateProjectResourcesOptions are options for ResourcesInterface.UpdateProjectResources().
type ResourcesUpdateProjectResourcesOptions struct{}

// ResourcesUpdateServiceResourcesOptions are options for ResourcesInterface.UpdateServiceResources().
type ResourcesUpdateServiceResourcesOptions struct{}

// ResourcesGetAllStageResourcesOptions are options for ResourcesInterface.GetAllStageResources().
type ResourcesGetAllStageResourcesOptions struct{}

// ResourcesGetAllServiceResourcesOptions are options for ResourcesInterface.GetAllServiceResources().
type ResourcesGetAllServiceResourcesOptions struct{}

// ResourcesGetResourceOptions are options for ResourcesInterface.GetResource().
type ResourcesGetResourceOptions struct {
	// URIOptions modify the resource's URI.
	URIOptions []URIOption
}

// ResourcesDeleteResourceOptions are options for ResourcesInterface.DeleteResource().
type ResourcesDeleteResourceOptions struct {
	// URIOptions modify the resource's URI.
	URIOptions []URIOption
}

// ResourcesUpdateResourceOptions are options for ResourcesInterface.UpdateResource().
type ResourcesUpdateResourceOptions struct {
	// URIOptions modify the resource's URI.
	URIOptions []URIOption
}

// ResourcesCreateResourceOptions are options for ResourcesInterface.CreateResource().
type ResourcesCreateResourceOptions struct {
	// URIOptions modify the resource's URI.
	URIOptions []URIOption
}

type ResourcesInterface interface {
	// CreateResources creates a resource for the specified entity.
	CreateResources(ctx context.Context, project string, stage string, service string, resources []*models.Resource, opts ResourcesCreateResourcesOptions) (*models.EventContext, *models.Error)

	// CreateProjectResources creates multiple project resources.
	CreateProjectResources(ctx context.Context, project string, resources []*models.Resource, opts ResourcesCreateProjectResourcesOptions) (string, error)

	// UpdateProjectResources updates multiple project resources.
	UpdateProjectResources(ctx context.Context, project string, resources []*models.Resource, opts ResourcesUpdateProjectResourcesOptions) (string, error)

	// UpdateServiceResources updates multiple service resources.
	UpdateServiceResources(ctx context.Context, project string, stage string, service string, resources []*models.Resource, opts ResourcesUpdateServiceResourcesOptions) (string, error)

	// GetAllStageResources returns a list of all resources.
	GetAllStageResources(ctx context.Context, project string, stage string, opts ResourcesGetAllStageResourcesOptions) ([]*models.Resource, error)

	// GetAllServiceResources returns a list of all resources.
	GetAllServiceResources(ctx context.Context, project string, stage string, service string, opts ResourcesGetAllServiceResourcesOptions) ([]*models.Resource, error)

	// GetResource returns a resource from the defined ResourceScope.
	GetResource(ctx context.Context, scope ResourceScope, opts ResourcesGetResourceOptions) (*models.Resource, error)

	// DeleteResource delete a resource from the URI defined by ResourceScope.
	DeleteResource(ctx context.Context, scope ResourceScope, opts ResourcesDeleteResourceOptions) error

	// UpdateResource updates a resource from the URI defined by ResourceScope.
	UpdateResource(ctx context.Context, resource *models.Resource, scope ResourceScope, opts ResourcesUpdateResourceOptions) (string, error)

	// CreateResource creates one or more resources at the URI defined by ResourceScope.
	CreateResource(ctx context.Context, resource []*models.Resource, scope ResourceScope, opts ResourcesCreateResourceOptions) (string, error)
}

// ResourceHandler handles resources
type ResourceHandler struct {
	baseURL    string
	authToken  string
	authHeader string
	httpClient *http.Client
	scheme     string
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
	buildURI := r.scheme + "://" + r.baseURL + scope.GetProjectPath() + scope.GetStagePath() + scope.GetServicePath() + scope.GetResourcePath()
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
	return NewResourceHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewResourceHandlerWithHTTPClient returns a new ResourceHandler which sends all requests directly to the configuration-service using the specified http.Client
func NewResourceHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *ResourceHandler {
	return createResourceHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedResourceHandler returns a new ResourceHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
func NewAuthenticatedResourceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ResourceHandler {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, configurationServiceBaseURL) {
		baseURL += "/" + configurationServiceBaseURL
	}

	return createResourceHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createResourceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ResourceHandler {
	return &ResourceHandler{
		baseURL:    httputils.TrimHTTPScheme(baseURL),
		authHeader: authHeader,
		authToken:  authToken,
		httpClient: httpClient,
		scheme:     scheme,
	}
}

func (r *ResourceHandler) getBaseURL() string {
	return r.baseURL
}

func (r *ResourceHandler) getAuthToken() string {
	return r.authToken
}

func (r *ResourceHandler) getAuthHeader() string {
	return r.authHeader
}

func (r *ResourceHandler) getHTTPClient() *http.Client {
	return r.httpClient
}

// CreateResources creates a resource for the specified entity.
func (r *ResourceHandler) CreateResources(ctx context.Context, project string, stage string, service string, resources []*models.Resource, opts ResourcesCreateResourcesOptions) (*models.EventContext, *models.Error) {
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
		return postWithEventContext(ctx, r.scheme+"://"+r.baseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+service+pathToResource, requestStr, r)
	} else if project != "" && stage != "" && service == "" {
		return postWithEventContext(ctx, r.scheme+"://"+r.baseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToResource, requestStr, r)
	} else {
		return postWithEventContext(ctx, r.scheme+"://"+r.baseURL+v1ProjectPath+"/"+project+"/"+pathToResource, requestStr, r)
	}
}

// CreateProjectResources creates multiple project resources.
func (r *ResourceHandler) CreateProjectResources(ctx context.Context, project string, resources []*models.Resource, opts ResourcesCreateProjectResourcesOptions) (string, error) {
	return r.CreateResourcesByURI(ctx, r.scheme+"://"+r.baseURL+v1ProjectPath+"/"+project+pathToResource, resources)
}

// UpdateProjectResources updates multiple project resources.
func (r *ResourceHandler) UpdateProjectResources(ctx context.Context, project string, resources []*models.Resource, opts ResourcesUpdateProjectResourcesOptions) (string, error) {
	return r.UpdateResourcesByURI(ctx, r.scheme+"://"+r.baseURL+v1ProjectPath+"/"+project+pathToResource, resources)
}

// UpdateServiceResources updates multiple service resources.
func (r *ResourceHandler) UpdateServiceResources(ctx context.Context, project string, stage string, service string, resources []*models.Resource, opts ResourcesUpdateServiceResourcesOptions) (string, error) {
	return r.UpdateResourcesByURI(ctx, r.scheme+"://"+r.baseURL+v1ProjectPath+"/"+project+pathToStage+"/"+stage+pathToService+"/"+url.QueryEscape(service)+pathToResource, resources)
}

func (r *ResourceHandler) CreateResourcesByURI(ctx context.Context, uri string, resources []*models.Resource) (string, error) {
	return r.writeResources(ctx, uri, "POST", resources)
}

func (r *ResourceHandler) UpdateResourcesByURI(ctx context.Context, uri string, resources []*models.Resource) (string, error) {
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

	resp, err := r.httpClient.Do(req)
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

func (r *ResourceHandler) UpdateResourceByURI(ctx context.Context, uri string, resource *models.Resource) (string, error) {
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

	resp, err := r.httpClient.Do(req)
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

// GetResource returns a resource from the defined ResourceScope.
func (r *ResourceHandler) GetResource(ctx context.Context, scope ResourceScope, opts ResourcesGetResourceOptions) (*models.Resource, error) {
	buildURI := r.buildResourceURI(scope)
	return r.GetResourceByURI(ctx, r.applyOptions(buildURI, opts.URIOptions))
}

//DeleteResource delete a resource from the URI defined by ResourceScope.
func (r *ResourceHandler) DeleteResource(ctx context.Context, scope ResourceScope, opts ResourcesDeleteResourceOptions) error {
	buildURI := r.buildResourceURI(scope)
	return r.DeleteResourceByURI(ctx, r.applyOptions(buildURI, opts.URIOptions))
}

//UpdateResource updates a resource from the URI defined by ResourceScope.
func (r *ResourceHandler) UpdateResource(ctx context.Context, resource *models.Resource, scope ResourceScope, opts ResourcesUpdateResourceOptions) (string, error) {
	buildURI := r.buildResourceURI(scope)
	return r.UpdateResourceByURI(ctx, r.applyOptions(buildURI, opts.URIOptions), resource)
}

//CreateResource creates one or more resources at the URI defined by ResourceScope.
func (r *ResourceHandler) CreateResource(ctx context.Context, resource []*models.Resource, scope ResourceScope, opts ResourcesCreateResourceOptions) (string, error) {
	buildURI := r.buildResourceURI(scope)
	return r.CreateResourcesByURI(ctx, r.applyOptions(buildURI, opts.URIOptions), resource)
}

func (r *ResourceHandler) GetResourceByURI(ctx context.Context, uri string) (*models.Resource, error) {
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

func (r *ResourceHandler) DeleteResourceByURI(ctx context.Context, uri string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequestWithContext(ctx, "DELETE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, r)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetAllStageResources returns a list of all resources.
func (r *ResourceHandler) GetAllStageResources(ctx context.Context, project string, stage string, opts ResourcesGetAllStageResourcesOptions) ([]*models.Resource, error) {
	myURL, err := url.Parse(r.scheme + "://" + r.getBaseURL() + v1ProjectPath + "/" + project + pathToStage + "/" + stage + pathToResource)
	if err != nil {
		return nil, err
	}
	return r.getAllResources(ctx, myURL)
}

// GetAllServiceResources returns a list of all resources.
func (r *ResourceHandler) GetAllServiceResources(ctx context.Context, project string, stage string, service string, opts ResourcesGetAllServiceResourcesOptions) ([]*models.Resource, error) {
	myURL, err := url.Parse(r.scheme + "://" + r.getBaseURL() + v1ProjectPath + "/" + project + pathToStage + "/" + stage +
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
