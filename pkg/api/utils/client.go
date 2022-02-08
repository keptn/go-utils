package api

import (
	"fmt"
	"net/http"
	"net/url"
)

var _ KeptnInterface = (*APISet)(nil)

type KeptnInterface interface {
	APIV1() APIV1Interface
	AuthV1() AuthV1Interface
	EventsV1() EventsV1Interface
	LogsV1() LogsV1Interface
	ProjectsV1() ProjectsV1Interface
	ResourcesV1() ResourcesV1Interface
	SecretsV1() SecretsV1Interface
	SequencesV1() SequencesV1Interface
	ServicesV1() ServicesV1Interface
	StagesV1() StagesV1Interface
	UniformV1() UniformV1Interface
	ShipyardControlV1() ShipyardControlV1Interface
}

// APISet contains the API utils for all keptn APIs
type APISet struct {
	endpointURL            *url.URL
	apiToken               string
	authHeader             string
	scheme                 string
	internal               bool
	httpClient             *http.Client
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
func (c *APISet) APIV1() APIV1Interface {
	return c.apiHandler
}

// AuthV1 retrieves the AuthHandler
func (c *APISet) AuthV1() AuthV1Interface {
	return c.authHandler
}

// EventsV1 retrieves the EventHandler
func (c *APISet) EventsV1() EventsV1Interface {
	return c.eventHandler
}

// LogsV1 retrieves the LogHandler
func (c *APISet) LogsV1() LogsV1Interface {
	return c.logHandler
}

// ProjectsV1 retrieves the ProjectHandler
func (c *APISet) ProjectsV1() ProjectsV1Interface {
	return c.projectHandler
}

// ResourcesV1 retrieves the ResourceHandler
func (c *APISet) ResourcesV1() ResourcesV1Interface {
	return c.resourceHandler
}

// SecretsV1 retrieves the SecretHandler
func (c *APISet) SecretsV1() SecretsV1Interface {
	return c.secretHandler
}

// SequencesV1 retrieves the SequenceControlHandler
func (c *APISet) SequencesV1() SequencesV1Interface {
	return c.sequenceControlHandler
}

// ServicesV1 retrieves the ServiceHandler
func (c *APISet) ServicesV1() ServicesV1Interface {
	return c.serviceHandler
}

// StagesV1 retrieves the StageHandler
func (c *APISet) StagesV1() StagesV1Interface {
	return c.stageHandler
}

// UniformV1 retrieves the UniformHandler
func (c *APISet) UniformV1() UniformV1Interface {
	return c.uniformHandler
}

// ShipyardControlV1 retrieves the ShipyardControllerHandler
func (c *APISet) ShipyardControlV1() ShipyardControlV1Interface {
	return c.shipyardControlHandler
}

// Token retrieves the API token
func (c *APISet) Token() string {
	return c.apiToken
}

// Endpoint retrieves the base API endpoint URL
func (c *APISet) Endpoint() *url.URL {
	return c.endpointURL
}

// WithAuthToken sets the given auth token.
// Optionally a custom auth header can be set (default x-token)
func WithAuthToken(authToken string, authHeader ...string) func(*APISet) {
	aHeader := "x-token"
	if len(authHeader) > 0 {
		aHeader = authHeader[0]
	}
	return func(a *APISet) {
		a.apiToken = authToken
		a.authHeader = aHeader
	}
}

// WithHTTPClient configures a custom http client to use
func WithHTTPClient(client *http.Client) func(*APISet) {
	return func(a *APISet) {
		a.httpClient = client
	}
}

// WithScheme sets the scheme
// If this option is not used, then default scheme "http" is used by the APISet
func WithScheme(scheme string) func(*APISet) {
	return func(a *APISet) {
		a.scheme = scheme
	}
}

// Internal configures the APISet to be used
// internally within the control plane
func Internal() func(*APISet) {
	return func(a *APISet) {
		a.internal = true
	}
}

// New creates a new APISet instance
func New(baseURL string, options ...func(*APISet)) (*APISet, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create apiset: %w", err)
	}
	as := &APISet{}
	for _, o := range options {
		o(as)
	}
	as.endpointURL = u
	as.httpClient = createInstrumentedClientTransport(as.httpClient)

	if as.scheme == "" {
		if as.endpointURL.Scheme != "" {
			as.scheme = u.Scheme
		} else {
			as.scheme = "http"
		}
	}

	as.apiHandler = createAuthenticatedAPIHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.authHandler = createAuthenticatedAuthHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.logHandler = createAuthenticatedLogHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.eventHandler = createAuthenticatedEventHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.projectHandler = createAuthProjectHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.resourceHandler = createAuthenticatedResourceHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.secretHandler = createAuthenticatedSecretHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.sequenceControlHandler = createAuthenticatedSequenceControlHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.serviceHandler = createAuthenticatedServiceHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.shipyardControlHandler = createAuthenticatedShipyardControllerHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.stageHandler = createAuthenticatedStageHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	as.uniformHandler = createAuthenticatedUniformHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme, as.internal)
	return as, nil
}

