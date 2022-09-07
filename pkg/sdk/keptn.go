package sdk

import (
	"context"
	"github.com/keptn/go-utils/pkg/sdk/connector/types"
	sdk "github.com/keptn/go-utils/pkg/sdk/internal/api"
	"github.com/keptn/go-utils/pkg/sdk/internal/config"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	apiv2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk/connector/controlplane"
)

const (
	shkeptnspecversion = "0.2.4"
	cloudeventsversion = "1.0"
)

type IKeptn interface {
	// Start starts the internal event handling logic and needs to be called by the user
	// after creating value of IKeptn
	Start() error
	// GetResourceHandler returns a handler to fetch data from the configuration service
	GetResourceHandler() ResourceHandler
	// SendStartedEvent sends a started event for the given input event to the Keptn API
	// The first parameter is the "parent" event from which common information like e.g.
	// the keptn context, task name, ... are taken and used for constructing the corresponding .started event
	SendStartedEvent(KeptnEvent) error
	// SendFinishedEvent sends a finished event for the given input event to the Keptn API.
	// The first parameter can be seen as the "parent" event from which common information like e.g.
	// the keptn context, task name, ... are taken and used for constructing the corresponding .finished event.
	// The second parameter is the new event data to be set on the newly constructed .finished event
	SendFinishedEvent(KeptnEvent, interface{}) error
	// Logger returns the logger used by the sdk
	// Per default DefaultLogger is used which internally just uses the go logging package
	// Another logger can be configured using the sdk.WithLogger function
	Logger() Logger
	// APIV1 returns API utils for all Keptn APIs
	APIV1() api.KeptnInterface
	// APIV2 returns API utils for all v2 Keptn APIs
	APIV2() apiv2.KeptnInterface
}

type TaskHandler interface {
	// Execute is called whenever the actual business-logic of the service shall be executed.
	// Thus, the core logic of the service shall be triggered/implemented in this method.
	//
	// Note, that the contract of the method is to return the payload of the .finished event to be sent out as well as a Error Pointer
	// or nil, if there was no error during execution.
	Execute(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error)
}

type KeptnEvent models.KeptnContextExtendedCE

type Error struct {
	StatusType keptnv2.StatusType
	ResultType keptnv2.ResultType
	Message    string
	Err        error
}

func (e Error) Error() string {
	return e.Message
}

// KeptnOption can be used to configure the keptn sdk
type KeptnOption func(*Keptn)

type ResourceHandler interface {
	GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error)
}

type resourceHandlerWrapper struct {
	resourceHandler apiv2.ResourcesInterface
}

func (rhw *resourceHandlerWrapper) GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
	v2Scope := apiv2.NewResourceScope().Project(scope.GetProject()).Stage(scope.GetStage()).Service(scope.GetService()).Resource(scope.GetResource())

	return rhw.resourceHandler.GetResource(context.Background(), *v2Scope, apiv2.ResourcesGetResourceOptions{})
}

type healthEndpointRunner func(port string, cp *controlplane.ControlPlane)

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var gracefulShutdownKey = gracefulShutdownKeyType{}

type wgInterface interface {
	Add(delta int)
	Done()
	Wait()
}

type nopWG struct {
	// --
}

func (w *nopWG) Add(delta int) {
	// --
}
func (w *nopWG) Done() {
	// --
}
func (w *nopWG) Wait() {
	// --
}

// WithTaskHandler registers a handler which is responsible for processing a .triggered event.
// Note, that if you want to have more control on configuring the behavior of the task handler,
// you can use WithTaskEventHandler instead
func WithTaskHandler(eventType string, handler TaskHandler, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) KeptnOption {
	return WithTaskEventHandler(eventType, handler, TaskHandlerOptions{
		Filters:               filters,
		SkipAutomaticResponse: false,
	})
}

// WithTaskEventHandler registers a handler which is responsible for processing a received .triggered event
func WithTaskEventHandler(eventType string, handler TaskHandler, options TaskHandlerOptions) KeptnOption {
	return func(k *Keptn) {
		k.taskRegistry.Add(eventType, taskEntry{taskHandler: handler, eventFilters: options.Filters, taskHandlerOpts: options})
	}
}

// WithAutomaticResponse sets the option to instruct the sdk to automatically send a .started and .finished event.
// Per default this behavior is turned on and can be disabled with this function. Note, that this affects ALL
// task handlers. If you want to disable automatic event responses for a specific task handler, this can be done
// with the respective TaskHandlerOptions passed to WithTaskEventHandler
func WithAutomaticResponse(autoResponse bool) KeptnOption {
	return func(k *Keptn) {
		k.automaticEventResponse = autoResponse
	}
}

// TaskHandlerOptions are specific options for a task handler
type TaskHandlerOptions struct {
	// Filters specifies functions that determine whether the event shall be handled or ignored
	Filters []func(IKeptn, KeptnEvent) bool
	// SkipAutomaticResponse determines whether automatic sending of .started/.finished events should be skipped
	SkipAutomaticResponse bool
}

