package events

import (
	"encoding/json"
)

// InternalProjectCreateEventType is a CloudEvent type for creating a new project
const InternalProjectCreateEventType = "sh.keptn.internal.event.project.create"

// InternalProjectDeleteEventType is a CloudEvent type for deleting a project
const InternalProjectDeleteEventType = "sh.keptn.internal.event.project.delete"

// InternalServiceCreateEventType is a CloudEvent type for creating a new service
const InternalServiceCreateEventType = "sh.keptn.internal.event.service.create"

// ConfigurationChangeEventType is a CloudEvent type for changing the configuration
const ConfigurationChangeEventType = "sh.keptn.event.configuration.change"

// DeploymentFinishedEventType is a CloudEvent for indicating that the deployment has finished
const DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"

// TestsFinishedEventType is a CloudEvent for indicating that tests have finished
const TestsFinishedEventType = "sh.keptn.events.tests-finished"

// StartEvaluationEventType is a CloudEvent for retrieving SLI values
const StartEvaluationEventType = "sh.keptn.event.start-evaluation"

// EvaluationDoneEventType is a CloudEvent for indicating that the evaluation has finished
const EvaluationDoneEventType = "sh.keptn.events.evaluation-done"

// ProblemOpenEventType is a CloudEvent type to inform about an open problem
const ProblemOpenEventType = "sh.keptn.event.problem.open"

// ProblemEventType is a CloudEvent type to inform about a problem
const ProblemEventType = "sh.keptn.events.problem"

// ConfigureMonitoringEventType is a CloudEvent for configuring monitoring
const ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"

// InternalGetSLIEventType is a CloudEvent for retrieving SLI values
const InternalGetSLIEventType = "sh.keptn.internal.event.get-sli"

// InternalGetSLIDoneEventType is a CloudEvent for submitting SLI values
const InternalGetSLIDoneEventType = "sh.keptn.internal.event.get-sli.done"

// ProjectCreateEventData represents the data for creating a new project
type ProjectCreateEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Shipyard is a base64 encoded string of the shipyard file
	Shipyard string `json:"shipyard"`
	// GitUser is the name of a git user of an upstream repository
	GitUser string `json:"gitUser,omitempty"`
	// GitToken is the authentication token for the git user
	GitToken string `json:"gitToken,omitempty"`
	// GitRemoteURL is the remote url of a repository
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`
}

// ProjectDeleteEventData represents the data for deleting a project
type ProjectDeleteEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
}

// ServiceCreateEventData represents the data for creating a new service
type ServiceCreateEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// HelmChart are the data of a Helm chart packed as tgz and base64 encoded
	HelmChart string `json:"helmChart"`
	// DeploymentStrategies contains the deployment strategy for the stages
	DeploymentStrategies map[string]DeploymentStrategy `json:"deploymentStrategies"`
}

// ConfigurationChangeEventData represents the data for changing the service configuration
type ConfigurationChangeEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// ValuesCanary contains new Helm values for canary
	ValuesCanary map[string]interface{} `json:"valuesCanary,omitempty"`
	// Canary contains a new configuration for canary releases
	Canary *Canary `json:"canary,omitempty"`
	// FileChangesUserChart provides new content for the user chart.
	// The key value pairs represent the URI within the chart (i.e. the key) and the new content (i.e. the value).
	FileChangesUserChart map[string]string `json:"fileChangesUserChart,omitempty"`
	// FileChangesGeneratedChart provides new content for the generated chart.
	// The key value pairs represent the URI within the chart (i.e. the key) and the new content (i.e. the value).
	FileChangesGeneratedChart map[string]string `json:"fileChangesGeneratedChart,omitempty"`
	// FileChangesUmbrellaChart provides new content for the umbrella chart.
	// The key value pairs represent the URI within the chart (i.e. the key) and the new content (i.e. the value).
	FileChangesUmbrellaChart map[string]string `json:"fileChangesUmbrellaChart,omitempty"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

// Canary describes the new configuration in a canary release
type Canary struct {
	// Value represents the traffic percentage on the canary
	Value int32 `json:"value,omitempty"`
	// Action represents the action of the canary
	Action CanaryAction `json:"action"`
}

// DeploymentFinishedEventData represents the data for a deployment finished event
type DeploymentFinishedEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Service is the name of the new service
	Service string `json:"service"`
	// TestStrategy is the testing strategy
	TestStrategy string `json:"teststrategy"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy string `json:"deploymentstrategy"`
	// Tag of the new deployed artifact
	Tag string `json:"tag"`
	// Image of the new deployed artifact
	Image string `json:"image"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