type InternalAPISet struct {
	httpClient             *http.Client
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

type InternalService int

const (
	ConfigurationService InternalService = iota
	ShipyardController
	ApiService
	SecretService
	Unknown
)

type InClusterAPIMappings map[InternalService]string

func NewInternal(apiMappings InClusterAPIMappings, client *http.Client) (*InternalAPISet, error) {
	if client == nil {
		client = &http.Client{}
	}
	as := &InternalAPISet{}

	as.httpClient = client
	as.apiHandler = createAuthenticatedAPIHandler(apiMappings[Unknown], "", "", as.httpClient, "http", true)
	as.authHandler = createAuthenticatedAuthHandler(apiMappings[ApiService], "", "", as.httpClient, "http", true)
	as.logHandler = createAuthenticatedLogHandler(apiMappings[ShipyardController], "", "", as.httpClient, "http", true)
	as.eventHandler = createAuthenticatedEventHandler(apiMappings[ApiService], "", "", as.httpClient, "http", true)
	as.projectHandler = createAuthProjectHandler(apiMappings[ShipyardController], "", "", as.httpClient, "http", true)
	as.resourceHandler = createAuthenticatedResourceHandler(apiMappings[ConfigurationService], "", "", as.httpClient, "http", true)
	as.secretHandler = createAuthenticatedSecretHandler(apiMappings[SecretService], "", "", as.httpClient, "http", true)
	as.sequenceControlHandler = createAuthenticatedSequenceControlHandler(apiMappings[ShipyardController], "", "", as.httpClient, "http", true)
	as.serviceHandler = createAuthenticatedServiceHandler(apiMappings[ShipyardController], "", "", as.httpClient, "http", true)
	as.shipyardControlHandler = createAuthenticatedShipyardControllerHandler(apiMappings[ShipyardController], "", "", as.httpClient, "http", true)
	as.stageHandler = createAuthenticatedStageHandler(apiMappings[ShipyardController], "", "", as.httpClient, "http", true)
	as.uniformHandler = createAuthenticatedUniformHandler(apiMappings[ShipyardController], "", "", as.httpClient, "http", true)

	return as, nil
}

// APIV1 retrieves the APIHandler
func (c *InternalAPISet) APIV1() APIV1Interface {
	return c.apiHandler
}

// AuthV1 retrieves the AuthHandler
func (c *InternalAPISet) AuthV1() AuthV1Interface {
	return c.authHandler
}

// EventsV1 retrieves the EventHandler
func (c *InternalAPISet) EventsV1() EventsV1Interface {
	return c.eventHandler
}

// LogsV1 retrieves the LogHandler
func (c *InternalAPISet) LogsV1() LogsV1Interface {
	return c.logHandler
}

// ProjectsV1 retrieves the ProjectHandler
func (c *InternalAPISet) ProjectsV1() ProjectsV1Interface {
	return c.projectHandler
}

// ResourcesV1 retrieves the ResourceHandler
func (c *InternalAPISet) ResourcesV1() ResourcesV1Interface {
	return c.resourceHandler
}

// SecretsV1 retrieves the SecretHandler
func (c *InternalAPISet) SecretsV1() SecretsV1Interface {
	return c.secretHandler
}

// SequencesV1 retrieves the SequenceControlHandler
func (c *InternalAPISet) SequencesV1() SequencesV1Interface {
	return c.sequenceControlHandler
}

// ServicesV1 retrieves the ServiceHandler
func (c *InternalAPISet) ServicesV1() ServicesV1Interface {
	return c.serviceHandler
}

// StagesV1 retrieves the StageHandler
func (c *InternalAPISet) StagesV1() StagesV1Interface {
	return c.stageHandler
}

// UniformV1 retrieves the UniformHandler
func (c *InternalAPISet) UniformV1() UniformV1Interface {
	return c.uniformHandler
}

// ShipyardControlV1 retrieves the ShipyardControllerHandler
func (c *InternalAPISet) ShipyardControlV1() ShipyardControlV1Interface {
	return c.shipyardControlHandler
}
