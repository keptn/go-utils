package events

import (
	"github.com/keptn/go-utils/pkg/models"
)

// ServiceCreateEventType is a CloudEvent type for creating a new service
const ServiceCreateEventType = "sh.keptn.event.service.create"

// InternalServiceCreateEventType is a CloudEvent type for creating a new service
const InternalServiceCreateEventType = "sh.keptn.internal.event.service.create"

// ProjectCreateEventType is a CloudEvent type for creating a new project
const ProjectCreateEventType = "sh.keptn.event.project.create"

// ProjectDeleteEventType is a CloudEvent type for deleting a project
const ProjectDeleteEventType = "sh.keptn.event.project.delete"

// InternalProjectCreateEventType is a CloudEvent type for creating a new project
const InternalProjectCreateEventType = "sh.keptn.internal.event.project.create"

// InternalProjectDeleteEventType is a CloudEvent type for deleting a project
const InternalProjectDeleteEventType = "sh.keptn.internal.event.project.delete"

// ConfigurationChangeEventType is a CloudEvent type for changing the configuration
const ConfigurationChangeEventType = "sh.keptn.event.configuration.change"

// ProblemOpenEventType is a CloudEvent type to inform about an open problem
const ProblemOpenEventType = "sh.keptn.event.problem.open"

// ConfigureMonitoringEventType is a CloudEvent for configuring monitoring
const ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"

// TestsFinishedEventType is a CloudEvent for indicating that tests have finished
const TestsFinishedEventType = "sh.keptn.event.tests.finished"

// TestFinishedEventType_0_5_0_Compatible is a CloudEvent for indicating that tests have finished
const TestFinishedEventType_0_5_0_Compatible = "sh.keptn.events.tests-finished"

// EvaluationStartEventType is a CloudEvent to trigger an evaluation step
const EvaluationStartEventType = "sh.keptn.event.evaluation.start"

// EvaluationDoneEventType is a CloudEvent for indicating that the evaluation has finished
const EvaluationDoneEventType = "sh.keptn.events.evaluation-done"

// DeploymentFinishedEventType is a CloudEvent for indicating that the deployment has finished
const DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"

// StartEvaluationEventType is a CloudEvent for retrieving SLI values
const StartEvaluationEventType = "sh.keptn.event.start-evaluation"

// InternalGetSLIEventType is a CloudEvent for retrieving SLI values
const InternalGetSLIEventType = "sh.keptn.internal.event.get-sli"

// InternalGetSLIEventType is a CloudEvent for submitting SLI values
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

// ProjectDeleteEventData represents the data for deleting a new project
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
}

// EvaluationStartEventData represents the data for a evaluation start event
type EvaluationStartEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
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
}

// Canary describes the new configuration in a canary release
type Canary struct {
	// Value represents the traffic percentage on the canary
	Value int32 `json:"value,omitempty"`
	// Action represents the action of the canary
	Action CanaryAction `json:"action"`
}

// ProblemEventData represents the data for describing a problem
type ProblemEventData struct {
	State          string `json:"state"`
	ProblemID      string `json:"problemID"`
	ProblemTitle   string `json:"problemtitle"`
	ProblemDetails string `json:"problemdetails"`
	ImpactedEntity string `json:"impactedEntity"`
	Tags           string `json:"tags,omitempty"`
	Project        string `json:"project,omitempty"`
	Stage          string `json:"stage,omitempty"`
	Service        string `json:"service,omitempty"`
}

// ConfigureMonitoringEventData represents the data necessary to configure monitoring for a service
type ConfigureMonitoringEventData struct {
	Type              string                    `json:"type"`
	Project           string                    `json:"project"`
	Service           string                    `json:"service"`
	ServiceIndicators *models.ServiceIndicators `json:"serviceIndicators"`
	ServiceObjectives *models.ServiceObjectives `json:"serviceObjectives"`
	Remediation       *models.Remediations      `json:"remediation"`
}

// InternalGetSLIEventData describes a set of SLIs to be retrieved by a data source
type InternalGetSLIEventData struct {
	SLIProvider   string       `json:"sliProvider"`
	Project       string       `json:"project"`
	Service       string       `json:"service"`
	Stage         string       `json:"stage"`
	Start         string       `json:"start"`
	End           string       `json:"end"`
	Indicators    []string     `json:"indicators"`
	CustomFilters []*SLIFilter `json:"customFilters"`
}

// InternalGetSLIDoneEventData contains a list of SLIs and their values
type InternalGetSLIDoneEventData struct {
	Project         string       `json:"project"`
	Service         string       `json:"service"`
	Stage           string       `json:"stage"`
	IndicatorValues []*SLIResult `json:"indicatorValues"`
}

type SLIFilter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SLIResult struct {
	Metric string  `json:"metric"`
	Value  float64 `json:"value"`
}
