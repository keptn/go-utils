package eventutils

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-openapi/strfmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"time"
)

const contentType = "application/json"
const eventTypePrefix = "sh.keptn.event"
const defaultSpecVersion = "1.0"
const defaultKeptnSpecVersion = "0.2.1"
const keptnContextExtension = "shkeptncontext"
const keptnTriggeredIdExtension = "triggeredid"
const keptnSpecVersionExtension = "shkeptnspecversion"

func KeptnEvent(eventType string, payload interface{}) *KeptnEventBuilder {
	return &KeptnEventBuilder{
		KeptnContextExtendedCE: models.KeptnContextExtendedCE{
			Contenttype:        contentType,
			Data:               payload,
			Shkeptnspecversion: defaultKeptnSpecVersion,
			Specversion:        defaultSpecVersion,
			Time:               strfmt.DateTime(time.Now().UTC()),
			Type:               strutils.Stringp(eventTypePrefix + "." + eventType),
		},
	}
}

type KeptnEventBuilder struct {
	models.KeptnContextExtendedCE
}

func (eb *KeptnEventBuilder) Build() models.KeptnContextExtendedCE {
	return eb.KeptnContextExtendedCE
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
