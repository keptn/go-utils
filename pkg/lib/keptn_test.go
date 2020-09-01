package keptn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/go-test/deep"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
)

func TestNewKeptn(t *testing.T) {
	incomingEvent := cloudevents.New(cloudevents.CloudEventsVersionV02)
	incomingEvent.SetSource("test")
	incomingEvent.SetExtension("shkeptncontext", "test-context")
	incomingEvent.SetDataContentType(cloudevents.ApplicationCloudEventsJSON)
	incomingEvent.SetID("test-id")

	keptnBase := &KeptnBaseEvent{
		Project:            "sockshop",
		Stage:              "dev",
		Service:            "carts",
		TestStrategy:       nil,
		DeploymentStrategy: nil,
		Tag:                nil,
		Image:              nil,
		Labels:             nil,
	}

	marshal, _ := json.Marshal(keptnBase)
	incomingEvent.Data = marshal

	incomingEvent.SetData(marshal)
	incomingEvent.DataEncoded = true
	incomingEvent.DataBinary = true

	type args struct {
		incomingEvent *cloudevents.Event
		opts          keptn.KeptnOpts
	}
	tests := []struct {
		name string
		args args
		want *Keptn
	}{
		{
			name: "Get 'in-cluster' KeptnBase",
			args: args{
				incomingEvent: &incomingEvent,
				opts:          keptn.KeptnOpts{},
			},
			want: &Keptn{
				KeptnBase: keptn.KeptnBase{
					KeptnContext:       "test-context",
					Event:              keptnBase,
					Logger:             keptn.NewLogger("test-context", "test-id", "keptn"),
					EventBrokerURL:     keptn.DefaultEventBrokerURL,
					UseLocalFileSystem: false,
					ResourceHandler:    api.NewResourceHandler(keptn.ConfigurationServiceURL),
					EventHandler:       api.NewEventHandler(keptn.ConfigurationServiceURL),
				},
			},
		},
		{
			name: "Get local KeptnBase",
			args: args{
				incomingEvent: &incomingEvent,
				opts: keptn.KeptnOpts{
					UseLocalFileSystem:      true,
					ConfigurationServiceURL: "",
					EventBrokerURL:          "",
				},
			},
			want: &Keptn{
				KeptnBase: keptn.KeptnBase{
					KeptnContext: "test-context",
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
					Logger:             keptn.NewLogger("test-context", "test-id", "keptn"),
					EventBrokerURL:     keptn.DefaultEventBrokerURL,
					UseLocalFileSystem: true,
					ResourceHandler:    api.NewResourceHandler(keptn.ConfigurationServiceURL),
					EventHandler:       api.NewEventHandler(keptn.ConfigurationServiceURL),
				},
			},
		},
		{
			name: "Get KeptnBase with custom configuration service URL",
			args: args{
				incomingEvent: &incomingEvent,
				opts: keptn.KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "",
				},
			},
			want: &Keptn{
				KeptnBase: keptn.KeptnBase{
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
					KeptnContext:       "test-context",
					EventBrokerURL:     keptn.DefaultEventBrokerURL,
					UseLocalFileSystem: false,
					ResourceHandler:    api.NewResourceHandler("custom-config:8080"),
					EventHandler:       api.NewEventHandler("custom-config:8080"),
					Logger:             keptn.NewLogger("test-context", "test-id", "keptn"),
				},
			},
		},
		{
			name: "Get KeptnBase with custom event brokerURL",
			args: args{
				incomingEvent: &incomingEvent,
				opts: keptn.KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "custom-eb:8080",
				},
			},
			want: &Keptn{
				KeptnBase: keptn.KeptnBase{
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
					KeptnContext:       "test-context",
					EventBrokerURL:     "custom-eb:8080",
					UseLocalFileSystem: false,
					ResourceHandler:    api.NewResourceHandler("custom-config:8080"),
					EventHandler:       api.NewEventHandler("custom-config:8080"),
					Logger:             keptn.NewLogger("test-context", "test-id", "keptn"),
				},
			},
		},
		{
			name: "Get KeptnBase with custom logger",
			args: args{
				incomingEvent: &incomingEvent,
				opts: keptn.KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "custom-eb:8080",
					LoggingOptions: &keptn.LoggingOpts{
						ServiceName: stringp("my-service"),
					},
				},
			},
			want: &Keptn{
				KeptnBase: keptn.KeptnBase{
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
					KeptnContext:       "test-context",
					EventBrokerURL:     "custom-eb:8080",
					UseLocalFileSystem: false,
					ResourceHandler:    api.NewResourceHandler("custom-config:8080"),
					EventHandler:       api.NewEventHandler("custom-config:8080"),
					Logger:             keptn.NewLogger("test-context", "test-id", "my-service"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := NewKeptn(tt.args.incomingEvent, tt.args.opts); deep.Equal(got, tt.want) != nil {
				fmt.Println(deep.Equal(got.Event, tt.want.Event))
				fmt.Println(deep.Equal(got.Logger, tt.want.Logger))
				fmt.Println(deep.Equal(got.ResourceHandler, tt.want.ResourceHandler))
				fmt.Println(deep.Equal(got.EventHandler, tt.want.EventHandler))
				t.Errorf("NewKeptn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeptn_GetKeptnResource(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)

			res := &models.Resource{
				ResourceContent: "dGVzdC1jb250ZW50Cg==",
				ResourceURI:     stringp("test-resource.file"),
			}
			marshal, _ := json.Marshal(res)
			w.Write(marshal)
		}),
	)
	defer ts.Close()

	type fields struct {
		KeptnBase          *KeptnBaseEvent
		KeptnContext       string
		eventBrokerURL     string
		useLocalFileSystem bool
		resourceHandler    *api.ResourceHandler
	}
	type args struct {
		resource string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get a resource",
			fields: fields{
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
				eventBrokerURL:     "",
				useLocalFileSystem: false,
				resourceHandler:    api.NewResourceHandler(ts.URL),
			},
			args: args{
				resource: "test-resource.file",
			},
			want:    "test-content",
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
				},
			}
			got, err := k.GetKeptnResource(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKeptnResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKeptnResource() got = %v, want %v", got, tt.want)
			}
			_ = os.RemoveAll(tt.args.resource)
		})
	}
}

func stringp(s string) *string {
	return &s
}
