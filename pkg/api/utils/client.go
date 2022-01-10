package api

import (
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

func (c *ApiSet) APIV1() *APIHandler {
	return c.apiHandler
}

func (c *ApiSet) AuthV1() *AuthHandler {
	return c.authHandler
}

func (c *ApiSet) EventsV1() *EventHandler {
	return c.eventHandler
}

func (c *ApiSet) LogsV1() *LogHandler {
	return c.logHandler
}

func (c *ApiSet) ProjectsV1() *ProjectHandler {
	return c.projectHandler
}

func (c *ApiSet) ResourcesV1() *ResourceHandler {
	return c.resourceHandler
}

func (c *ApiSet) SecretsV1() *SecretHandler {
	return c.secretHandler
}

func (c *ApiSet) SequencesV1() *SequenceControlHandler {
	return c.sequenceControlHandler
}

func (c *ApiSet) ServicesV1() *ServiceHandler {
	return c.serviceHandler
}

func (c *ApiSet) StagesV1() *StageHandler {
	return c.stageHandler
}

func (c *ApiSet) UniformV1() *UniformHandler {
	return c.uniformHandler
}

func (c *ApiSet) ShipyardControlHandlerV1() *ShipyardControllerHandler {
	return c.shipyardControlHandler
}

func (c *ApiSet) Token() string {
	return c.apiToken
}

func (c *ApiSet) Endpoint() *url.URL {
	return c.endpointURL
}

func NewApiSet(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) (*ApiSet, error) {
	var as ApiSet
	as.apiHandler = NewAuthenticatedAPIHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.authHandler = NewAuthenticatedAuthHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.logHandler = NewAuthenticatedLogHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.authHandler = NewAuthenticatedAuthHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.eventHandler = NewAuthenticatedEventHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.projectHandler = NewAuthenticatedProjectHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.resourceHandler = NewAuthenticatedResourceHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.secretHandler = NewAuthenticatedSecretHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.sequenceControlHandler = NewAuthenticatedSequenceControlHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.serviceHandler = NewAuthenticatedServiceHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.shipyardControlHandler = NewAuthenticatedShipyardControllerHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.stageHandler = NewAuthenticatedStageHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.uniformHandler = NewAuthenticatedUniformHandler(baseURL, authToken, "x-token", httpClient, scheme)
	as.apiToken = authToken
	url, _ := url.Parse(baseURL)
	as.endpointURL = url
	return &as, nil
}
