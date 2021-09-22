package v0_3_0

import (
	"context"
	"errors"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/keptn/go-utils/config"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const keptnStartedEventSuffix = ".started"
const keptnStatusChangedEventSuffix = ".status.changed"
const keptnFinishedEventSuffix = ".finished"

const keptnContextCEExtension = "shkeptncontext"
const keptnSpecVersionCEExtension = "shkeptnspecversion"
const triggeredIDCEExtension = "triggeredid"

type Keptn struct {
	keptn.KeptnBase
}

func NewKeptn(incomingEvent *cloudevents.Event, opts keptn.KeptnOpts) (*Keptn, error) {
	extension, _ := incomingEvent.Context.GetExtension("shkeptncontext")
	shkeptncontext := extension.(string)

	// create a base KeptnBase Event
	keptnBase := &keptnv2.EventData{}

	if err := incomingEvent.DataAs(keptnBase); err != nil {
		return nil, err
	}

	k := &Keptn{
		KeptnBase: keptn.KeptnBase{
			Event:              keptnBase,
			CloudEvent:         incomingEvent,
			KeptnContext:       shkeptncontext,
			UseLocalFileSystem: opts.UseLocalFileSystem,
			ResourceHandler:    nil,
		}}

	csURL := keptn.ConfigurationServiceURL
	if opts.ConfigurationServiceURL != "" {
		csURL = opts.ConfigurationServiceURL
	}

	if opts.EventBrokerURL != "" && opts.EventSender == nil {
		httpSender, err := keptnv2.NewHTTPEventSender(opts.EventBrokerURL)
		if err != nil {
			return nil, fmt.Errorf("could not initialize Keptn Handler: %s", err.Error())
		}
		k.KeptnBase.EventSender = httpSender
	} else if opts.EventSender != nil {
		k.KeptnBase.EventSender = opts.EventSender
	} else {
		httpSender, err := keptnv2.NewHTTPEventSender(keptnv2.DefaultHTTPEventEndpoint)
		if err != nil {
			return nil, fmt.Errorf("could not initialize Keptn Handler: %s", err.Error())
		}
		k.KeptnBase.EventSender = httpSender
	}

	datastoreURL := keptn.DatastoreURL
	if opts.DatastoreURL != "" {
		datastoreURL = opts.DatastoreURL
	}

	k.ResourceHandler = api.NewResourceHandler(csURL)
	k.EventHandler = api.NewEventHandler(datastoreURL)

	loggingServiceName := keptn.DefaultLoggingServiceName
	if opts.LoggingOptions != nil && opts.LoggingOptions.ServiceName != nil {
		loggingServiceName = *opts.LoggingOptions.ServiceName
	}
	k.Logger = keptn.NewLogger(k.KeptnContext, incomingEvent.Context.GetID(), loggingServiceName)

	return k, nil
}

// GetShipyard returns the shipyard definition of a project
func (k *Keptn) GetShipyard() (*keptnv2.Shipyard, error) {
	shipyardResource, err := k.ResourceHandler.GetProjectResource(k.Event.GetProject(), "shipyard.yaml")
	if err != nil {
		return nil, err
	}

	shipyard := keptnv2.Shipyard{}
	err = yaml.Unmarshal([]byte(shipyardResource.ResourceContent), &shipyard)
	if err != nil {
		return nil, err
	}
	return &shipyard, nil
}

// SendCloudEvent sends a cloudevent to the event broker
func (k *Keptn) SendCloudEvent(ctx context.Context, event cloudevents.Event) error {
	event.SetExtension(keptnSpecVersionCEExtension, config.GetKeptnGoUtilsConfig().ShKeptnSpecVersion)
	if k.UseLocalFileSystem {
		log.Println(fmt.Printf("%v", string(event.Data())))
		return nil
	}

	ctx = cloudevents.WithEncodingStructured(ctx)
	return k.EventSender.Send(ctx, event)
}

// SendTaskStartedEvent sends a .started event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskStartedEvent(ctx context.Context, data keptn.EventProperties, source string) (string, error) {
	if k.CloudEvent == nil {
		return "", fmt.Errorf("no incoming .triggered CloudEvent provided to the Keptn Handler")
	}
	outEventType, err := keptnv2.GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnStartedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .started event type for base event: %s", err.Error())
	}

	return k.sendEventWithBaseEventContext(ctx, data, source, err, outEventType)
}

