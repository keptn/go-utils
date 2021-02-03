package v0_2_0

import (
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

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
			}
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	type args struct {
		event cloudevents.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "",
			fields: getKeptnFields(ts),
			args: args{
				event: getTestEvent(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpSender, _ := NewHTTPEventSender(ts.URL)
			k := &Keptn{
				KeptnBase: keptn.KeptnBase{
					KeptnContext:       tt.fields.KeptnContext,
					Event:              tt.fields.KeptnBase,
					EventSender:        httpSender,
					UseLocalFileSystem: tt.fields.useLocalFileSystem,
					ResourceHandler:    tt.fields.resourceHandler,
					EventHandler:       tt.fields.eventHandler,
				},
			}
			if err := k.SendCloudEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("SendCloudEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getTestEvent() cloudevents.Event {

	event := cloudevents.NewEvent()
	event.SetType("test-type")
	event.SetSource("test-source")
	return event
}

func TestKeptn_SendCloudEvent(t *testing.T) {
	failOnFirstTry := true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			if failOnFirstTry {
				failOnFirstTry = false
				w.WriteHeader(500)
				w.Write([]byte(`{}`))
			}
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	eventNew := cloudevents.NewEvent()
	eventNew.SetSource("https://test-source")
	eventNew.SetID("8039eac3-9fb2-454f-8b2e-77f8310a81f1")
	eventNew.SetType("sh.keptn.events.test")
	eventNew.SetExtension("shkeptncontext", "test-context")
	eventNew.SetData(cloudevents.ApplicationJSON, map[string]string{"project": "sockshop"})

	httpSender, _ := NewHTTPEventSender(ts.URL)
	k := Keptn{
		KeptnBase: keptn.KeptnBase{
			EventSender: httpSender,
		},
	}

	if err := k.SendCloudEvent(eventNew); err != nil {
		t.Errorf("SendCloudEvent() error = %v", err)
	}
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
		Deployment: DeploymentWithStrategy{
			DeploymentStrategy: "direct",
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
