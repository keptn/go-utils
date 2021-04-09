package eventutils

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateKeptnEvent_MissingInformation(t *testing.T) {
	type TestData struct {
		v0_2_0.EventData
		Content string `json:"content"`
	}

	t.Run("missing project", func(t *testing.T) {
		testData := TestData{
			EventData: v0_2_0.EventData{
				Stage:   "my-stage",
				Service: "my-service",
			},
		}
		_, err := KeptnEvent("sh.keptn.event.dev.delivery.triggered", testData).Build()
		assert.NotNil(t, err)
	})

	t.Run("missing stage", func(t *testing.T) {
		testData := TestData{
			EventData: v0_2_0.EventData{
				Project: "my-project",
				Service: "my-service",
			},
		}
		_, err := KeptnEvent("sh.keptn.event.dev.delivery.triggered", testData).Build()
		assert.NotNil(t, err)
	})

	t.Run("missing service", func(t *testing.T) {
		testData := TestData{
			EventData: v0_2_0.EventData{
				Project: "my-project",
				Stage:   "my-stage",
			},
		}
		_, err := KeptnEvent("sh.keptn.event.dev.delivery.triggered", testData).Build()
		assert.NotNil(t, err)
	})
}
func TestCreateSimpleKeptnEvent(t *testing.T) {

	type TestData struct {
		v0_2_0.EventData
		Content string `json:"content"`
	}

	testData := TestData{
		EventData: v0_2_0.EventData{
			Project: "my-project",
			Stage:   "my-stabe",
			Service: "my-service",
		},
		Content: "some-content",
	}

	event, err := KeptnEvent("sh.keptn.event.dev.delivery.triggered", testData).Build()
	require.Nil(t, err)
	require.Equal(t, "application/json", event.Contenttype)
	require.Equal(t, testData, event.Data)
	require.Equal(t, "", event.Shkeptncontext)
	require.Equal(t, time.Now().UTC().Round(time.Minute), time.Time(event.Time).Round(time.Minute))
	require.Equal(t, defaultKeptnSpecVersion, event.Shkeptnspecversion)
	require.Equal(t, defaultSpecVersion, event.Specversion)
	require.Equal(t, "", event.Triggeredid)
	require.Equal(t, strutils.Stringp("sh.keptn.event.dev.delivery.triggered"), event.Type)
}

func TestCreateKeptnEvent(t *testing.T) {

	event, _ := KeptnEvent("sh.keptn.event.dev.delivery.triggered", map[string]interface{}{}).
		WithID("my-id").
		WithKeptnContext("my-keptn-context").
		WithSource("my-source").
		WithTriggeredID("my-triggered-id").
		WithKeptnSpecVersion("2.0").
		Build()

	require.Equal(t, "application/json", event.Contenttype)
	require.Equal(t, map[string]interface{}{}, event.Data)
	require.Equal(t, "2.0", event.Shkeptnspecversion)
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
