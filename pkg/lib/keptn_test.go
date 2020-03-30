package keptn

import (
	"encoding/json"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewKeptn(t *testing.T) {
	incomingEvent := cloudevents.New(cloudevents.CloudEventsVersionV02)
	incomingEvent.SetSource("test")
	incomingEvent.SetExtension("shkeptncontext", "test-context")
	incomingEvent.SetDataContentType(cloudevents.ApplicationCloudEventsJSON)

	keptnBase := &KeptnBase{
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
		opts          KeptnOpts
	}
	tests := []struct {
		name string
		args args
		want *Keptn
	}{
		{
			name: "Get 'in-cluster' Keptn",
			args: args{
				incomingEvent: &incomingEvent,
				opts:          KeptnOpts{},
			},
			want: &Keptn{
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
				KeptnContext:       "test-context",
				eventBrokerURL:     defaultEventBrokerURL,
				useLocalFileSystem: false,
				resourceHandler: &api.ResourceHandler{
					BaseURL:    configurationServiceURL,
					AuthHeader: "",
					AuthToken:  "",
					HTTPClient: &http.Client{},
					Scheme:     "http",
				},
			},
		},
		{
			name: "Get local Keptn",
			args: args{
				incomingEvent: &incomingEvent,
				opts: KeptnOpts{
					UseLocalFileSystem:      true,
					ConfigurationServiceURL: "",
					EventBrokerURL:          "",
				},
			},
			want: &Keptn{
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
				KeptnContext:       "test-context",
				eventBrokerURL:     defaultEventBrokerURL,
				useLocalFileSystem: true,
				resourceHandler: &api.ResourceHandler{
					BaseURL:    configurationServiceURL,
					AuthHeader: "",
					AuthToken:  "",
					HTTPClient: &http.Client{},
					Scheme:     "http",
				},
			},
		},
		{
			name: "Get Keptn with custom configuration service URL",
			args: args{
				incomingEvent: &incomingEvent,
				opts: KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "",
				},
			},
			want: &Keptn{
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
				KeptnContext:       "test-context",
				eventBrokerURL:     defaultEventBrokerURL,
				useLocalFileSystem: false,
				resourceHandler: &api.ResourceHandler{
					BaseURL:    "custom-config:8080",
					AuthHeader: "",
					AuthToken:  "",
					HTTPClient: &http.Client{},
					Scheme:     "http",
				},
			},
		},
		{
			name: "Get Keptn with custom event brokerURL",
			args: args{
				incomingEvent: &incomingEvent,
				opts: KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "custom-eb:8080",
				},
			},
			want: &Keptn{
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
				KeptnContext:       "test-context",
				eventBrokerURL:     "custom-eb:8080",
				useLocalFileSystem: false,
				resourceHandler: &api.ResourceHandler{
					BaseURL:    "custom-config:8080",
					AuthHeader: "",
					AuthToken:  "",
					HTTPClient: &http.Client{},
					Scheme:     "http",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := NewKeptn(tt.args.incomingEvent, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
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
		KeptnBase          *KeptnBase
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
				KeptnBase:          tt.fields.KeptnBase,
				KeptnContext:       tt.fields.KeptnContext,
				eventBrokerURL:     tt.fields.eventBrokerURL,
				useLocalFileSystem: tt.fields.useLocalFileSystem,
				resourceHandler:    tt.fields.resourceHandler,
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

// generateStringWithSpecialChars generates a string of the given length
// and containing at least one special character and digit.
func generateStringWithSpecialChars(length int) string {
	rand.Seed(time.Now().UnixNano())

	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials

	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]

	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}

	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	str := string(buf)

	return str
}

// TestInvalidKeptnEntityName tests whether a random string containing a special character or digit
// does not pass the name validation.
func TestInvalidKeptnEntityName(t *testing.T) {
	invalidName := generateStringWithSpecialChars(8)
	if ValidateKeptnEntityName(invalidName) {
		t.Fatalf("%s starts with upper case letter(s) or contains special character(s), but passed the name validation", invalidName)
	}
}

func TestInvalidKeptnEntityName2(t *testing.T) {
	if ValidateKeptnEntityName("sockshop-") {
		t.Fatalf("project name must not end with hyphen")
	}
}

func TestValidKeptnEntityName(t *testing.T) {
	if !ValidateKeptnEntityName("sockshop-test") {
		t.Fatalf("project should be valid")
	}
}

// TestAddResourceContentToSLIMap
func TestAddResourceContentToSLIMap(t *testing.T) {
	SLIs := make(map[string]string)
	resource := &models.Resource{}
	resourceURI := "provider/sli.yaml"
	resource.ResourceURI = &resourceURI
	resource.ResourceContent = `--- 
indicators: 
  error_rate: "builtin:service.errors.total.count:merge(0):avg?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  response_time_p50: "builtin:service.response.time:merge(0):percentile(50)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  response_time_p90: "builtin:service.response.time:merge(0):percentile(90)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  response_time_p95: "builtin:service.response.time:merge(0):percentile(95)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  throughput: "builtin:service.requestCount.total:merge(0):count?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"	
`
	SLIs, _ = addResourceContentToSLIMap(SLIs, resource)

	if len(SLIs) != 5 {
		t.Errorf("Unexpected lenght of SLI map")
	}
}

// TestAddResourceContentToSLIMap
func TestAddMultipleResourceContentToSLIMap(t *testing.T) {
	SLIs := make(map[string]string)
	resource := &models.Resource{}
	resourceURI := "provider/sli.yaml"
	resource.ResourceURI = &resourceURI
	resource.ResourceContent = `--- 
indicators: 
  error_rate: "not defined"
  response_time_p50: "builtin:service.response.time:merge(0):percentile(50)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  response_time_p90: "builtin:service.response.time:merge(0):percentile(90)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  response_time_p95: "builtin:service.response.time:merge(0):percentile(95)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  throughput: "builtin:service.requestCount.total:merge(0):count?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"	
`
	SLIs, _ = addResourceContentToSLIMap(SLIs, resource)

	resource.ResourceContent = `--- 
indicators: 
  error_rate: "builtin:service.errors.total.count:merge(0):avg?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  failure_rate: "builtin:service.requestCount.total:merge(0):count?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"	
`
	SLIs, _ = addResourceContentToSLIMap(SLIs, resource)

	if len(SLIs) != 6 {
		t.Errorf("Unexpected length of SLI map")
	}

	if SLIs["error_rate"] != "builtin:service.errors.total.count:merge(0):avg?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)" {
		t.Errorf("Unexpected value of error_rate SLI")
	}
}

func stringp(s string) *string {
	return &s
}
