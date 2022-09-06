package v0_2_0

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/keptn/go-utils/config"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type fields struct {
	KeptnContext       string
	KeptnBase          *EventData
	eventBrokerURL     string
	useLocalFileSystem bool
	resourceHandler    *api.ResourceHandler
	eventHandler       *api.EventHandler
}

func getKeptnFields(ts *httptest.Server) fields {
	return fields{
		KeptnBase: &EventData{
			Project: "sockshop",
			Stage:   "dev",
			Service: "carts",
		},
		eventBrokerURL:     ts.URL,
		useLocalFileSystem: false,
		resourceHandler:    api.NewResourceHandler(ts.URL),
	}
}

func TestKeptn_SendCloudEventWithRetry(t *testing.T) {
	failOnFirstTry := true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			if failOnFirstTry {
				failOnFirstTry = false
				w.WriteHeader(500)
				w.Write([]byte(`{}`))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	httpSender, _ := NewHTTPEventSender(ts.URL)

	err := httpSender.Send(context.TODO(), getTestEvent())

	require.Nil(t, err)
}

func TestKeptn_SendCloudEventWithOneRetry(t *testing.T) {
	nrRequests := 0
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			if nrRequests < 2 {
				nrRequests++
				w.WriteHeader(500)
				w.Write([]byte(`{}`))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	httpSender, _ := NewHTTPEventSender(ts.URL, WithSendRetries(1))

	err := httpSender.Send(context.TODO(), getTestEvent())

	require.NotNil(t, err)

	// reset nrRequests
	nrRequests = 0

	//initialize a new sender with 2 retries
	httpSender, _ = NewHTTPEventSender(ts.URL, WithSendRetries(5))

	err = httpSender.Send(context.TODO(), getTestEvent())

	require.Nil(t, err)
}

func getTestEvent() cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetType("test-type")
	event.SetSource("test-source")
	return event
}

func TestKeptn_SendCloudEvent(t *testing.T) {
	testEventSender := &TestSender{}
	eventNew := cloudevents.NewEvent()
	eventNew.SetSource("https://test-source")
	eventNew.SetID("8039eac3-9fb2-454f-8b2e-77f8310a81f1")
	eventNew.SetType("sh.keptn.events.test")
	eventNew.SetExtension("shkeptncontext", "test-context")
	eventNew.SetData(cloudevents.ApplicationJSON, map[string]string{"project": "sockshop"})

	k := Keptn{
		KeptnBase: keptn.KeptnBase{
			EventSender: testEventSender,
		},
	}

	err := k.SendCloudEvent(eventNew)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(testEventSender.SentEvents))
	keptnEvent, err := ToKeptnEvent(testEventSender.SentEvents[0])
	assert.Nil(t, err)
	assert.Equal(t, config.GetKeptnGoUtilsConfig().ShKeptnSpecVersion, keptnEvent.Shkeptnspecversion)
	assert.Equal(t, defaultSpecVersion, keptnEvent.Specversion)
	assert.Equal(t, "sh.keptn.events.test", *keptnEvent.Type)
	assert.Equal(t, "test-context", keptnEvent.Shkeptncontext)
	assert.Equal(t, "8039eac3-9fb2-454f-8b2e-77f8310a81f1", keptnEvent.ID)
	assert.Equal(t, "https://test-source", *keptnEvent.Source)
	assert.Equal(t, map[string]interface{}{"project": "sockshop"}, keptnEvent.Data)
}

func TestEventDataAs(t *testing.T) {
	eventData := DeploymentTriggeredEventData{
		EventData: EventData{
			Project: "p",
			Stage:   "s",
			Service: "s",
			Labels:  map[string]string{"1": "2"},
			Status:  StatusSucceeded,
			Result:  ResultPass,
			Message: "m",
		},
		ConfigurationChange: ConfigurationChange{
			Values: map[string]interface{}{"image": "my-image:tag"},
		},
	}

	ce := models.KeptnContextExtendedCE{
		Data: eventData,
	}

	var decodedEventData DeploymentTriggeredEventData
	err := EventDataAs(ce, &decodedEventData)
	assert.Nil(t, err)
	assert.Equal(t, eventData, decodedEventData)
}

