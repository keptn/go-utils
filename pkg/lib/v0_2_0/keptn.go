package v0_2_0

import (
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"gopkg.in/yaml.v2"
	"log"
)

type Keptn struct {
	keptn.KeptnBase
}

func NewKeptn(incomingEvent *cloudevents.Event, opts keptn.KeptnOpts) (*Keptn, error) {
	extension, _ := incomingEvent.Context.GetExtension("shkeptncontext")
	shkeptncontext := extension.(string)

	// create a base KeptnBase Event
	keptnBase := &EventData{}

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
		httpSender, err := NewHTTPEventSender(opts.EventBrokerURL)
		if err != nil {
			return nil, fmt.Errorf("could not initialize Keptn Handler: %s", err.Error())
		}
		k.KeptnBase.EventSender = httpSender
	} else if opts.EventSender != nil {
		k.KeptnBase.EventSender = opts.EventSender
	} else {
		httpSender, err := NewHTTPEventSender(DefaultHTTPEventEndpoint)
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
func (k *Keptn) GetShipyard() (*Shipyard, error) {
	shipyardResource, err := k.ResourceHandler.GetProjectResource(k.Event.GetProject(), "shipyard.yaml")
	if err != nil {
		return nil, err
	}

	shipyard := Shipyard{}
	err = yaml.Unmarshal([]byte(shipyardResource.ResourceContent), &shipyard)
	if err != nil {
		return nil, err
	}
	return &shipyard, nil
}

// SendCloudEvent sends a cloudevent to the event broker
func (k *Keptn) SendCloudEvent(event cloudevents.Event) error {
	if k.UseLocalFileSystem {
		log.Println(fmt.Printf("%v", string(event.Data())))
		return nil
	}

	return k.EventSender.SendEvent(event)
}

// SendTaskStartedEvent sends a .started event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskStartedEvent(data keptn.EventProperties, source string) (string, error) {

	outEventType, err := GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnStartedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .started event type for base event: %s", err.Error())
	}

	return k.sendEventWithBaseEventContext(data, source, err, outEventType)
}

// SendTaskStartedEvent sends a .status.changed event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskStatusChangedEvent(data keptn.EventProperties, source string) (string, error) {
	outEventType, err := GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnStatusChangedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .started event type for base event: %s", err.Error())
	}

	return k.sendEventWithBaseEventContext(data, source, err, outEventType)
}

// SendTaskStartedEvent sends a .finished event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskFinishedEvent(data keptn.EventProperties, source string) (string, error) {
	outEventType, err := GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnFinishedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .started event type for base event: %s", err.Error())
	}

	return k.sendEventWithBaseEventContext(data, source, err, outEventType)
}

func (k *Keptn) sendEventWithBaseEventContext(data keptn.EventProperties, source string, err error, outEventType string) (string, error) {
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

	if err := k.EventSender.SendEvent(*ce); err != nil {
		return "", fmt.Errorf("could not send CloudEvent: %s", err.Error())
	}
	return ce.ID(), nil
}

func (k *Keptn) createCloudEventWithContextAndPayload(outEventType string, keptnContext interface{}, source string, data keptn.EventProperties) (*cloudevents.Event, error) {
	ce := cloudevents.NewEvent()
	ce.SetID(uuid.New().String())
	ce.SetType(outEventType)
	ce.SetDataContentType(cloudevents.ApplicationJSON)
	ce.SetExtension(keptnContextCEExtension, keptnContext)
	ce.SetExtension(triggeredIDCEExtenstion, k.CloudEvent.ID())
	ce.SetSource(source)

	var eventData keptn.EventProperties
	if data != nil {
		eventData = data
	} else {
		eventData = &EventData{}
	}

	eventData = ensureContextAttributesAreSet(k.Event, eventData)
	if err := ce.SetData(cloudevents.ApplicationJSON, eventData); err != nil {
		return nil, fmt.Errorf("could not set data of CloudEvent: %s", err.Error())
	}
	return &ce, nil
}

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