// WithGracefulShutdown sets the option to ensure running tasks/handlers will finish in case of interrupt or forced termination
// Per default this behavior is turned on and can be disabled with this function
func WithGracefulShutdown(gracefulShutdown bool) KeptnOption {
	return func(k *Keptn) {
		k.gracefulShutdown = gracefulShutdown
	}
}

// WithLogger configures keptn to use another logger
func WithLogger(logger Logger) KeptnOption {
	return func(k *Keptn) {
		k.logger = logger
	}
}

// Keptn is the default implementation of IKeptn
type Keptn struct {
	controlPlane           *controlplane.ControlPlane
	eventSender            controlplane.EventSender
	resourceHandler        ResourceHandler
	api                    api.KeptnInterface
	apiV2                  apiv2.KeptnInterface
	source                 string
	taskRegistry           *taskRegistry
	syncProcessing         bool
	automaticEventResponse bool
	gracefulShutdown       bool
	logger                 Logger
	env                    config.EnvConfig
	healthEndpointRunner   healthEndpointRunner
}

// NewKeptn creates a new Keptn
func NewKeptn(source string, opts ...KeptnOption) *Keptn {
	keptn := &Keptn{
		source:                 source,
		taskRegistry:           newTaskMap(),
		automaticEventResponse: true,
		gracefulShutdown:       true,
		syncProcessing:         false,
		logger:                 newDefaultLogger(),
		env:                    config.NewEnvConfig(),
		healthEndpointRunner:   newHealthEndpointRunner,
	}

	for _, opt := range opts {
		opt(keptn)
	}

	var env config.EnvConfig
	if err := envconfig.Process("", &env); err != nil {
		keptn.logger.Fatalf("failed to process env vars: %v", err)
	}

	httpClientFactory := sdk.CreateClientGetter(env)
	initializationResult, err := sdk.Initialize(env, httpClientFactory, keptn.logger)
	if err != nil {
		keptn.logger.Fatalf("failed to initialize keptn sdk: %v", err)
	}

	keptn.api = initializationResult.KeptnAPI
	keptn.apiV2 = initializationResult.KeptnAPIV2
	keptn.controlPlane = initializationResult.ControlPlane
	keptn.eventSender = initializationResult.EventSenderCallback
	keptn.resourceHandler = &resourceHandlerWrapper{resourceHandler: initializationResult.KeptnAPIV2.Resources()}

	return keptn
}

func (k *Keptn) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	k.logger.Debug("Handling event ", event)
	eventSender, ok := ctx.Value(types.EventSenderKey).(controlplane.EventSender)
	if !ok {
		k.logger.Errorf("Unable to get event sender. Skip processing of event %s", event.ID)
		return nil
	}

	if event.Type == nil {
		k.logger.Errorf("Unable to get event type. Skip processing of event %s", event.ID)
		return nil
	}

	if !keptnv2.IsTaskEventType(*event.Type) {
		k.logger.Errorf("Event type %s does not match format for task events. Skip Processing of event %s", *event.Type, event.ID)
		return nil
	}
	wg, ok := ctx.Value(gracefulShutdownKey).(wgInterface)
	if !ok {
		k.logger.Errorf("Unable to get graceful shutdown wait group. Skip processing of event %s", event.ID)
		return nil
	}
	wg.Add(1)
	k.runEventTaskAction(func() {
		{
			defer wg.Done()
			if handler, ok := k.taskRegistry.Contains(*event.Type); ok {
				keptnEvent := &KeptnEvent{}
				if err := keptnv2.Decode(&event, keptnEvent); err != nil {
					errorLogEvent, err := createErrorLogEvent(k.source, event, nil, &Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed})
					if err != nil {
						k.logger.Errorf("Unable to create '.error.log' event from '.triggered' event: %v", err)
						return
					}
					// no started event sent yet, so it only makes sense to Send an error log event at this point
					if err := eventSender(*errorLogEvent); err != nil {
						k.logger.Errorf("Unable to send '.finished' event: %v", err)
						return
					}
					return
				}

				// execute the filtering functions of the task handler to determine whether the incoming event should be handled
				// only if all functions return true, the event will be handled
				for _, filterFn := range handler.eventFilters {
					if !filterFn(k, *keptnEvent) {
						k.logger.Infof("Will not handle incoming %s event", *event.Type)
						return
					}
				}

				// automatic response of events is enabled if it is turned on globally, and not disabled for the specific handler
				autoResponse := k.automaticEventResponse && !k.taskRegistry.Get(*event.Type).taskHandlerOpts.SkipAutomaticResponse

				// only respond with .started event if the incoming event is a task.triggered event
				if keptnv2.IsTaskEventType(*event.Type) && keptnv2.IsTriggeredEventType(*event.Type) && autoResponse {
					startedEvent, err := createStartedEvent(k.source, event)
					if err != nil {
						k.logger.Errorf("Unable to create '.started' event from '.triggered' event: %v", err)
						return
					}
					if err := eventSender(*startedEvent); err != nil {
						k.logger.Errorf("Unable to send '.started' event: %v", err)
						return
					}
				}

				result, err := handler.taskHandler.Execute(k, *keptnEvent)
				if err != nil {
					k.logger.Errorf("Error during task execution %v", err.Err)
					if autoResponse {
						errorEvent, err := createErrorEvent(k.source, event, result, err)
						if err != nil {
							k.logger.Errorf("Unable to create '.error' event: %v", err)
							return
						}
						if err := eventSender(*errorEvent); err != nil {
							k.logger.Errorf("Unable to send '.error' event: %v", err)
							return
						}
					}
					return
				}
				if result == nil {
					k.logger.Infof("no finished data set by task executor for event %s. Skipping sending finished event", *event.Type)
				} else if keptnv2.IsTaskEventType(*event.Type) && keptnv2.IsTriggeredEventType(*event.Type) && autoResponse {
					finishedEvent, err := createFinishedEvent(k.source, event, result)
					if err != nil {
						k.logger.Errorf("Unable to create '.finished' event: %v", err)
						return
					}
					if err := eventSender(*finishedEvent); err != nil {
						k.logger.Errorf("Unable to send '.finished' event: %v", err)
						return
					}
				}
			}
		}
	})
	return nil
}

