package eventutils

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-openapi/strfmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"strings"
	"time"
)

const contentType = "application/json"
const defaultSpecVersion = "1.0"
const defaultKeptnSpecVersion = "0.2.1"
const keptnContextExtension = "shkeptncontext"
const keptnTriggeredIdExtension = "triggeredid"
const keptnSpecVersionExtension = "shkeptnspecversion"

func KeptnEvent(eventType string, payload interface{}) *KeptnEventBuilder {

	ce := models.KeptnContextExtendedCE{
		Contenttype:        contentType,
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

	commonEventData := v0_2_0.EventData{}
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
	event.SetExtension(keptnContextExtension, keptnEvent.Shkeptncontext)
	event.SetExtension(keptnTriggeredIdExtension, keptnEvent.Triggeredid)
	event.SetExtension(keptnSpecVersionExtension, keptnEvent.Shkeptnspecversion)
	return event
}

func ToKeptnEvent(event cloudevents.Event) (models.KeptnContextExtendedCE, error) {
	keptnContext := ""
	if err := event.ExtensionAs(keptnContextExtension, &keptnContext); err != nil {
		return models.KeptnContextExtendedCE{}, err
	}

	triggeredID := ""
	if err := event.ExtensionAs(keptnTriggeredIdExtension, &triggeredID); err != nil {
		return models.KeptnContextExtendedCE{}, err
	}

	keptnSpecVersion := ""
	if err := event.ExtensionAs(keptnSpecVersionExtension, &keptnSpecVersion); err != nil {
		return models.KeptnContextExtendedCE{}, err
	}

	var data interface{}
	event.DataAs(&data)

	keptnEvent := models.KeptnContextExtendedCE{
		Contenttype:        contentType,
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

func IsTaskEventType(eventType string) bool {
	parts := strings.Split(eventType, ".")
	return len(parts) == 5
}

func IsSequenceEventType(eventType string) bool {
	parts := strings.Split(eventType, ".")
	return len(parts) == 6
}

func IsValidEventType(eventType string) bool {
	return IsSequenceEventType(eventType) || IsTaskEventType(eventType)
}

func ParseSequenceEventType(sequenceTriggeredEventType string) (string, string, string, error) {
	parts := strings.Split(sequenceTriggeredEventType, ".")
	if IsSequenceEventType(sequenceTriggeredEventType) {
		return parts[3], parts[4], parts[5], nil
	}
	return "", "", "", fmt.Errorf("%s is not a valid keptn sequence triggered event type", sequenceTriggeredEventType)
}

func ParseTaskEventType(taskEventType string) (string, string, error) {
	if !IsTaskEventType(taskEventType) {
		return "", "", fmt.Errorf("%s is not a valid keptn task event type", taskEventType)
	}
	parts := strings.Split(taskEventType, ".")
	return parts[3], parts[4], nil
}

func ParseEventKind(eventType string) (string, error) {
	if !IsValidEventType(eventType) {
		return "", fmt.Errorf("%s is not a valid keptn event type", eventType)
	}
	parts := strings.Split(eventType, ".")
	return parts[len(parts)-1], nil
}

func ParseEventTypeWithoutKind(eventType string) (string, error) {
	if !IsValidEventType(eventType) {
		return "", fmt.Errorf("%s is not a valid keptn event type", eventType)
	}
	kind, _ := ParseEventKind(eventType)
	return strings.TrimSuffix(eventType, "."+kind), nil
}

func ReplaceEventTypeKind(eventType, newKind string) (string, error) {
	if !IsValidEventType(eventType) {
		return "", fmt.Errorf("%s is not a valid keptn event type", eventType)
	}
	parts := strings.Split(eventType, ".")
	parts[len(parts)-1] = newKind
	return strings.Join(parts, "."), nil
}
