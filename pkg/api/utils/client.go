package api

import (
	"fmt"
	"net/http"
	"net/url"
)

type KeptnInterface interface {
	APIV1() *APIHandler
	AuthV1() *AuthHandler
	EventsV1() *EventHandler
	LogsV1() *LogHandler
	ProjectsV1() *ProjectHandler
	ResourcesV1() *ResourceHandler
	SecretsV1() *SecretHandler
	SequencesV1() *SequenceControlHandler
	ServicesV1() *ServiceHandler
	StagesV1() *StageHandler
	UniformV1() *UniformHandler
}

// ApiSet contains the API utils for all keptn APIs
type ApiSet struct {
	endpointURL            *url.URL
	apiToken               string
	apiHandler             *APIHandler
	authHandler            *AuthHandler
	eventHandler           *EventHandler
	logHandler             *LogHandler
	projectHandler         *ProjectHandler
	resourceHandler        *ResourceHandler
	secretHandler          *SecretHandler
	sequenceControlHandler *SequenceControlHandler
	serviceHandler         *ServiceHandler
	stageHandler           *StageHandler
	uniformHandler         *UniformHandler
	shipyardControlHandler *ShipyardControllerHandler
}

// APIV1 retrieves the APIHandler
func (c *ApiSet) APIV1() *APIHandler {
	return c.apiHandler
}

// AuthV1 retrieves the AuthHandler
func (c *ApiSet) AuthV1() *AuthHandler {
	return c.authHandler
}

// EventsV1 retrieves the EventHandler
func (c *ApiSet) EventsV1() *EventHandler {
	return c.eventHandler
}

// LogsV1 retrieves the LogHandler
func (c *ApiSet) LogsV1() *LogHandler {
	return c.logHandler
}

// ProjectsV1 retrieves the ProjectHandler
func (c *ApiSet) ProjectsV1() *ProjectHandler {
	return c.projectHandler
}

// ResourcesV1 retrieves the ResourceHandler
func (c *ApiSet) ResourcesV1() *ResourceHandler {
	return c.resourceHandler
}

// SecretsV1 retrieves the SecretHandler
func (c *ApiSet) SecretsV1() *SecretHandler {
	return c.secretHandler
}

// SequencesV1 retrieves the SequenceControlHandler
func (c *ApiSet) SequencesV1() *SequenceControlHandler {
	return c.sequenceControlHandler
}

// ServicesV1 retrieves the ServiceHandler
func (c *ApiSet) ServicesV1() *ServiceHandler {
	return c.serviceHandler
}

// StagesV1 retrieves the StageHandler
func (c *ApiSet) StagesV1() *StageHandler {
	return c.stageHandler
}

// UniformV1 retrieves the UniformHandler
func (c *ApiSet) UniformV1() *UniformHandler {
	return c.uniformHandler
}

// ShipyardControlHandlerV1 retrieves the ShipyardControllerHandler
func (c *ApiSet) ShipyardControlHandlerV1() *ShipyardControllerHandler {
	return c.shipyardControlHandler
}

// Token retrieves the API token
func (c *ApiSet) Token() string {
	return c.apiToken
}

// Endpoint retrieves the base API endpoint URL
func (c *ApiSet) Endpoint() *url.URL {
	return c.endpointURL
}

// NewApiSet creates a new ApiSet
func NewApiSet(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) (*ApiSet, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create apiset: %w", err)
	}
	var as ApiSet
	as.endpointURL = u
	as.apiToken = authToken
	as.apiHandler = NewAuthenticatedAPIHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.authHandler = NewAuthenticatedAuthHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.logHandler = NewAuthenticatedLogHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.authHandler = NewAuthenticatedAuthHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.eventHandler = NewAuthenticatedEventHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.projectHandler = NewAuthenticatedProjectHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.resourceHandler = NewAuthenticatedResourceHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.secretHandler = NewAuthenticatedSecretHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.sequenceControlHandler = NewAuthenticatedSequenceControlHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.serviceHandler = NewAuthenticatedServiceHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.shipyardControlHandler = NewAuthenticatedShipyardControllerHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.stageHandler = NewAuthenticatedStageHandler(baseURL, authToken, authHeader, httpClient, scheme)
	as.uniformHandler = NewAuthenticatedUniformHandler(baseURL, authToken, authHeader, httpClient, scheme)
	return &as, nil
}
