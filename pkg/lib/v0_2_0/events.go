package v0_2_0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/config"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/keptn"

	ceObs "github.com/cloudevents/sdk-go/observability/opentelemetry/v2/client"
	ceObsHttp "github.com/cloudevents/sdk-go/observability/opentelemetry/v2/http"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	httpprotocol "github.com/cloudevents/sdk-go/v2/protocol/http"
)

const MAX_SEND_RETRIES = 3

const DefaultHTTPEventEndpoint = "http://localhost:8081/event"

const defaultSpecVersion = "1.0"

const keptnEventTypePrefix = "sh.keptn.event."
const keptnTriggeredEventSuffix = ".triggered"
const keptnStartedEventSuffix = ".started"
const keptnStatusChangedEventSuffix = ".status.changed"
const keptnFinishedEventSuffix = ".finished"
const keptnInvalidatedEventSuffix = ".invalidated"

const keptnContextCEExtension = "shkeptncontext"
const keptnSpecVersionCEExtension = "shkeptnspecversion"
const triggeredIDCEExtension = "triggeredid"
const keptnGitCommitIDCEExtension = "gitcommitid"

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
	// Creates a HTTP protocol wrapped with the OpenTelemetry transport
	p, err := ceObsHttp.NewObservedHTTP()
	if err != nil {
		return nil, fmt.Errorf("failed to create protocol: %s", err.Error())
	}

	// An HTTP client with an ObservabilityService
	// THe ObsService will generate spans on send/receive operations
	c, err := cloudevents.NewClient(
		p, cloudevents.WithTimeNow(),
		cloudevents.WithUUIDs(),
		client.WithObservabilityService(ceObs.NewOTelObservabilityService()),
	)

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
	return httpSender.Send(context.Background(), event)
}

func (httpSender HTTPEventSender) Send(ctx context.Context, event cloudevents.Event) error {
	ctx = cloudevents.ContextWithTarget(ctx, httpSender.EventsEndpoint)
	ctx = cloudevents.WithEncodingStructured(ctx)
	var result protocol.Result
	for i := 0; i <= MAX_SEND_RETRIES; i++ {
		result = httpSender.Client.Send(ctx, event)
		httpResult, ok := result.(*httpprotocol.Result)
		switch {
		case ok:
			if httpResult.StatusCode >= 200 && httpResult.StatusCode < 300 {
				return nil
			}
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		case cloudevents.IsUndelivered(result):
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		default:
			return nil
		}
	}
	return errors.New("Failed to send cloudevent: " + result.Error())
}

// EventSender fakes the sending of CloudEvents
type TestSender struct {
	SentEvents []cloudevents.Event
	Reactors   map[string]func(event cloudevents.Event) error
}

// SendEvent fakes the sending of CloudEvents
func (s *TestSender) SendEvent(event cloudevents.Event) error {
	return s.Send(context.TODO(), event)
}

func (s *TestSender) Send(ctx context.Context, event cloudevents.Event) error {
	if s.Reactors != nil {
		for eventTypeSelector, reactor := range s.Reactors {
			if eventTypeSelector == "*" || eventTypeSelector == event.Type() {
				if err := reactor(event); err != nil {
					return err
				}
			}
		}
	}
	s.SentEvents = append(s.SentEvents, event)
	return nil
}

// AssertSentEventTypes checks if the given event types have been passed to the SendEvent function
func (s *TestSender) AssertSentEventTypes(eventTypes []string) error {
	if len(s.SentEvents) != len(eventTypes) {
		return fmt.Errorf("expected %d event, got %d", len(s.SentEvents), len(eventTypes))
	}
	for index, event := range s.SentEvents {
		if event.Type() != eventTypes[index] {
			return fmt.Errorf("received event type '%s' != %s", event.Type(), eventTypes[index])
		}
	}
	return nil
}

