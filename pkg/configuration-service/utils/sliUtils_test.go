package utils

import (
	"testing"

	"github.com/keptn/go-utils/pkg/configuration-service/models"
)

// TestAddResourceContentToSLIMap
func TestAddResourceContentToSLIMap(t *testing.T) {
	SLIs := make(map[string]string)
	resource := &models.Resource{}
	resourceURI := "dynatrace/sli.yaml"
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
	resourceURI := "dynatrace/sli.yaml"
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
