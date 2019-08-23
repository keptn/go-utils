package utils

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ResourceHandler handles resources
type ResourceHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
}

// Resource contains resource information
type Resource struct {
	ResourceContent string `json:"resourceContent"`
	ResourceURI     string `json:"resourceUri"`
}

type versionResponse struct {
	Version string `json:"version"`
}

type resourceRequest struct {
	Resources []*Resource `json:"resources"`
}

// NewResourceHandler returns a new ResourceHandler
func NewResourceHandler(baseURL string) *ResourceHandler {
	return &ResourceHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
	}
}

// NewAuthenticatedResourceHandler returns a new ResourceHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedResourceHandler(baseURL string, authToken string, authHeader string) *ResourceHandler {
	return &ResourceHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
	}
}

// CreateProjectResources creates multiple project resources
func (r *ResourceHandler) CreateProjectResources(project string, resources []*Resource) (string, error) {
	return r.createResources("http://"+r.BaseURL+"/v1/project/"+project+"/resource", resources)
}

// GetProjectResource retrieves a project resource from the configuration service
func (r *ResourceHandler) GetProjectResource(project string, resourceURI string) (*Resource, error) {
	return r.getResource("http://" + r.BaseURL + "/v1/project/" + project + "/resource/" + url.QueryEscape(resourceURI))
}

// UpdateProjectResource updates a project resource
func (r *ResourceHandler) UpdateProjectResource(project string, resource *Resource) (string, error) {
	return r.updateResource("http://"+r.BaseURL+"/v1/project/"+project+"/resource/"+url.QueryEscape(resource.ResourceURI), resource)
}

// DeleteProjectResource deletes a project resource
func (r *ResourceHandler) DeleteProjectResource(project string, resourceURI string) error {
	return r.deleteResource("http://" + r.BaseURL + "/v1/project/" + project + "/resource/" + url.QueryEscape(resourceURI))
}

// UpdateProjectResources updates multiple project resources
func (r *ResourceHandler) UpdateProjectResources(project string, resources []*Resource) (string, error) {
	return r.updateResources("http://"+r.BaseURL+"/v1/project/"+project+"/resource", resources)
}

// CreateStageResources creates a stage resource
func (r *ResourceHandler) CreateStageResources(project string, stage string, resources []*Resource) (string, error) {
	return r.createResources("http://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", resources)
}

// GetStageResource retrieves a stage resource from the configuration service
func (r *ResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*Resource, error) {
	return r.getResource("http://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/resource/" + url.QueryEscape(resourceURI))
}

// UpdateStageResource updates a stage resource
func (r *ResourceHandler) UpdateStageResource(project string, stage string, resource *Resource) (string, error) {
	return r.updateResource("http://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource/"+url.QueryEscape(resource.ResourceURI), resource)
}

// UpdateStageResources updates multiple stage resources
func (r *ResourceHandler) UpdateStageResources(project string, stage string, resources []*Resource) (string, error) {
	return r.updateResources("http://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/resource", resources)
}

// DeleteStageResource deletes a stage resource
func (r *ResourceHandler) DeleteStageResource(project string, stage string, resourceURI string) error {
	return r.deleteResource("http://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/resource/" + url.QueryEscape(resourceURI))
}

// CreateServiceResources creates a service resource
func (r *ResourceHandler) CreateServiceResources(project string, stage string, service string, resources []*Resource) (string, error) {
	return r.createResources("http://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/resource", resources)
}

// GetServiceResource retrieves a service resource from the configuration service
func (r *ResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*Resource, error) {
	return r.getResource("http://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/service/" + url.QueryEscape(service) + "/resource/" + url.QueryEscape(resourceURI))
}

// UpdateServiceResource updates a service resource
func (r *ResourceHandler) UpdateServiceResource(project string, stage string, service string, resource *Resource) (string, error) {
	return r.updateResource("http://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+url.QueryEscape(service)+"/resource/"+url.QueryEscape(resource.ResourceURI), resource)
}

// UpdateServiceResources updates multiple service resources
func (r *ResourceHandler) UpdateServiceResources(project string, stage string, service string, resources []*Resource) (string, error) {
	return r.updateResources("http://"+r.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service/"+url.QueryEscape(service)+"/resource", resources)
}

// DeleteServiceResource deletes a service resource
func (r *ResourceHandler) DeleteServiceResource(project string, stage string, service string, resourceURI string) error {
	return r.deleteResource("http://" + r.BaseURL + "/v1/project/" + project + "/stage/" + stage + "/service/" + url.QueryEscape(service) + "/resource/" + url.QueryEscape(resourceURI))
}

func (r *ResourceHandler) createResources(uri string, resources []*Resource) (string, error) {
	return r.writeResources(uri, "POST", resources)
}

func (r *ResourceHandler) updateResources(uri string, resources []*Resource) (string, error) {
	return r.writeResources(uri, "PUT", resources)
}

func (r *ResourceHandler) writeResources(uri string, method string, resources []*Resource) (string, error) {
	for i := range resources {
		resources[i].ResourceContent = b64.StdEncoding.EncodeToString([]byte(resources[i].ResourceContent))
	}
	resReq := &resourceRequest{
		Resources: resources,
	}

	resourceStr, err := json.Marshal(resReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(resourceStr))
	req.Header.Set("Content-Type", "application/json")
	if r.AuthHeader != "" && r.AuthToken != "" {
		req.Header.Set(r.AuthHeader, r.AuthToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var version versionResponse
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

func (r *ResourceHandler) updateResource(uri string, resource *Resource) (string, error) {
	return r.writeResource(uri, "PUT", resource)
}

func (r *ResourceHandler) writeResource(uri string, method string, resource *Resource) (string, error) {
	resource.ResourceContent = b64.StdEncoding.EncodeToString([]byte(resource.ResourceContent))
	resourceStr, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(resourceStr))
	req.Header.Set("Content-Type", "application/json")
	if r.AuthHeader != "" && r.AuthToken != "" {
		req.Header.Set(r.AuthHeader, r.AuthToken)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var version versionResponse
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

func (r *ResourceHandler) getResource(uri string) (*Resource, error) {
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("Content-Type", "application/json")
	if r.AuthHeader != "" && r.AuthToken != "" {
		req.Header.Set(r.AuthHeader, r.AuthToken)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var resource Resource
	err = json.Unmarshal(body, &resource)
	if err != nil {
		return nil, err
	}
	decodedStr, err := b64.StdEncoding.DecodeString(resource.ResourceContent)
	if err != nil {
		return nil, err
	}
	resource.ResourceContent = string(decodedStr)
	return &resource, nil
}

func (r *ResourceHandler) deleteResource(uri string) error {
	req, err := http.NewRequest("DELETE", uri, nil)
	req.Header.Set("Content-Type", "application/json")
	if r.AuthHeader != "" && r.AuthToken != "" {
		req.Header.Set(r.AuthHeader, r.AuthToken)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
