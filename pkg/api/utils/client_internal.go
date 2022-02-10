package api

import "net/http"

type InternalAPISet struct {
	*APISet
	apimap InClusterAPIMappings
}

type InternalService int

const (
	ConfigurationService InternalService = iota
	ShipyardController
	ApiService
	SecretService
	MongoDBDatastore
)

type InClusterAPIMappings map[InternalService]string

var DefaultInClusterAPIMappings = InClusterAPIMappings{
	ConfigurationService: "configuration-service:8080",
	ShipyardController:   "shipyard-controller:8080",
	ApiService:           "api-service:8080",
	SecretService:        "secret-service:8080",
	MongoDBDatastore:     "mongodb-datastore:8080",
}

// NewInternal creates a new InternalAPISet usable for calling keptn services from within the control plane
func NewInternal(client *http.Client, apiMappings ...InClusterAPIMappings) (*InternalAPISet, error) {
	var apimap InClusterAPIMappings
	if len(apiMappings) > 0 {
		apimap = apiMappings[0]
	} else {
		apimap = DefaultInClusterAPIMappings
	}

	if client == nil {
		client = &http.Client{}
	}

	as := &InternalAPISet{APISet: &APISet{}}
	as.internal = true
	as.httpClient = client

	as.apiHandler = createAuthenticatedAPIHandler(apimap[ShipyardController], "", "", as.httpClient, "http", as.internal)
	as.authHandler = createAuthenticatedAuthHandler(apimap[ApiService], "", "", as.httpClient, "http", as.internal)
	as.logHandler = createAuthenticatedLogHandler(apimap[ShipyardController], "", "", as.httpClient, "http", as.internal)
	as.eventHandler = createAuthenticatedEventHandler(apimap[MongoDBDatastore], "", "", as.httpClient, "http", as.internal)
	as.projectHandler = createAuthProjectHandler(apimap[ShipyardController], "", "", as.httpClient, "http", as.internal)
	as.resourceHandler = createAuthenticatedResourceHandler(apimap[ConfigurationService], "", "", as.httpClient, "http", as.internal)
	as.secretHandler = createAuthenticatedSecretHandler(apimap[SecretService], "", "", as.httpClient, "http", as.internal)
	as.sequenceControlHandler = createAuthenticatedSequenceControlHandler(apimap[ShipyardController], "", "", as.httpClient, "http", as.internal)
	as.serviceHandler = createAuthenticatedServiceHandler(apimap[ShipyardController], "", "", as.httpClient, "http", as.internal)
	as.shipyardControlHandler = createAuthenticatedShipyardControllerHandler(apimap[ShipyardController], "", "", as.httpClient, "http", as.internal)
	as.stageHandler = createAuthenticatedStageHandler(apimap[ShipyardController], "", "", as.httpClient, "http", as.internal)
	as.uniformHandler = createAuthenticatedUniformHandler(apimap[ShipyardController], "", "", as.httpClient, "http", as.internal)
	return as, nil
}