func (k *Keptn) RegistrationData() controlplane.RegistrationData {
	subscriptions := []models.EventSubscription{}
	subjects := []string{}
	if k.env.PubSubTopic != "" {
		subjects = strings.Split(k.env.PubSubTopic, ",")
	}

	for _, s := range subjects {
		subscriptions = append(subscriptions, models.EventSubscription{Event: s})
	}
	return controlplane.RegistrationData{
		Name: k.source,
		MetaData: models.MetaData{
			Hostname:           k.env.K8sNodeName,
			IntegrationVersion: k.env.K8sDeploymentVersion,
			Location:           k.env.Location,
			DistributorVersion: "0.15.0", // note: to be deleted when bridge stops requiring this info
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      k.env.K8sNamespace,
				PodName:        k.env.K8sPodName,
				DeploymentName: k.env.K8sDeploymentName,
			},
		},
		Subscriptions: subscriptions,
	}
}

func (k *Keptn) Start() error {
	if k.env.HealthEndpointEnabled {
		k.healthEndpointRunner(k.env.HealthEndpointPort, k.controlPlane)
	}
	ctx, wg := k.getContext(k.gracefulShutdown)
	err := k.controlPlane.Register(ctx, k)
	// add additional waiting time to ensure the waitGroup has been increased for all events that have been received between receiving SIGTERM and this point
	<-time.After(5 * time.Second)
	wg.Wait()

	return err
}

func (k *Keptn) GetResourceHandler() ResourceHandler {
	return k.resourceHandler
}

func (k *Keptn) SendStartedEvent(parentEvent KeptnEvent) error {
	startedEvent, err := createStartedEvent(k.source, models.KeptnContextExtendedCE(parentEvent))
	if err != nil {
		return err
	}
	return k.eventSender(*startedEvent)
}

func (k *Keptn) SendFinishedEvent(parentEvent KeptnEvent, newEventData interface{}) error {
	finishedEvent, err := createFinishedEvent(k.source, models.KeptnContextExtendedCE(parentEvent), newEventData)
	if err != nil {
		return err
	}
	return k.eventSender(*finishedEvent)
}

// APIV1 retrieves the APIV1 client
// Deprecated use APIV2 instead
func (k *Keptn) APIV1() api.KeptnInterface {
	return k.api
}

// APIV2 retrieves the APIV2 client
func (k *Keptn) APIV2() apiv2.KeptnInterface {
	return k.apiV2
}

func (k *Keptn) Logger() Logger {
	return k.logger
}

func (k *Keptn) runEventTaskAction(fn func()) {
	if k.syncProcessing {
		fn()
	} else {
		go fn()
	}
}

func (k *Keptn) getContext(graceful bool) (context.Context, wgInterface) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	var wg wgInterface
	if graceful {
		wg = &sync.WaitGroup{}
	} else {
		wg = &nopWG{}
	}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), gracefulShutdownKey, wg))
	go func() {
		<-ch
		cancel()
	}()
	return ctx, wg
}

func noOpHealthEndpointRunner(port string, cp *controlplane.ControlPlane) {}

func newHealthEndpointRunner(port string, cp *controlplane.ControlPlane) {
	go func() {
		api.RunHealthEndpoint(port, api.WithReadinessConditionFunc(func() bool {
			return cp.IsRegistered()
		}))
	}()
}
