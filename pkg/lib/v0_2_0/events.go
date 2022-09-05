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

const (
	shkeptnspecversion = "0.2.4"
	cloudeventsversion = "1.0"
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

type HTTPSenderOption func(httpSender *HTTPEventSender)

// WithSendRetries allows to specify the number of retries that are performed if the receiver of an event returns a HTTP error code
func WithSendRetries(retries int) HTTPSenderOption {
	return func(httpSender *HTTPEventSender) {
		httpSender.nrRetries = retries
	}
}

// HTTPEventSender sends CloudEvents via HTTP
type HTTPEventSender struct {
	// EventsEndpoint is the http endpoint the events are sent to
	EventsEndpoint string
	// Client is an implementation of the cloudevents.Client interface
	Client cloudevents.Client
	// nrRetries is the number of retries that are attempted if the endpoint an event is forwarded to returns an http code outside the 2xx range
	nrRetries int
}

// NewHTTPEventSender creates a new HTTPSender
func NewHTTPEventSender(endpoint string, opts ...HTTPSenderOption) (*HTTPEventSender, error) {
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
		nrRetries:      MAX_SEND_RETRIES,
	}

	for _, o := range opts {
		o(httpSender)
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
	for i := 0; i <= httpSender.nrRetries; i++ {
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
	return fmt.Errorf("could not send cloudevent after %d retries. Received result from the receiver: %w", httpSender.nrRetries, result)
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
	event.SetTime(keptnEvent.Time)
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

// CreateStartedEvent takes a parent event (e.g. .triggered event) and creates a corresponding .started event
func CreateStartedEvent(source string, parentEvent models.KeptnContextExtendedCE) (*models.KeptnContextExtendedCE, error) {
	if parentEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from parent event %s", parentEvent.ID)
	}
	startedEventType, err := ReplaceEventTypeKind(*parentEvent.Type, "started")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.started' event for parent event %s: %w", parentEvent.ID, err)
	}
	eventData := EventData{}
	parentEvent.DataAs(&eventData)
	return createEvent(source, startedEventType, parentEvent, eventData), nil
}

// CreateFinishedEvent takes a parent event (e.g. .triggered event) and creates a corresponding .finished event
func CreateFinishedEvent(source string, parentEvent models.KeptnContextExtendedCE, eventData interface{}) (*models.KeptnContextExtendedCE, error) {
	if parentEvent.Type == nil {
		return nil, fmt.Errorf("unable to get keptn event type from event %s", parentEvent.ID)
	}

	if parentEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from parent event %s", parentEvent.ID)
	}
	finishedEventType, err := ReplaceEventTypeKind(*parentEvent.Type, "finished")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event: %v from %s", err, *parentEvent.Type)
	}
	var genericEventData map[string]interface{}
	err = Decode(eventData, &genericEventData)
	if err != nil || genericEventData == nil {
		return nil, fmt.Errorf("unable to decode generic event data")
	}

	if genericEventData["status"] == nil || genericEventData["status"] == "" {
		genericEventData["status"] = "succeeded"
	}

	if genericEventData["result"] == nil || genericEventData["result"] == "" {
		genericEventData["result"] = "pass"
	}
	return createEvent(source, finishedEventType, parentEvent, genericEventData), nil
}

// CreateFinishedEventWithError takes a parent event (e.g. .triggered event) and creates a corresponding errored .finished event
func CreateFinishedEventWithError(source string, parentEvent models.KeptnContextExtendedCE, eventData interface{}, errVal *Error) (*models.KeptnContextExtendedCE, error) {
	if errVal == nil {
		errVal = &Error{}
	}
	commonEventData := EventData{}
	if eventData == nil {
		parentEvent.DataAs(&commonEventData)
	}
	commonEventData.Result = errVal.ResultType
	commonEventData.Status = errVal.StatusType
	commonEventData.Message = errVal.Message

	finishedEventType, err := ReplaceEventTypeKind(*parentEvent.Type, "finished")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event for parent event %s: %w", parentEvent.ID, err)
	}
	return createEvent(source, finishedEventType, parentEvent, commonEventData), nil
}

// CreateErrorEvent takes a parent event (e.g. .triggered event) and creates a corresponding errored event
func CreateErrorEvent(source string, parentEvent models.KeptnContextExtendedCE, eventData interface{}, errVal *Error) (*models.KeptnContextExtendedCE, error) {
	if errVal == nil {
		errVal = &Error{}
	}

	if IsTaskEventType(*parentEvent.Type) && IsTriggeredEventType(*parentEvent.Type) {
		errorFinishedEvent, err := CreateFinishedEventWithError(source, parentEvent, eventData, errVal)
		if err != nil {
			return nil, err
		}
		return errorFinishedEvent, nil
	}
	errorLogEvent, err := CreateErrorLogEvent(source, parentEvent, eventData, errVal)
	if err != nil {
		return nil, err
	}
	return errorLogEvent, nil
}

// CreateErrorEvent takes a parent event (e.g. .triggered event) and creates a corresponding errored .log event
func CreateErrorLogEvent(source string, parentEvent models.KeptnContextExtendedCE, eventData interface{}, errVal *Error) (*models.KeptnContextExtendedCE, error) {
	if parentEvent.Type == nil {
		return nil, fmt.Errorf("unable to get keptn event type from parent event %s", parentEvent.ID)
	}

	if parentEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from parent event %s", parentEvent.ID)
	}
	if errVal == nil {
		errVal = &Error{}
	}

	if IsTaskEventType(*parentEvent.Type) && IsTriggeredEventType(*parentEvent.Type) {
		errorFinishedEvent, err := CreateFinishedEventWithError(source, parentEvent, eventData, errVal)
		if err != nil {
			return nil, err
		}
		return errorFinishedEvent, nil
	}
	errorEventData := ErrorLogEvent{}
	if eventData == nil {
		parentEvent.DataAs(&errorEventData)
	}
	if IsTaskEventType(*parentEvent.Type) {
		taskName, _, err := ParseTaskEventType(*parentEvent.Type)
		if err == nil && taskName != "" {
			errorEventData.Task = taskName
		}
	}
	errorEventData.Message = errVal.Message
	if parentEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from parent event %s", parentEvent.ID)
	}
	return createEvent(source, ErrorLogEventName, parentEvent, errorEventData), nil
}

func createEvent(source string, eventType string, parentEvent models.KeptnContextExtendedCE, eventData interface{}) *models.KeptnContextExtendedCE {
	return &models.KeptnContextExtendedCE{
		ID:                 uuid.NewString(),
		Triggeredid:        parentEvent.ID,
		Shkeptncontext:     parentEvent.Shkeptncontext,
		Contenttype:        cloudevents.ApplicationJSON,
		Data:               eventData,
		Source:             strutils.Stringp(source),
		Shkeptnspecversion: shkeptnspecversion,
		Specversion:        cloudeventsversion,
		Time:               time.Now().UTC(),
		Type:               strutils.Stringp(eventType),
	}
}
