package utils

import (
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
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

// NewResourceHandler returns a new ResourceHandler
func NewResourceHandler(baseURL string) *ResourceHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &ResourceHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     "https",
	}
}

// NewAuthenticatedResourceHandler returns a new ResourceHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedResourceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ResourceHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
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
		return post(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/resource", requestStr, r)
	} else if project != "" && stage != "" && service == "" {
		return post(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", requestStr, r)
	} else {
		return post(r.Scheme+"://"+r.BaseURL+"/v1/project/"+project+"/resource", requestStr, r)
	}
}
