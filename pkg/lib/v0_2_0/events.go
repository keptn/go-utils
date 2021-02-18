package v0_2_0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	httpprotocol "github.com/cloudevents/sdk-go/v2/protocol/http"
)

const MAX_SEND_RETRIES = 3

const DefaultHTTPEventEndpoint = "http://localhost:8081/event"

const keptnEventTypePrefix = "sh.keptn.event."
const keptnTriggeredEventSuffix = ".triggered"
const keptnStartedEventSuffix = ".started"
const keptnStatusChangedEventSuffix = ".status.changed"
const keptnFinishedEventSuffix = ".finished"

const keptnContextCEExtension = "shkeptncontext"
const keptnSpecVersionCEExtension = "shkeptnspecversion"
const triggeredIDCEExtension = "triggeredid"

// HTTPEventSender sends CloudEvents via HTTP
type HTTPEventSender struct {
	// EventsEndpoint is the http endpoint the events are sent to
	EventsEndpoint string
	// Client is an implementation of the cloudevents.Client interface
	Client cloudevents.Client
}

// NewHTTPEventSender creates a new HTTPSender
func NewHTTPEventSender(endpoint string) (*HTTPEventSender, error) {
	if endpoint == "" {
		endpoint = DefaultHTTPEventEndpoint
	}
	p, err := cloudevents.NewHTTP()
	if err != nil {
		return nil, fmt.Errorf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, fmt.Errorf("failed to create client, %s", err.Error())
	}

	httpSender := &HTTPEventSender{
		EventsEndpoint: endpoint,
		Client:         c,
	}
	return httpSender, nil
}

// SendEvent sends a CloudEvent
func (httpSender HTTPEventSender) SendEvent(event cloudevents.Event) error {
	ctx := cloudevents.ContextWithTarget(context.Background(), httpSender.EventsEndpoint)
	ctx = cloudevents.WithEncodingStructured(ctx)

	var result protocol.Result
	for i := 0; i <= MAX_SEND_RETRIES; i++ {
		result = httpSender.Client.Send(ctx, event)
		httpResult, ok := result.(*httpprotocol.Result)
		if ok {
			if httpResult.StatusCode >= 200 && httpResult.StatusCode < 300 {
				return nil
			}
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		} else if cloudevents.IsUndelivered(result) {
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		} else {
			return nil
		}
	}
	return errors.New("Failed to send cloudevent: " + result.Error())
}

// GetTriggeredEventType returns for the given task the name of the triggered event type
func GetTriggeredEventType(task string) string {

	return keptnEventTypePrefix + task + keptnTriggeredEventSuffix
}

// GetStartedEventType returns for the given task the name of the started event type
func GetStartedEventType(task string) string {
	return keptnEventTypePrefix + task + keptnStartedEventSuffix
}

// GetStatusChangedEventType returns for the given task the name of the status.changed event type
func GetStatusChangedEventType(task string) string {
	return keptnEventTypePrefix + task + keptnStatusChangedEventSuffix
}

// GetFinishedEventType returns for the given task the name of the finished event type
func GetFinishedEventType(task string) string {
	return keptnEventTypePrefix + task + keptnFinishedEventSuffix
}

func GetEventTypeForTriggeredEvent(baseTriggeredEventType, newEventTypeSuffix string) (string, error) {
	if !strings.HasSuffix(baseTriggeredEventType, keptnTriggeredEventSuffix) {
		return "", errors.New("provided baseTriggeredEventType is not a .triggered event type")
	}
	trimmed := strings.TrimSuffix(baseTriggeredEventType, keptnTriggeredEventSuffix)
	return trimmed + newEventTypeSuffix, nil
}

// EventData contains mandatory fields of all Keptn CloudEvents
type EventData struct {
	Project string            `json:"project,omitempty"`
	Stage   string            `json:"stage,omitempty"`
	Service string            `json:"service,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`

	Status  StatusType `json:"status,omitempty" jsonschema:"enum=succeeded,enum=errored,enum=unknown"`
	Result  ResultType `json:"result,omitempty" jsonschema:"enum=pass,enum=warning,enum=fail"`
	Message string     `json:"message,omitempty"`
}

func (e *EventData) GetProject() string {
	return e.Project
}

func (e *EventData) GetStage() string {
	return e.Stage
}

func (e *EventData) GetService() string {
	return e.Service
}

func (e *EventData) GetLabels() map[string]string {
	return e.Labels
}

func (e *EventData) SetProject(project string) {
	e.Project = project
}

func (e *EventData) SetStage(stage string) {
	e.Stage = stage
}

func (e *EventData) SetService(service string) {
	e.Service = service
}

func (e *EventData) SetLabels(labels map[string]string) {
	e.Labels = labels
}

// Decode decodes the given raw interface to the target pointer specified
// by the out parameter
func Decode(in, out interface{}) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}

// EventDataAs decodes the event data of the given keptn cloud event to the
// target pointer specified by the out parameter
func EventDataAs(in models.KeptnContextExtendedCE, out interface{}) error {
	return Decode(in.Data, out)
}
