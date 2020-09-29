package keptn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
)

type fields struct {
	KeptnContext       string
	KeptnBase          *KeptnBaseEvent
	eventBrokerURL     string
	useLocalFileSystem bool
	resourceHandler    *api.ResourceHandler
	eventHandler       *api.EventHandler
}

func getKeptnFields(ts *httptest.Server) fields {
	return fields{
		KeptnBase: &KeptnBaseEvent{
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
				KeptnBase: keptn.KeptnBase{
					KeptnContext:       tt.fields.KeptnContext,
					Event:              tt.fields.KeptnBase,
					EventBrokerURL:     tt.fields.eventBrokerURL,
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
		KeptnBase: keptn.KeptnBase{
			KeptnContext: "",
			Event: &KeptnBaseEvent{
				Project:            "sockshop",
				Stage:              "dev",
				Service:            "carts",
				TestStrategy:       nil,
				DeploymentStrategy: nil,
				Tag:                nil,
				Image:              nil,
				Labels:             nil,
			},
			EventBrokerURL:     ts.URL,
			UseLocalFileSystem: false,
			ResourceHandler:    api.NewResourceHandler(ts.URL),
		},
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

			wantEventType = InternalGetSLIDoneEventType
			go func() {
				if err := k.SendInternalGetSLIDoneEvent(tt.args.incomingEvent, nil, nil, tt.args.labels, nil, tt.args.eventSource); (err != nil) != tt.wantErr {
					t.Errorf("SendInternalGetSLIDoneEvent() error = %v, wantErr %v", err, tt.wantErr)
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
