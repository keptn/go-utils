package v0_2_0

import (
	"errors"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/config"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"gopkg.in/yaml.v3"
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
	event.SetExtension(keptnSpecVersionCEExtension, config.GetKeptnGoUtilsConfig().ShKeptnSpecVersion)
	if k.UseLocalFileSystem {
		log.Println(fmt.Printf("%v", string(event.Data())))
		return nil
	}

	return k.EventSender.SendEvent(event)
}

// SendTaskStartedEvent sends a .started event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskStartedEvent(data keptn.EventProperties, source string) (string, error) {
	if k.CloudEvent == nil {
		return "", fmt.Errorf("no incoming .triggered CloudEvent provided to the Keptn Handler")
	}
	outEventType, err := GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnStartedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .started event type for base event: %s", err.Error())
	}

	return k.sendEventWithBaseEventContext(data, source, err, outEventType)
}

// SendTaskStartedEvent sends a .status.changed event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskStatusChangedEvent(data keptn.EventProperties, source string) (string, error) {
	if k.CloudEvent == nil {
		return "", fmt.Errorf("no incoming .triggered CloudEvent provided to the Keptn Handler")
	}
	outEventType, err := GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnStatusChangedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .status.changed event type for base event: %s", err.Error())
	}

	return k.sendEventWithBaseEventContext(data, source, err, outEventType)
}

// SendTaskFinishedEvent sends a .finished event for the incoming .triggered event the KeptnHandler was initialized with.
// It returns the ID of the sent CloudEvent or an error
func (k *Keptn) SendTaskFinishedEvent(data keptn.EventProperties, source string) (string, error) {
	if k.CloudEvent == nil {
		return "", fmt.Errorf("no incoming .triggered CloudEvent provided to the Keptn Handler")
	}
	outEventType, err := GetEventTypeForTriggeredEvent(k.CloudEvent.Type(), keptnFinishedEventSuffix)
	if err != nil {
		return "", fmt.Errorf("could not determine .finished event type for base event: %s", err.Error())
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

// createCloudEventWithContextAndPayload initializes a new CloudEvent and ensures that context attributes such as the triggeredID, keptnContext, gitcommitid
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

	// if available, add the gitcommitid extension to the CloudEvent context
	if keptnGitCommitID, err := k.CloudEvent.Context.GetExtension(keptnGitCommitIDCEExtension); err == nil && keptnGitCommitID != "" {
		ce.SetExtension(keptnGitCommitIDCEExtension, keptnGitCommitID)
	}

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

// ensureContextAttributesAreSet makes sure all properties that remain constant over the course of a task sequence execution are set in the outgoing event
func ensureContextAttributesAreSet(srcEvent, newEvent keptn.EventProperties) keptn.EventProperties {
	newEvent.SetProject(srcEvent.GetProject())
	newEvent.SetStage(srcEvent.GetStage())
	newEvent.SetService(srcEvent.GetService())
	labels := srcEvent.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	// make sure labels from triggered event are included. Existing labels cannot be changed, but new ones can be added
	for key, value := range newEvent.GetLabels() {
		if labels[key] == "" {
			labels[key] = value
		}
	}
	newEvent.SetLabels(labels)
	return newEvent
}
