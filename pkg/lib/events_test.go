package keptn

import (
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type fields struct {
	KeptnContext       string
	KeptnBase          *KeptnBase
	eventBrokerURL     string
	useLocalFileSystem bool
	resourceHandler    *api.ResourceHandler
	eventHandler       *api.EventHandler
}

func getKeptnFields(ts *httptest.Server) fields {
	return fields{
		KeptnBase: &KeptnBase{
			Project:            "sockshop",
			Stage:              "dev",
			Service:            "carts",
			TestStrategy:       nil,
			DeploymentStrategy: nil,
			Tag:                nil,
			Image:              nil,
			Labels:             nil,
		},
		eventBrokerURL:     ts.URL,
		useLocalFileSystem: false,
		resourceHandler:    api.NewResourceHandler(ts.URL),
	}
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

	source, _ := url.Parse("https://test-source")
	contentType := "application/json"

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
				event: cloudevents.Event{
					Context: cloudevents.EventContextV02{
						ID:          uuid.New().String(),
						Time:        &types.Timestamp{Time: time.Now()},
						Type:        ConfigurationChangeEventType,
						Source:      types.URLRef{URL: *source},
						ContentType: &contentType,
						Extensions:  map[string]interface{}{"shkeptncontext": "test-context"},
					}.AsV02(),
					Data: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Keptn{
				KeptnContext:       tt.fields.KeptnContext,
				KeptnBase:          tt.fields.KeptnBase,
				eventBrokerURL:     tt.fields.eventBrokerURL,
				useLocalFileSystem: tt.fields.useLocalFileSystem,
				resourceHandler:    tt.fields.resourceHandler,
				eventHandler:       tt.fields.eventHandler,
			}
			if err := k.SendCloudEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("SendCloudEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeptn_SendFunctions(t *testing.T) {

	receivedCorrectType := make(chan bool)
	var wantEventType string
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)

			payload, _ := ioutil.ReadAll(r.Body)

			ce := &cloudevents.Event{}

			_ = json.Unmarshal(payload, &ce)
			fmt.Println(ce.Context.GetType())
			receivedCorrectType <- ce.Context.GetType() == wantEventType
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	k := &Keptn{
		KeptnBase: &KeptnBase{
			Project:            "sockshop",
			Stage:              "dev",
			Service:            "carts",
			TestStrategy:       nil,
			DeploymentStrategy: nil,
			Tag:                nil,
			Image:              nil,
			Labels:             nil,
		},
		eventBrokerURL:     ts.URL,
		useLocalFileSystem: false,
		resourceHandler:    api.NewResourceHandler(ts.URL),
	}

	type args struct {
		incomingEvent      *cloudevents.Event
		teststrategy       string
		deploymentstrategy string
		startedAt          time.Time
		result             string
		labels             map[string]string
		eventSource        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				incomingEvent:      nil,
				teststrategy:       "functional",
				deploymentstrategy: "direct",
				startedAt:          time.Time{},
				result:             "pass",
				labels:             nil,
				eventSource:        "test-service",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			wantEventType = TestsFinishedEventType
			go func() {
				if err := k.SendTestsFinishedEvent(tt.args.incomingEvent, tt.args.teststrategy, tt.args.deploymentstrategy, tt.args.startedAt, tt.args.result, tt.args.labels, tt.args.eventSource); (err != nil) != tt.wantErr {
					t.Errorf("SendTestsFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			verifyReceivedEventType(receivedCorrectType, t)

			wantEventType = ConfigurationChangeEventType
			go func() {
				if err := k.SendConfigurationChangeEvent(tt.args.incomingEvent, tt.args.labels, tt.args.eventSource); (err != nil) != tt.wantErr {
					t.Errorf("SendTestsFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			verifyReceivedEventType(receivedCorrectType, t)

			wantEventType = DeploymentFinishedEventType
			go func() {
				if err := k.SendDeploymentFinishedEvent(tt.args.incomingEvent, "functional", "direct", "", "", "", "", tt.args.labels, tt.args.eventSource); (err != nil) != tt.wantErr {
					t.Errorf("SendTestsFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			verifyReceivedEventType(receivedCorrectType, t)
		})
	}
}

func verifyReceivedEventType(receivedCorrectType chan bool, t *testing.T) {
	select {
	case success := <-receivedCorrectType:
		if !success {
			t.Errorf("SendTestsFinishedEvent(): did not receive correct cloud event type")
		}

	case <-time.After(5 * time.Second):
		t.Errorf("SendTestsFinishedEvent(): timed out waiting for event")
	}
}

func Test_getExpBackoffTime(t *testing.T) {
	type args struct {
		retryNr int
	}
	type durationRange struct {
		min time.Duration
		max time.Duration
	}
	tests := []struct {
		name string
		args args
		want durationRange
	}{
		{
			name: "Get exponential backoff time (1)",
			args: args{
				retryNr: 1,
			},
			want: durationRange{
				min: 375.0 * time.Millisecond,
				max: 1125.0 * time.Millisecond,
			},
		},
		{
			name: "Get exponential backoff time (2)",
			args: args{
				retryNr: 2,
			},
			want: durationRange{
				min: 750.0 * time.Millisecond,
				max: 2250.0 * time.Millisecond,
			},
		},
		{
			name: "Get exponential backoff time (3)",
			args: args{
				retryNr: 3,
			},
			want: durationRange{
				min: 1125.0 * time.Millisecond,
				max: 3375.0 * time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getExpBackoffTime(tt.args.retryNr)
			if got < tt.want.min || got > tt.want.max {
				t.Errorf("getExpBackoffTime() = %v, want [%v,%v]", got, tt.want.min, tt.want.max)
			}
		})
	}
}

func TestEventTypeExists(t *testing.T) {
	tests := []struct{
		InputTyp string
		Result	bool
	} {
		{"sh.keptn.internal.event.project.create", true},
		{"sh.keptn.internal.event.project.delete", true},
		{"sh.keptn.internal.event.service.create", true},
		{"sh.keptn.events.tests-finished", true},
		{"sh.keptn.events.problem", true},
		{"sh.keptn.internal.event.get-sli", true},
		{"sh.kepntt.internal.event.get-sli", false},
		{"event.get-sli", false},
		{"internal.event.get-sli", false},
		{"create", false},
		{"", false},
	}

	for _, tt := range tests {
		got := EventTypeExists(tt.InputTyp)
		if got != tt.Result {
			t.Errorf("EventTypeExists() = %t, wanted %t", got, tt.Result)
		}
	}
}