func TestGetEventTypeForTriggeredEvent(t *testing.T) {
	type args struct {
		baseTriggeredEventType string
		newEventTypeSuffix     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get .started event",
			args: args{
				baseTriggeredEventType: GetTriggeredEventType(EvaluationTaskName),
				newEventTypeSuffix:     keptnStartedEventSuffix,
			},
			want:    GetStartedEventType(EvaluationTaskName),
			wantErr: false,
		},
		{
			name: "get .status.changed event",
			args: args{
				baseTriggeredEventType: GetTriggeredEventType(EvaluationTaskName),
				newEventTypeSuffix:     keptnStatusChangedEventSuffix,
			},
			want:    GetStatusChangedEventType(EvaluationTaskName),
			wantErr: false,
		},
		{
			name: "get .finished event",
			args: args{
				baseTriggeredEventType: GetTriggeredEventType(EvaluationTaskName),
				newEventTypeSuffix:     keptnFinishedEventSuffix,
			},
			want:    GetFinishedEventType(EvaluationTaskName),
			wantErr: false,
		},
		{
			name: "no .triggered event as input - return error",
			args: args{
				baseTriggeredEventType: GetStartedEventType(EvaluationTaskName),
				newEventTypeSuffix:     keptnFinishedEventSuffix,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEventTypeForTriggeredEvent(tt.args.baseTriggeredEventType, tt.args.newEventTypeSuffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventTypeForTriggeredEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEventTypeForTriggeredEvent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateKeptnEvent_MissingInformation(t *testing.T) {
	type TestData struct {
		EventData
		Content string `json:"content"`
	}

	t.Run("missing project", func(t *testing.T) {
		testData := TestData{
			EventData: EventData{
				Stage:   "my-stage",
				Service: "my-service",
			},
		}
		_, err := KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
		assert.NotNil(t, err)
	})

	t.Run("missing stage", func(t *testing.T) {
		testData := TestData{
			EventData: EventData{
				Project: "my-project",
				Service: "my-service",
			},
		}
		_, err := KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
		assert.NotNil(t, err)
	})

	t.Run("missing service", func(t *testing.T) {
		testData := TestData{
			EventData: EventData{
				Project: "my-project",
				Stage:   "my-stage",
			},
		}
		_, err := KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
		assert.NotNil(t, err)
	})
}
func TestCreateSimpleKeptnEvent(t *testing.T) {

	type TestData struct {
		EventData
		Content string `json:"content"`
	}

	testData := TestData{
		EventData: EventData{
			Project: "my-project",
			Stage:   "my-stabe",
			Service: "my-service",
		},
		Content: "some-content",
	}

	event, err := KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
	require.Nil(t, err)
	require.Equal(t, "application/json", event.Contenttype)
	require.Equal(t, testData, event.Data)
	require.Equal(t, "", event.Shkeptncontext)
	require.Equal(t, time.Now().UTC().Round(time.Minute), event.Time.Round(time.Minute))
	require.NotEmpty(t, event.Shkeptnspecversion)
	require.Equal(t, defaultSpecVersion, event.Specversion)
	require.Equal(t, "", event.Triggeredid)
	require.Equal(t, "", event.GitCommitID)
	require.Equal(t, strutils.Stringp("sh.keptn.event.dev.delivery.triggered"), event.Type)
}

func TestCreateKeptnEvent(t *testing.T) {

	event, _ := KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", map[string]interface{}{}).
		WithID("my-id").
		WithKeptnContext("my-keptn-context").
		WithTriggeredID("my-triggered-id").
		WithGitCommitID("my-commit-id").
		WithKeptnSpecVersion("2.0").
		Build()

	require.Equal(t, "application/json", event.Contenttype)
	require.Equal(t, map[string]interface{}{}, event.Data)
	require.Equal(t, "2.0", event.Shkeptnspecversion)
	require.Equal(t, defaultSpecVersion, event.Specversion)
	require.Equal(t, "my-id", event.ID)
	require.Equal(t, time.Now().UTC().Round(time.Minute), event.Time.Round(time.Minute))
	require.Equal(t, strutils.Stringp("source"), event.Source)
	require.Equal(t, "my-keptn-context", event.Shkeptncontext)
	require.Equal(t, "my-triggered-id", event.Triggeredid)
	require.Equal(t, "my-commit-id", event.GitCommitID)
	require.Equal(t, strutils.Stringp("sh.keptn.event.dev.delivery.triggered"), event.Type)
}

func TestToCloudEvent(t *testing.T) {

	type TestData struct {
		Content string `json:"content"`
	}

	time := time.Now().UTC()

	expected := cloudevents.NewEvent()
	expected.SetType("sh.keptn.event.dev.delivery.triggered")
	expected.SetID("my-id")
	expected.SetSource("source")
	expected.SetTime(time)
	expected.SetData(cloudevents.ApplicationJSON, TestData{Content: "testdata"})
	expected.SetDataContentType(cloudevents.ApplicationJSON)
	expected.SetSpecVersion(defaultSpecVersion)
	expected.SetExtension(keptnContextCEExtension, "my-keptn-context")
	expected.SetExtension(triggeredIDCEExtension, "my-triggered-id")
	expected.SetExtension(keptnGitCommitIDCEExtension, "git-commit-id")
	expected.SetExtension(keptnSpecVersionCEExtension, config.GetKeptnGoUtilsConfig().ShKeptnSpecVersion)

	keptnEvent := models.KeptnContextExtendedCE{
		Contenttype:        "application/json",
		Data:               TestData{Content: "testdata"},
		ID:                 "my-id",
		Shkeptncontext:     "my-keptn-context",
		Source:             strutils.Stringp("source"),
		Shkeptnspecversion: config.GetKeptnGoUtilsConfig().ShKeptnSpecVersion,
		Specversion:        defaultSpecVersion,
		Triggeredid:        "my-triggered-id",
		GitCommitID:        "git-commit-id",
		Type:               strutils.Stringp("sh.keptn.event.dev.delivery.triggered"),
		Time:               time,
	}
	cloudevent := ToCloudEvent(keptnEvent)
	assert.Equal(t, expected, cloudevent)

}

func TestToKeptnEvent(t *testing.T) {

	type TestData struct {
		Content string `json:"content"`
	}

	time := time.Now().UTC()
	expected := models.KeptnContextExtendedCE{
		Contenttype:        "application/json",
		Data:               map[string]interface{}{"content": "testdata"},
		ID:                 "my-id",
		Shkeptncontext:     "my-keptn-context",
		Shkeptnspecversion: config.GetKeptnGoUtilsConfig().ShKeptnSpecVersion,
		Source:             strutils.Stringp("my-source"),
		Specversion:        defaultSpecVersion,
		Time:               time,
		Triggeredid:        "my-triggered-id",
		GitCommitID:        "my-commit-id",
		Type:               strutils.Stringp("sh.keptn.event.dev.delivery.triggered"),
	}

	ce := cloudevents.NewEvent()
	ce.SetType("sh.keptn.event.dev.delivery.triggered")
	ce.SetID("my-id")
	ce.SetSource("my-source")
	ce.SetTime(time)
	ce.SetDataContentType(cloudevents.ApplicationJSON)
	ce.SetSpecVersion(defaultSpecVersion)
	ce.SetData(cloudevents.ApplicationJSON, TestData{Content: "testdata"})
	ce.SetExtension(keptnContextCEExtension, "my-keptn-context")
	ce.SetExtension(triggeredIDCEExtension, "my-triggered-id")
	ce.SetExtension(keptnGitCommitIDCEExtension, "my-commit-id")
	ce.SetExtension(keptnSpecVersionCEExtension, config.GetKeptnGoUtilsConfig().ShKeptnSpecVersion)

	keptnEvent, err := ToKeptnEvent(ce)

	require.Nil(t, err)
	require.Equal(t, expected, keptnEvent)
}

func TestIsTaskEventType(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid", args{"sh.keptn.event.deployment.triggered"}, true},
		{"too long", args{"sh.keptn.event.deployment.triggered.triggered"}, false},
		{"empty", args{""}, false},
		{"only dots", args{"...."}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTaskEventType(tt.args.eventType); got != tt.want {
				t.Errorf("IsTaskEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSequenceEventType(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid", args{"sh.keptn.event.dev.delivery.triggered"}, true},
		{"too long", args{"sh.keptn.event.dev.delivery.triggered.triggered"}, false},
		{"empty", args{""}, false},
		{"only dots", args{"....."}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSequenceEventType(tt.args.eventType); got != tt.want {
				t.Errorf("IsSequenceEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSequenceEventType(t *testing.T) {
	type args struct {
		sequenceTriggeredEventType string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		want2   string
		wantErr bool
	}{
		{"valid", args{"sh.keptn.event.dev.delivery.triggered"}, "dev", "delivery", "triggered", false},
		{"too long", args{"sh.keptn.event.dev.delivery.triggered.triggered"}, "", "", "", true},
		{"empty", args{""}, "", "", "", true},
		{"only dots", args{"....."}, "", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := ParseSequenceEventType(tt.args.sequenceTriggeredEventType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSequenceEventType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseSequenceEventType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseSequenceEventType() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ParseSequenceEventType() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestParseTaskEventType(t *testing.T) {
	type args struct {
		taskEventType string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{"valid", args{"sh.keptn.event.task.triggered"}, "task", "triggered", false},
		{"too long", args{"sh.keptn.event.task.triggered.triggered"}, "", "", true},
		{"empty", args{""}, "", "", true},
		{"only dots", args{"....."}, "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseTaskEventType(tt.args.taskEventType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTaskEventType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTaskEventType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseTaskEventType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParseEventKind(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"valid", args{"sh.keptn.event.task.triggered"}, "triggered", false},
		{"empty", args{""}, "", true},
		{"only dots", args{"....."}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEventKind(tt.args.eventType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEventKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseEventKind() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseEventTypeWithoutKind(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"valid", args{"sh.keptn.event.task.triggered"}, "sh.keptn.event.task", false},
		{"empty", args{""}, "", true},
		{"only dots", args{"....."}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEventTypeWithoutKind(tt.args.eventType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEventTypeWithoutKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseEventTypeWithoutKind() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceEventTypeKind(t *testing.T) {
	type args struct {
		eventType string
		newKind   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"valid", args{"sh.keptn.event.task.triggered", "started"}, "sh.keptn.event.task.started", false},
		{"valid", args{"sh.keptn.event.task.triggered", ""}, "sh.keptn.event.task", false},
		{"empty", args{"", ""}, "", true},
		{"only dots", args{".....", "started"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceEventTypeKind(tt.args.eventType, tt.args.newKind)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceEventTypeKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReplaceEventTypeKind() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CreateFinishedEvent(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
		eventData   interface{}
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing event type",
			args: args{
				source:      "source",
				parentEvent: models.KeptnContextExtendedCE{},
				eventData:   nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "missing event context",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Type: strutils.Stringp("sh.keptn.event.evaluation.triggered"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "event type cannot be replaced",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("something.weird"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "passed event data missing",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "defalult status succeeded",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: EventData{},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return ce.Data.(map[string]interface{})["status"] == string(StatusSucceeded) && ce.Data.(map[string]interface{})["result"] == string(ResultPass)
			},
			wantErr: false,
		},
		{
			name: "defalult status succeeded #2",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: EventData{Status: ""},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return ce.Data.(map[string]interface{})["status"] == string(StatusSucceeded) && ce.Data.(map[string]interface{})["result"] == string(ResultPass)
			},
			wantErr: false,
		},
		{
			name: "defalult result pass",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: EventData{},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return ce.Data.(map[string]interface{})["result"] == "pass"
			},
			wantErr: false,
		},
		{
			name: "defalult result pass #2",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: EventData{Result: ""},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return ce.Data.(map[string]interface{})["result"] == "pass"
			},
			wantErr: false,
		},
		{
			name: "passed status reused",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: EventData{Status: StatusErrored, Result: ResultFailed},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				fmt.Println(ce.Data.(map[string]interface{})["status"])
				return ce.Data.(map[string]interface{})["status"] == string(StatusErrored) && ce.Data.(map[string]interface{})["result"] == string(ResultFailed)
			},
			wantErr: false,
		},
		{
			name: "correct event type",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: EventData{Status: StatusSucceeded, Result: ResultPass},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				fmt.Println(ce.Data.(map[string]interface{})["status"])
				return *ce.Type == "sh.keptn.event.eval.finished"
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateFinishedEvent(tt.args.source, tt.args.parentEvent, tt.args.eventData)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}

func Test_CreateStartedEvent(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
		eventData   interface{}
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing keptn context",
			args: args{
				source:      "source",
				parentEvent: models.KeptnContextExtendedCE{},
				eventData:   nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "missing event type",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abce",
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "non-replacable event type",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abce",
					Type:           strutils.Stringp("somethin.weird"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "passed event data",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Data:           EventData{},
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: EventData{Status: StatusSucceeded, Result: ResultPass},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == "sh.keptn.event.eval.started" &&
					reflect.DeepEqual(ce.Data, EventData{Status: StatusSucceeded, Result: ResultPass})
			},
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == "sh.keptn.event.eval.started"
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateStartedEvent(tt.args.source, tt.args.parentEvent, tt.args.eventData)
			if (err != nil) != tt.wantErr {
				t.Errorf("createStartedEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}

func Test_CreateErrorLogEvent(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
		eventData   interface{}
		errVal      *Error
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing event type",
			args: args{
				source:      "soure",
				parentEvent: models.KeptnContextExtendedCE{},
				eventData:   nil,
				errVal:      nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "missing keptn context",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Type: strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: nil,
				errVal:    nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "creates finished event for events of type .triggered",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
				},
				eventData: nil,
				errVal:    nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return *ce.Type == "sh.keptn.event.eval.finished" },
			wantErr:     false,
		},
		{
			name: "creates error event for events other than .triggered",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
					Data:           ErrorLogEvent{Message: "msg"},
				},
				eventData: nil,
				errVal:    nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == ErrorLogEventName &&
					reflect.DeepEqual(ce.Data, ErrorLogEvent{Task: "eval"})
			},
			wantErr: false,
		},
		{
			name: "creates error event for events other than .triggered with custom err message",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: nil,
				errVal:    &Error{Message: "silly msg"},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == ErrorLogEventName &&
					reflect.DeepEqual(ce.Data, ErrorLogEvent{Message: "silly msg", Task: "eval"})
			},
			wantErr: false,
		},
		{
			name: "creates error event for events other than .triggered with custom eventData",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: ErrorLogEvent{Message: "silly msg2", Task: "eval33"},
				errVal:    nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == ErrorLogEventName &&
					reflect.DeepEqual(ce.Data, ErrorLogEvent{Message: "silly msg2", Task: "eval33"})
			},
			wantErr: false,
		},
		{
			name: "creates error event for events other than .triggered with custom eventData and err",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: ErrorLogEvent{Message: "silly msg2", Task: "eval33"},
				errVal:    &Error{Message: "silly msg"},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == ErrorLogEventName &&
					reflect.DeepEqual(ce.Data, ErrorLogEvent{Message: "silly msg2", Task: "eval33"})
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateErrorLogEvent(tt.args.source, tt.args.parentEvent, tt.args.eventData, tt.args.errVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("createErrorLogEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}

func Test_CreateErrorEvent(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
		eventData   interface{}
		err         *Error
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing keptn context",
			args: args{
				source: "",
				parentEvent: models.KeptnContextExtendedCE{
					Type: strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "missing event type",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abce",
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "creates error event for events other than .triggered",
			args: args{
				source: "",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return *ce.Type == ErrorLogEventName },
			wantErr:     false,
		},
		{
			name: "creates finished event for events of type .triggered",
			args: args{
				source: "",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.triggered"),
					Data:           EventData{Project: "proj"},
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == "sh.keptn.event.eval.finished" &&
					reflect.DeepEqual(ce.Data, EventData{Project: "proj"})
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateErrorEvent(tt.args.source, tt.args.parentEvent, tt.args.eventData, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("createErrorEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}

func Test_CreateFinishedEventWithError(t *testing.T) {
	type args struct {
		source      string
		parentEvent models.KeptnContextExtendedCE
		eventData   interface{}
		err         *Error
	}
	tests := []struct {
		name        string
		args        args
		assertEvent func(*models.KeptnContextExtendedCE) bool
		wantErr     bool
	}{
		{
			name: "missing keptn context",
			args: args{
				source: "",
				parentEvent: models.KeptnContextExtendedCE{
					Type: strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "missing event type",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abce",
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "non-replacable event type",
			args: args{
				source: "source",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abce",
					Type:           strutils.Stringp("somethin.weird"),
				},
				eventData: nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool { return ce == nil },
			wantErr:     true,
		},
		{
			name: "creates error event with custom err",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: nil,
				err:       &Error{Message: "silly msg", StatusType: StatusAborted, ResultType: ResultFailed},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == "sh.keptn.event.eval.finished" &&
					reflect.DeepEqual(ce.Data, EventData{Message: "silly msg", Status: StatusAborted, Result: ResultFailed})
			},
			wantErr: false,
		},
		{
			name: "creates error event with custom eventData",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: EventData{Message: "silly msg", Status: StatusAborted, Result: ResultFailed},
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == "sh.keptn.event.eval.finished" &&
					reflect.DeepEqual(ce.Data, EventData{Message: "silly msg", Status: StatusAborted, Result: ResultFailed})
			},
			wantErr: false,
		},
		{
			name: "creates error event with custom eventData and err",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
				},
				eventData: EventData{Message: "silly msg", Status: StatusAborted, Result: ResultFailed},
				err:       &Error{Message: "silly msg22", StatusType: StatusSucceeded, ResultType: ResultPass},
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				return *ce.Type == "sh.keptn.event.eval.finished" &&
					reflect.DeepEqual(ce.Data, EventData{Message: "silly msg", Status: StatusAborted, Result: ResultFailed})
			},
			wantErr: false,
		},
		{
			name: "creates error event with eventdata and err nil",
			args: args{
				source: "soure",
				parentEvent: models.KeptnContextExtendedCE{
					Shkeptncontext: "abcde",
					Type:           strutils.Stringp("sh.keptn.event.eval.started"),
					Data:           EventData{Project: "proj"},
				},
				eventData: nil,
				err:       nil,
			},
			assertEvent: func(ce *models.KeptnContextExtendedCE) bool {
				fmt.Println(ce.Data)
				return *ce.Type == "sh.keptn.event.eval.finished" &&
					reflect.DeepEqual(ce.Data, EventData{Project: "proj"})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateFinishedEventWithError(tt.args.source, tt.args.parentEvent, tt.args.eventData, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFinishedEventWithError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.assertEvent(got))
		})
	}
}