// AddReactor adds custom logic that should be applied when SendEvent is called for the given event type
func (s *TestSender) AddReactor(eventTypeSelector string, reactor func(event cloudevents.Event) error) {
	if s.Reactors == nil {
		s.Reactors = map[string]func(event cloudevents.Event) error{}
	}
	s.Reactors[eventTypeSelector] = reactor
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

func IsFinishedEventType(eventType string) bool {
	return strings.HasSuffix(eventType, keptnFinishedEventSuffix)
}

func IsStartedEventType(eventType string) bool {
	return strings.HasSuffix(eventType, keptnStartedEventSuffix)
}

func IsTriggeredEventType(eventType string) bool {
	return strings.HasSuffix(eventType, keptnTriggeredEventSuffix)
}

// EventData contains mandatory fields of all Keptn CloudEvents
type EventData struct {
	Project string            `json:"project,omitempty"`
	Stage   string            `json:"stage,omitempty"`
	Service string            `json:"service,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
	Status  StatusType        `json:"status,omitempty" jsonschema:"enum=succeeded,enum=errored,enum=unknown"`
	Result  ResultType        `json:"result,omitempty" jsonschema:"enum=pass,enum=warning,enum=fail"`
	Message string            `json:"message,omitempty"`
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

// KeptnEvent creates a builder for a new KeptnContextExtendedCE
func KeptnEvent(eventType string, source string, payload interface{}) *KeptnEventBuilder {
	cfg := config.GetKeptnGoUtilsConfig()
	ce := models.KeptnContextExtendedCE{
		ID:                 uuid.NewString(),
		Contenttype:        cloudevents.ApplicationJSON,
		Data:               payload,
		Source:             strutils.Stringp(source),
		Shkeptnspecversion: cfg.ShKeptnSpecVersion,
		Specversion:        defaultSpecVersion,
		Time:               time.Now().UTC(),
		Type:               strutils.Stringp(eventType),
	}

	return &KeptnEventBuilder{ce}
}

// KeptnEventBuilder is used for constructing a new KeptnContextExtendedCE
type KeptnEventBuilder struct {
	models.KeptnContextExtendedCE
}

// Build creates a value of KeptnContextExtendedCE from the current builder
// It also does basic validation like the presence of project, service and stage in the event data
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

// WithKeptnSpecVersion can be used to override the keptn spec version
func (eb *KeptnEventBuilder) WithKeptnSpecVersion(keptnSpecVersion string) *KeptnEventBuilder {
	eb.Shkeptnspecversion = keptnSpecVersion
	return eb
}

// WithKeptnContext can be used to set a keptn context
func (eb *KeptnEventBuilder) WithKeptnContext(keptnContext string) *KeptnEventBuilder {
	eb.Shkeptncontext = keptnContext
	return eb
}

// WithTriggeredID can be used to set the triggered ID
func (eb *KeptnEventBuilder) WithTriggeredID(triggeredID string) *KeptnEventBuilder {
	eb.Triggeredid = triggeredID
	return eb
}

// WithGitCommitID can be used to set the git commit ID
func (eb *KeptnEventBuilder) WithGitCommitID(gitCommitID string) *KeptnEventBuilder {
	eb.GitCommitID = gitCommitID
	return eb
}

// WithID can be used to override the ID, which is auto generated by default
func (eb *KeptnEventBuilder) WithID(id string) *KeptnEventBuilder {
	eb.ID = id
	return eb
}

// ToCloudEvent takes a KeptnContextExtendedCE and converts it to an ordinary CloudEvent
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
	event.SetExtension(keptnGitCommitIDCEExtension, keptnEvent.GitCommitID)
	event.SetExtension(keptnSpecVersionCEExtension, keptnEvent.Shkeptnspecversion)
	return event
}

// ToKeptnEvent takes a CloudEvent and converts it into a KeptnContextExtendedCE
func ToKeptnEvent(event cloudevents.Event) (models.KeptnContextExtendedCE, error) {
	var keptnContext string
	event.ExtensionAs(keptnContextCEExtension, &keptnContext)

	var triggeredID string
	event.ExtensionAs(triggeredIDCEExtension, &triggeredID)

	var keptnSpecVersion string
	event.ExtensionAs(keptnSpecVersionCEExtension, &keptnSpecVersion)

	var gitCommitID string
	event.ExtensionAs(keptnGitCommitIDCEExtension, &gitCommitID)

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
		Time:               event.Time(),
		Triggeredid:        triggeredID,
		GitCommitID:        gitCommitID,
		Type:               strutils.Stringp(event.Type()),
	}

	return keptnEvent, nil
}