// TestsFinishedEventData represents the data for a test finished event
type TestsFinishedEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// TestStrategy is the testing strategy
	TestStrategy string `json:"teststrategy"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy string `json:"deploymentstrategy"`
	// Start indicates the starting timestamp of the tests
	Start string `json:"start"`
	// End indicates the end timestamp of the tests
	End string `json:"end"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
	// Result shows the status of the test
	Result string `json:"result"`
}

// StartEvaluationEventData represents the data for a test finished event
type StartEvaluationEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// TestStrategy is the testing strategy
	TestStrategy string `json:"teststrategy"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy string `json:"deploymentstrategy"`
	// Start indicates the starting timestamp of the tests
	Start string `json:"start"`
	// End indicates the end timestamp of the tests
	End string `json:"end"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

// EvaluationDoneEventData contains information about evaluation results
type EvaluationDoneEventData struct {
	EvaluationDetails *EvaluationDetails `json:"evaluationdetails"`
	// Result is the result of an evaluation; possible values are: pass, warning, fail
	Result string `json:"result"`
	// Project is the name of the project
	Project string `json:"project"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Service is the name of the new service
	Service string `json:"service"`
	// TestStrategy is the testing strategy
	TestStrategy string `json:"teststrategy"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy string `json:"deploymentstrategy"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

type EvaluationDetails struct {
	TimeStart        string                 `json:"timeStart"`
	TimeEnd          string                 `json:"timeEnd"`
	Result           string                 `json:"result"`
	Score            float64                `json:"score"`
	SLOFileContent   string                 `json:"sloFileContent"`
	IndicatorResults []*SLIEvaluationResult `json:"indicatorResults"`
}

type SLIFilter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SLIResult struct {
	Metric  string  `json:"metric"`
	Value   float64 `json:"value"`
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
}

type SLIEvaluationResult struct {
	Score   float64      `json:"score"`
	Value   *SLIResult   `json:"value"`
	Targets []*SLITarget `json:"targets"`
	Status  string       `json:"status"` // pass | warning | fail
}

type SLITarget struct {
	Criteria    string  `json:"criteria"`
	TargetValue float64 `json:"targetValue"`
	Violated    bool    `json:"violated"`
}

// ProblemEventData represents the data for describing a problem
type ProblemEventData struct {
	// State is the state of the problem; possible values are: OPEN, RESOLVED
	State string `json:"State,omitempty"`
	// ProblemID is a unique system identifier of the reported problem
	ProblemID string `json:"ProblemID"`
	// ProblemTitle is the display number of the reported problem.
	ProblemTitle string `json:"ProblemTitle"`
	// ProblemDetails are all problem event details including root cause
	ProblemDetails json.RawMessage `json:"ProblemDetails"`
	// PID is a unique system identifier of the reported problem.
	PID string `json:"PID"`
	// ImpcatedEntity is an identifier of the impacted entity
	ImpactedEntity string `json:"ImpactedEntity,omitempty"`
	// Tags is a comma separated list of tags that are defined for all impacted entities.
	Tags string `json:"Tags,omitempty"`
	// Project is the name of the project
	Project string `json:"project,omitempty"`
	// Stage is the name of the stage
	Stage string `json:"stage,omitempty"`
	// Service is the name of the new service
	Service string `json:"service,omitempty"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

// ConfigureMonitoringEventData represents the data necessary to configure monitoring for a service
type ConfigureMonitoringEventData struct {
	Type string `json:"type"`
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	//ServiceIndicators *models.ServiceIndicators `json:"serviceIndicators"`
	//ServiceObjectives *models.ServiceObjectives `json:"serviceObjectives"`
	//Remediation       *models.Remediations      `json:"remediation"`
}

// InternalGetSLIEventData describes a set of SLIs to be retrieved by a data source
type InternalGetSLIEventData struct {
	// SLIProvider is the name of the SLI provider which is queried
	SLIProvider string `json:"sliProvider"`
	// Project is the name of the project
	Project string `json:"project"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Service is the name of the new service
	Service string `json:"service"`
	Start   string `json:"start"`
	End     string `json:"end"`
	// TestStrategy is the testing strategy
	TestStrategy string `json:"teststrategy"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy string       `json:"deploymentstrategy"`
	Deployment         string       `json:"deployment"`
	Indicators         []string     `json:"indicators"`
	CustomFilters      []*SLIFilter `json:"customFilters"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

// InternalGetSLIDoneEventData contains a list of SLIs and their values
type InternalGetSLIDoneEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Service is the name of the new service
	Service string `json:"service"`
	Start   string `json:"start"`
	End     string `json:"end"`
	// TestStrategy is the testing strategy
	TestStrategy    string       `json:"teststrategy"`
	IndicatorValues []*SLIResult `json:"indicatorValues"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy string `json:"deploymentstrategy"`
	Deployment         string `json:"deployment"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}