// SendTaskStartedEvent sends a .status.changed event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskStatusChangedEvent(ctx context.Context, data keptn.EventProperties, source string) (string, error) {
	if k.CloudEvent == nil {
		return "", fmt.Errorf("no incoming .triggered CloudEvent provided to the Keptn Handler")
	}
	outEventType, err := keptnv2.GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnStatusChangedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .status.changed event type for base event: %s", err.Error())
	}

	return k.sendEventWithBaseEventContext(ctx, data, source, err, outEventType)
}

// SendTaskStartedEvent sends a .finished event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskFinishedEvent(ctx context.Context, data keptn.EventProperties, source string) (string, error) {
	if k.CloudEvent == nil {
		return "", fmt.Errorf("no incoming .triggered CloudEvent provided to the Keptn Handler")
	}
	outEventType, err := keptnv2.GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnFinishedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .finished event type for base event: %s", err.Error())
	}

	return k.sendEventWithBaseEventContext(ctx, data, source, err, outEventType)
}

func (k *Keptn) sendEventWithBaseEventContext(ctx context.Context, data keptn.EventProperties, source string, err error, outEventType string) (string, error) {
	if source == "" {
		return "", errors.New("must provide non-empty source")
	}
	keptnContext, err := k.CloudEvent.Context.GetExtension(keptnContextCEExtension)
	if err != nil {
		return "", fmt.Errorf("could not determine shkeptncontext of base event: %s", err.Error())
	}

	ce, err := k.createCloudEventWithContextAndPayload(outEventType, keptnContext, source, data)
	if err != nil {
		return "", fmt.Errorf("could not initialize CloudEvent: %s", err.Error())
	}

	ctx = cloudevents.WithEncodingStructured(ctx)
	if err := k.EventSender.Send(ctx, *ce); err != nil {
		return "", fmt.Errorf("could not send CloudEvent: %s", err.Error())
	}
	return ce.ID(), nil
}

// createCloudEventWithContextAndPayload initializes a new CloudEvent and ensures that context attributes such as the triggeredID, keptnContext,
// as well as project, stage, service and labels are included in the resulting event
func (k *Keptn) createCloudEventWithContextAndPayload(outEventType string, keptnContext interface{}, source string, data keptn.EventProperties) (*cloudevents.Event, error) {
	ce := cloudevents.NewEvent()
	ce.SetID(uuid.New().String())
	ce.SetType(outEventType)
	ce.SetDataContentType(cloudevents.ApplicationJSON)
	// use existing  keptnContext for the new cloudevent
	ce.SetExtension(keptnContextCEExtension, keptnContext)
	// the triggeredID links the sent event to the received .triggered event
	ce.SetExtension(triggeredIDCEExtension, k.CloudEvent.ID())
	ce.SetSource(source)

	// if available, add the keptnspecversion extension to the CloudEvent context
	if keptnSpecVersion, err := k.CloudEvent.Context.GetExtension(keptnSpecVersionCEExtension); err == nil && keptnSpecVersion != "" {
		ce.SetExtension(keptnSpecVersionCEExtension, keptnSpecVersion)
	}

	var eventData keptn.EventProperties
	if data != nil {
		eventData = data
	} else {
		eventData = &keptnv2.EventData{}
	}

	eventData = ensureContextAttributesAreSet(k.Event, eventData)
	if err := ce.SetData(cloudevents.ApplicationJSON, eventData); err != nil {
		return nil, fmt.Errorf("could not set data of CloudEvent: %s", err.Error())
	}
	return &ce, nil
}

// ensureContextAttributesAreSet makes sure all properties that remain constant over the course of a task sequence execution are set in the outgoing event
func ensureContextAttributesAreSet(srcEvent, newEvent keptn.EventProperties) keptn.EventProperties {
	newEvent.SetProject(srcEvent.GetProject())
	newEvent.SetStage(srcEvent.GetStage())
	newEvent.SetService(srcEvent.GetService())
	labels := srcEvent.GetLabels()

	// make sure labels from triggered event are included. Existing labels cannot be changed, but new ones can be added
	for key, value := range newEvent.GetLabels() {
		if labels[key] == "" {
			labels[key] = value
		}
	}
	newEvent.SetLabels(labels)
	return newEvent
}
