package v0_2_0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/go-openapi/strfmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	httpprotocol "github.com/cloudevents/sdk-go/v2/protocol/http"
)

const MAX_SEND_RETRIES = 3

const DefaultHTTPEventEndpoint = "http://localhost:8081/event"

const defaultSpecVersion = "1.0"
const defaultKeptnSpecVersion = "0.2.1"

const keptnEventTypePrefix = "sh.keptn.event."
const keptnTriggeredEventSuffix = ".triggered"
const keptnStartedEventSuffix = ".started"
const keptnStatusChangedEventSuffix = ".status.changed"
const keptnFinishedEventSuffix = ".finished"
const keptnInvalidatedEventSuffix = ".invalidated"

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

// GetInvalidatedEventType returns for the given task the name of the finished event type
func GetInvalidatedEventType(task string) string {
	return keptnEventTypePrefix + task + keptnInvalidatedEventSuffix
}

// IsTaskEventType checks whether the given eventType is a task event type like e.g. "sh.keptn.event.task.triggered"
func IsTaskEventType(eventType string) bool {
	parts := strings.Split(eventType, ".")
	if len(parts) != 5 {
		return false
	}
	for _, p := range parts {
		if p == "" {
			return false
		}
	}
	return true
}

// IsSequenceEventType checks whether the given event type is a sequence event type like e.g. "sh.keptn.event.stage.sequence.triggered"
func IsSequenceEventType(eventType string) bool {
	parts := strings.Split(eventType, ".")
	if len(parts) != 6 {
		return false
	}
	for _, p := range parts {
		if p == "" {
			return false
		}
	}
	return true
}

// IsValidEventType checks whether the given event type is a valid event type, i.e. a valid task event type or sequence event type
func IsValidEventType(eventType string) bool {
	return IsSequenceEventType(eventType) || IsTaskEventType(eventType)
}

// ParseSequenceEventType parses the given sequence event type and returns the stage name, sequence name, event type as well as an error which
// is eventually nil
func ParseSequenceEventType(sequenceTriggeredEventType string) (string, string, string, error) {
	parts := strings.Split(sequenceTriggeredEventType, ".")
	if IsSequenceEventType(sequenceTriggeredEventType) {
		return parts[3], parts[4], parts[5], nil
	}
	return "", "", "", fmt.Errorf("%s is not a valid keptn sequence triggered event type", sequenceTriggeredEventType)
}

// ParseTaskEventType parses the given task event type and returns the task name, event type as well as an error which
// is eventually nil
func ParseTaskEventType(taskEventType string) (string, string, error) {
	if !IsTaskEventType(taskEventType) {
		return "", "", fmt.Errorf("%s is not a valid keptn task event type", taskEventType)
	}
	parts := strings.Split(taskEventType, ".")
	return parts[3], parts[4], nil
}

// ParseEventKind parses the given event type and returns the last element which is the "kind" of the event (e.g. triggered, finished, ...)
func ParseEventKind(eventType string) (string, error) {
	if !IsValidEventType(eventType) {
		return "", fmt.Errorf("%s is not a valid keptn event type", eventType)
	}
	parts := strings.Split(eventType, ".")
	return parts[len(parts)-1], nil
}

// ParseEventTypeWithoutKind parses the given event type and trims away the last element of the event which is the "kind"  of the event (e.g. triggered, finished, ...)
func ParseEventTypeWithoutKind(eventType string) (string, error) {
	if !IsValidEventType(eventType) {
		return "", fmt.Errorf("%s is not a valid keptn event type", eventType)
	}
	kind, _ := ParseEventKind(eventType)
	return strings.TrimSuffix(eventType, "."+kind), nil
}

