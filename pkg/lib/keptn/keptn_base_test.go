package keptn

import (
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
)

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
		t.Errorf("Unexpected length of SLI map")
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

// TestAddEmptyIndicatorsResourceContentToSLIMap
func TestAddEmptyIndicatorsResourceContentToSLIMap(t *testing.T) {
	SLIs := make(map[string]string)
	resource := &models.Resource{}
	resourceURI := "provider/sli.yaml"
	resource.ResourceURI = &resourceURI
	resource.ResourceContent = `---
indicators:
`
	SLIs, err := addResourceContentToSLIMap(SLIs, resource)
	if SLIs != nil {
		t.Errorf("Unexpected length of SLI map")
	}
	if err == nil || !strings.Contains(err.Error(), "missing required field: indicators") {
		t.Errorf("Unexpected error message")
	}
}

func TestGetServiceEndpoint(t *testing.T) {
	type args struct {
		service string
	}
	tests := []struct {
		name        string
		args        args
		envVarValue string
		want        url.URL
		wantErr     bool
	}{
		{
			name: "get http endpoint",
			args: args{
				service: "CONFIGURATION_SERVICE",
			},
			envVarValue: "http://resource-service",
			want: url.URL{
				Scheme: "http",
				Host:   "resource-service",
			},
			wantErr: false,
		},
		{
			name: "get https endpoint",
			args: args{
				service: "CONFIGURATION_SERVICE",
			},
			envVarValue: "https://resource-service",
			want: url.URL{
				Scheme: "https",
				Host:   "resource-service",
			},
			wantErr: false,
		},
		{
			name: "get http endpoint from service-name only",
			args: args{
				service: "CONFIGURATION_SERVICE",
			},
			envVarValue: "resource-service",
			want: url.URL{
				Scheme: "http",
				Host:   "resource-service",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.args.service, tt.envVarValue)
			got, err := GetServiceEndpoint(tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServiceEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServiceEndpoint() got = %v, want %v", got.Host, tt.want.Host)
			}
		})
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
			got := GetExpBackoffTime(tt.args.retryNr)
			if got < tt.want.min || got > tt.want.max {
				t.Errorf("getExpBackoffTime() = %v, want [%v,%v]", got, tt.want.min, tt.want.max)
			}
		})
	}
}
