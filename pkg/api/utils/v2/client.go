package v2

import (
	"fmt"
	"net/http"
	"net/url"
)

var _ KeptnInterface = (*APISet)(nil)

type KeptnInterface interface {
	APIV1() APIInterface
	AuthV1() AuthInterface
	EventsV1() EventsInterface
	LogsV1() LogsInterface
	ProjectsV1() ProjectsInterface
	ResourcesV1() ResourcesInterface
	SecretsV1() SecretsInterface
	SequencesV1() SequencesInterface
	ServicesV1() ServicesInterface
	StagesV1() StagesInterface
	UniformV1() UniformInterface
	ShipyardControlV1() ShipyardControlInterface
}

// APISet contains the API utils for all Keptn APIs
type APISet struct {
	endpointURL            *url.URL
	apiToken               string
	authHeader             string
	scheme                 string
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
func (c *APISet) APIV1() APIInterface {
	return c.apiHandler
}

// AuthV1 retrieves the AuthHandler
func (c *APISet) AuthV1() AuthInterface {
	return c.authHandler
}

// EventsV1 retrieves the EventHandler
func (c *APISet) EventsV1() EventsInterface {
	return c.eventHandler
}

// LogsV1 retrieves the LogHandler
func (c *APISet) LogsV1() LogsInterface {
	return c.logHandler
}

// ProjectsV1 retrieves the ProjectHandler
func (c *APISet) ProjectsV1() ProjectsInterface {
	return c.projectHandler
}

// ResourcesV1 retrieves the ResourceHandler
func (c *APISet) ResourcesV1() ResourcesInterface {
	return c.resourceHandler
}

// SecretsV1 retrieves the SecretHandler
func (c *APISet) SecretsV1() SecretsInterface {
	return c.secretHandler
}

// SequencesV1 retrieves the SequenceControlHandler
func (c *APISet) SequencesV1() SequencesInterface {
	return c.sequenceControlHandler
}

// ServicesV1 retrieves the ServiceHandler
func (c *APISet) ServicesV1() ServicesInterface {
	return c.serviceHandler
}

// StagesV1 retrieves the StageHandler
func (c *APISet) StagesV1() StagesInterface {
	return c.stageHandler
}

// UniformV1 retrieves the UniformHandler
func (c *APISet) UniformV1() UniformInterface {
	return c.uniformHandler
}

// ShipyardControlV1 retrieves the ShipyardControllerHandler
func (c *APISet) ShipyardControlV1() ShipyardControlInterface {
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

// New creates a new APISet instance
func New(baseURL string, options ...func(*APISet)) (*APISet, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create apiset: %w", err)
	}
	as := &APISet{}
	for _, o := range options {
		if o != nil {
			o(as)
		}
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

	as.apiHandler = createAuthenticatedAPIHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.authHandler = createAuthenticatedAuthHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.logHandler = createAuthenticatedLogHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.eventHandler = createAuthenticatedEventHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.projectHandler = createAuthProjectHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.resourceHandler = createAuthenticatedResourceHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.secretHandler = createAuthenticatedSecretHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.sequenceControlHandler = createAuthenticatedSequenceControlHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.serviceHandler = createAuthenticatedServiceHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.shipyardControlHandler = createAuthenticatedShipyardControllerHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.stageHandler = createAuthenticatedStageHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.uniformHandler = createAuthenticatedUniformHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	return as, nil
}
