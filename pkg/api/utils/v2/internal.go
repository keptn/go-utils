package v2

import (
	"context"
	"net/http"

	"github.com/keptn/go-utils/pkg/api/models"
)

// InternalAPISet is an implementation of KeptnInterface
// which can be used from within the Keptn control plane
type InternalAPISet struct {
	apimap                 InClusterAPIMappings
	httpClient             *http.Client
	apiHandler             *InternalAPIHandler
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

// InternalService is used to enumerate internal Keptn services
type InternalService int

const (
	ConfigurationService InternalService = iota
	ShipyardController
	ApiService
	SecretService
	MongoDBDatastore
)

// InClusterAPIMappings maps a keptn service name to its reachable domain name
type InClusterAPIMappings map[InternalService]string

// DefaultInClusterAPIMappings gives you the default InClusterAPIMappings
var DefaultInClusterAPIMappings = InClusterAPIMappings{
	ConfigurationService: "configuration-service:8080",
	ShipyardController:   "shipyard-controller:8080",
	ApiService:           "api-service:8080",
	SecretService:        "secret-service:8080",
	MongoDBDatastore:     "mongodb-datastore:8080",
}

// NewInternal creates a new InternalAPISet usable for calling Keptn services from within the control plane
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

	as := &InternalAPISet{}
	as.httpClient = client

	as.apiHandler = &InternalAPIHandler{
		shipyardControllerApiHandler: NewAPIHandlerWithHTTPClient(
			apimap[ShipyardController],
			&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))}),
	}

	as.authHandler = NewAuthHandlerWithHTTPClient(
		apimap[ApiService],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.logHandler = NewLogHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: getClientTransport(as.httpClient.Transport)})

	as.eventHandler = NewEventHandlerWithHTTPClient(
		apimap[MongoDBDatastore],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.projectHandler = NewProjectHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.resourceHandler = NewResourceHandlerWithHTTPClient(
		apimap[ConfigurationService],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.secretHandler = NewSecretHandlerWithHTTPClient(
		apimap[SecretService],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.sequenceControlHandler = NewSequenceControlHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.serviceHandler = NewServiceHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.shipyardControlHandler = NewShipyardControllerHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(getClientTransport(as.httpClient.Transport))})

	as.stageHandler = NewStageHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: wrapOtelTransport(as.httpClient.Transport)})

	as.uniformHandler = NewUniformHandlerWithHTTPClient(
		apimap[ShipyardController],
		&http.Client{Transport: getClientTransport(as.httpClient.Transport)})

	return as, nil
}

// API retrieves the APIHandler
func (c *InternalAPISet) API() APIInterface {
	return c.apiHandler
}

// Auth retrieves the AuthHandler
func (c *InternalAPISet) Auth() AuthInterface {
	return c.authHandler
}

// Events retrieves the EventHandler
func (c *InternalAPISet) Events() EventsInterface {
	return c.eventHandler
}

// Logs retrieves the LogHandler
func (c *InternalAPISet) Logs() LogsInterface {
	return c.logHandler
}

// Projects retrieves the ProjectHandler
func (c *InternalAPISet) Projects() ProjectsInterface {
	return c.projectHandler
}

// Resources retrieves the ResourceHandler
func (c *InternalAPISet) Resources() ResourcesInterface {
	return c.resourceHandler
}

// Secrets retrieves the SecretHandler
func (c *InternalAPISet) Secrets() SecretsInterface {
	return c.secretHandler
}

// Sequences retrieves the SequenceControlHandler
func (c *InternalAPISet) Sequences() SequencesInterface {
	return c.sequenceControlHandler
}

// Services retrieves the ServiceHandler
func (c *InternalAPISet) Services() ServicesInterface {
	return c.serviceHandler
}

// Stages retrieves the StageHandler
func (c *InternalAPISet) Stages() StagesInterface {
	return c.stageHandler
}

// Uniform retrieves the UniformHandler
func (c *InternalAPISet) Uniform() UniformInterface {
	return c.uniformHandler
}

// ShipyardControl retrieves the ShipyardControllerHandler
func (c *InternalAPISet) ShipyardControl() ShipyardControlInterface {
	return c.shipyardControlHandler
}

// InternalAPIHandler is used instead of APIHandler from go-utils because we cannot support
// (unauthenticated) internal calls to the api-service at the moment. So this implementation
// will panic as soon as a client wants to call these methods
type InternalAPIHandler struct {
	shipyardControllerApiHandler *APIHandler
}

func (i *InternalAPIHandler) SendEvent(_ context.Context, event models.KeptnContextExtendedCE, _ APISendEventOptions) (*models.EventContext, *models.Error) {
	panic("SendEvent() is not not supported for internal usage")
}

func (i *InternalAPIHandler) TriggerEvaluation(ctx context.Context, project string, stage string, service string, evaluation models.Evaluation, opts APITriggerEvaluationOptions) (*models.EventContext, *models.Error) {
	return i.shipyardControllerApiHandler.TriggerEvaluation(ctx, project, stage, service, evaluation, opts)
}

func (i *InternalAPIHandler) CreateProject(ctx context.Context, project models.CreateProject, opts APICreateProjectOptions) (string, *models.Error) {
	return i.shipyardControllerApiHandler.CreateProject(ctx, project, opts)
}

func (i *InternalAPIHandler) UpdateProject(ctx context.Context, project models.CreateProject, opts APIUpdateProjectOptions) (string, *models.Error) {
	return i.shipyardControllerApiHandler.UpdateProject(ctx, project, opts)
}

func (i *InternalAPIHandler) DeleteProject(ctx context.Context, project models.Project, opts APIDeleteProjectOptions) (*models.DeleteProjectResponse, *models.Error) {
	return i.shipyardControllerApiHandler.DeleteProject(ctx, project, opts)
}

func (i *InternalAPIHandler) CreateService(ctx context.Context, project string, service models.CreateService, opts APICreateServiceOptions) (string, *models.Error) {
	return i.shipyardControllerApiHandler.CreateService(ctx, project, service, opts)
}

func (i *InternalAPIHandler) DeleteService(ctx context.Context, project string, service string, opts APIDeleteServiceOptions) (*models.DeleteServiceResponse, *models.Error) {
	return i.shipyardControllerApiHandler.DeleteService(ctx, project, service, opts)
}

func (i *InternalAPIHandler) GetMetadata(_ context.Context, _ APIGetMetadataOptions) (*models.Metadata, *models.Error) {
	panic("GetMetadata() is not not supported for internal usage")
}
