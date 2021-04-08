package eventutils

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBuilder(t *testing.T) {
	KeptnEvent("type", nil).Build()
}
func TestCreateSimpleKeptnEvent(t *testing.T) {

	event := KeptnEvent("dev.delivery.triggered", map[string]interface{}{}).Build()
	require.Equal(t, "application/json", event.Contenttype)
	require.Equal(t, map[string]interface{}{}, event.Data)
	require.Equal(t, "", event.Shkeptncontext)
	require.Equal(t, time.Now().UTC().Round(time.Minute), time.Time(event.Time).Round(time.Minute))
	require.Equal(t, defaultKeptnSpecVersion, event.Shkeptnspecversion)
	require.Equal(t, defaultSpecVersion, event.Specversion)
	require.Equal(t, "", event.Triggeredid)
	require.Equal(t, strutils.Stringp("sh.keptn.event.dev.delivery.triggered"), event.Type)
}

func TestCreateKeptnEvent(t *testing.T) {

	event := KeptnEvent("dev.delivery.triggered", map[string]interface{}{}).
		WithID("my-id").
		WithKeptnContext("my-keptn-context").
		WithSource("my-source").
		WithTriggeredID("my-triggered-id").
		Build()

	require.Equal(t, "application/json", event.Contenttype)
	require.Equal(t, map[string]interface{}{}, event.Data)
	require.Equal(t, defaultKeptnSpecVersion, event.Shkeptnspecversion)
	require.Equal(t, defaultSpecVersion, event.Specversion)
	require.Equal(t, "my-id", event.ID)
	require.Equal(t, time.Now().UTC().Round(time.Minute), time.Time(event.Time).Round(time.Minute))
	require.Equal(t, strutils.Stringp("my-source"), event.Source)
	require.Equal(t, "my-keptn-context", event.Shkeptncontext)
	require.Equal(t, "my-triggered-id", event.Triggeredid)
	require.Equal(t, strutils.Stringp("sh.keptn.event.dev.delivery.triggered"), event.Type)
}

func TestToCloudEvent(t *testing.T) {

	type TestData struct {
		Content string `json:"content"`
	}

	expected := cloudevents.NewEvent()
	expected.SetType("sh.keptn.event.dev.delivery.triggered")
	expected.SetID("my-id")
	expected.SetSource("my-source")
	expected.SetData(cloudevents.ApplicationJSON, TestData{Content: "testdata"})
	expected.SetDataContentType(contentType)
	expected.SetSpecVersion(defaultSpecVersion)
	expected.SetExtension(keptnContextExtension, "my-keptn-context")
	expected.SetExtension(keptnTriggeredIdExtension, "my-triggered-id")
	expected.SetExtension(keptnSpecVersionExtension, defaultKeptnSpecVersion)

	keptnEvent := models.KeptnContextExtendedCE{
		Contenttype:        "application/json",
		Data:               TestData{Content: "testdata"},
		ID:                 "my-id",
		Shkeptncontext:     "my-keptn-context",
		Source:             strutils.Stringp("my-source"),
		Shkeptnspecversion: defaultKeptnSpecVersion,
		Specversion:        defaultSpecVersion,
		Triggeredid:        "my-triggered-id",
		Type:               strutils.Stringp("sh.keptn.event.dev.delivery.triggered"),
	}
	cloudevent := ToCloudEvent(keptnEvent)
	assert.Equal(t, expected, cloudevent)

}

func TestToKeptnEvent(t *testing.T) {

	type TestData struct {
		Content string `json:"content"`
	}

	expected := models.KeptnContextExtendedCE{
		Contenttype:        "application/json",
		Data:               map[string]interface{}{"content": "testdata"},
		ID:                 "my-id",
		Shkeptncontext:     "my-keptn-context",
		Source:             strutils.Stringp("my-source"),
		Shkeptnspecversion: defaultKeptnSpecVersion,
		Specversion:        defaultSpecVersion,
		Triggeredid:        "my-triggered-id",
		Type:               strutils.Stringp("sh.keptn.event.dev.delivery.triggered"),
	}

	ce := cloudevents.NewEvent()
	ce.SetType("sh.keptn.event.dev.delivery.triggered")
	ce.SetID("my-id")
	ce.SetSource("my-source")
	ce.SetDataContentType(contentType)
	ce.SetSpecVersion(defaultSpecVersion)
	ce.SetData(cloudevents.ApplicationJSON, TestData{Content: "testdata"})
	ce.SetExtension(keptnContextExtension, "my-keptn-context")
	ce.SetExtension(keptnTriggeredIdExtension, "my-triggered-id")
	ce.SetExtension(keptnSpecVersionExtension, defaultKeptnSpecVersion)

	keptnEvent, err := ToKeptnEvent(ce)

	require.Nil(t, err)
	assert.Equal(t, expected, keptnEvent)
}