// ReplaceEventTypeKind replaces the last element of the event which is the "kind" of the event (e.g. triggered, finished, ...) with a new value
// This is useful e.g. to transform a .started event type into a .finished event type
func ReplaceEventTypeKind(eventType, newKind string) (string, error) {
	if !IsValidEventType(eventType) {
		return "", fmt.Errorf("%s is not a valid keptn event type", eventType)
	}
	if newKind == "" {
		return ParseEventTypeWithoutKind(eventType)
	}
	parts := strings.Split(eventType, ".")
	parts[len(parts)-1] = newKind
	return strings.Join(parts, "."), nil
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

func KeptnEvent(eventType string, payload interface{}) *KeptnEventBuilder {

	ce := models.KeptnContextExtendedCE{
		Contenttype:        cloudevents.ApplicationJSON,
		Data:               payload,
		Source:             strutils.Stringp(""),
		Shkeptnspecversion: defaultKeptnSpecVersion,
		Specversion:        defaultSpecVersion,
		Time:               strfmt.DateTime(time.Now().UTC()),
		Type:               strutils.Stringp(eventType),
	}

	return &KeptnEventBuilder{ce}
}

type KeptnEventBuilder struct {
	models.KeptnContextExtendedCE
}

func (eb *KeptnEventBuilder) Build() (models.KeptnContextExtendedCE, error) {

	commonEventData := EventData{}
	if err := eb.DataAs(&commonEventData); err != nil {
		return eb.KeptnContextExtendedCE, err
	}
	if commonEventData.Project == "" || commonEventData.Service == "" || commonEventData.Stage == "" {
		return eb.KeptnContextExtendedCE, fmt.Errorf("cannot create keptn cloud event as it does not contain project, service and stage information")
	}

	return eb.KeptnContextExtendedCE, nil
}

func (eb *KeptnEventBuilder) WithKeptnSpecVersion(keptnSpecVersion string) *KeptnEventBuilder {
	eb.Shkeptnspecversion = keptnSpecVersion
	return eb
}

func (eb *KeptnEventBuilder) WithKeptnContext(keptnContext string) *KeptnEventBuilder {
	eb.Shkeptncontext = keptnContext
	return eb
}

func (eb *KeptnEventBuilder) WithSource(source string) *KeptnEventBuilder {
	eb.Source = &source
	return eb
}

func (eb *KeptnEventBuilder) WithTriggeredID(triggeredID string) *KeptnEventBuilder {
	eb.Triggeredid = triggeredID
	return eb
}

func (eb *KeptnEventBuilder) WithID(id string) *KeptnEventBuilder {
	eb.ID = id
	return eb
}

func ToCloudEvent(keptnEvent models.KeptnContextExtendedCE) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetType(*keptnEvent.Type)
	event.SetID(keptnEvent.ID)
	event.SetSource(*keptnEvent.Source)
	event.SetDataContentType(keptnEvent.Contenttype)
	event.SetSpecVersion(keptnEvent.Specversion)
	event.SetData(cloudevents.ApplicationJSON, keptnEvent.Data)
	event.SetExtension(keptnContextCEExtension, keptnEvent.Shkeptncontext)
	event.SetExtension(triggeredIDCEExtension, keptnEvent.Triggeredid)
	event.SetExtension(keptnSpecVersionCEExtension, keptnEvent.Shkeptnspecversion)
	return event
}

func ToKeptnEvent(event cloudevents.Event) (models.KeptnContextExtendedCE, error) {
	keptnContext := ""
	if err := event.ExtensionAs(keptnContextCEExtension, &keptnContext); err != nil {
		return models.KeptnContextExtendedCE{}, err
	}

	triggeredID := ""
	if err := event.ExtensionAs(triggeredIDCEExtension, &triggeredID); err != nil {
		return models.KeptnContextExtendedCE{}, err
	}

	keptnSpecVersion := ""
	if err := event.ExtensionAs(keptnSpecVersionCEExtension, &keptnSpecVersion); err != nil {
		return models.KeptnContextExtendedCE{}, err
	}

	var data interface{}
	event.DataAs(&data)

	keptnEvent := models.KeptnContextExtendedCE{
		Contenttype:        cloudevents.ApplicationJSON,
		Data:               data,
		ID:                 event.ID(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: keptnSpecVersion,
		Source:             strutils.Stringp(event.Source()),
		Specversion:        event.SpecVersion(),
		Time:               strfmt.DateTime(event.Time()),
		Triggeredid:        triggeredID,
		Type:               strutils.Stringp(event.Type()),
	}

	return keptnEvent, nil
}
