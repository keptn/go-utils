package v2

import (
	"fmt"
	"net/http"
	"net/url"
)

var _ KeptnInterface = (*APISet)(nil)

type KeptnInterface interface {
	API() APIInterface
	Auth() AuthInterface
	Events() EventsInterface
	Logs() LogsInterface
	Projects() ProjectsInterface
	Resources() ResourcesInterface
	Secrets() SecretsInterface
	Sequences() SequencesInterface
	Services() ServicesInterface
	Stages() StagesInterface
	Uniform() UniformInterface
	ShipyardControl() ShipyardControlInterface
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

// API retrieves the APIHandler
func (c *APISet) API() APIInterface {
	return c.apiHandler
}

// Auth retrieves the AuthHandler
func (c *APISet) Auth() AuthInterface {
	return c.authHandler
}

// Events retrieves the EventHandler
func (c *APISet) Events() EventsInterface {
	return c.eventHandler
}

// Logs retrieves the LogHandler
func (c *APISet) Logs() LogsInterface {
	return c.logHandler
}

// Projects retrieves the ProjectHandler
func (c *APISet) Projects() ProjectsInterface {
	return c.projectHandler
}

// Resources retrieves the ResourceHandler
func (c *APISet) Resources() ResourcesInterface {
	return c.resourceHandler
}

// Secrets retrieves the SecretHandler
func (c *APISet) Secrets() SecretsInterface {
	return c.secretHandler
}

// Sequences retrieves the SequenceControlHandler
func (c *APISet) Sequences() SequencesInterface {
	return c.sequenceControlHandler
}

// Services retrieves the ServiceHandler
func (c *APISet) Services() ServicesInterface {
	return c.serviceHandler
}

// Stages retrieves the StageHandler
func (c *APISet) Stages() StagesInterface {
	return c.stageHandler
}

// Uniform retrieves the UniformHandler
func (c *APISet) Uniform() UniformInterface {
	return c.uniformHandler
}

// ShipyardControl retrieves the ShipyardControllerHandler
func (c *APISet) ShipyardControl() ShipyardControlInterface {
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

	as.apiHandler = NewAuthenticatedAPIHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.authHandler = NewAuthenticatedAuthHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.logHandler = NewAuthenticatedLogHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.eventHandler = NewAuthenticatedEventHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.projectHandler = NewAuthenticatedProjectHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.resourceHandler = NewAuthenticatedResourceHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.secretHandler = NewAuthenticatedSecretHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.sequenceControlHandler = NewAuthenticatedSequenceControlHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.serviceHandler = NewAuthenticatedServiceHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.shipyardControlHandler = NewAuthenticatedShipyardControllerHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.stageHandler = NewAuthenticatedStageHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	as.uniformHandler = NewAuthenticatedUniformHandler(baseURL, as.apiToken, as.authHeader, as.httpClient, as.scheme)
	return as, nil
}
