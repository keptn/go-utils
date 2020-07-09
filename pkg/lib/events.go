package keptn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
)

import (
	"encoding/json"
)

const MAX_SEND_RETRIES = 3

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

// ApprovalTriggeredEvent is a CloudEvent triggering a new approval
const ApprovalTriggeredEventType = "sh.keptn.event.approval.triggered"

// ApprovalFinishedEventType is a CloudEvent confirming an open approval
const ApprovalFinishedEventType = "sh.keptn.event.approval.finished"

// ActionTriggeredEventType is a CloudEvent for triggering actions
const ActionTriggeredEventType = "sh.keptn.event.action.triggered"

// ActionStartedEventType is a CloudEvent for confirming the start of an action
const ActionStartedEventType = "sh.keptn.event.action.started"

// ActionFinishedEventType is a CloudEvent for indicating that an action has been finished
const ActionFinishedEventType = "sh.keptn.event.action.finished"

// RemediationTriggeredEventType is a CloudEvent for triggering remediations
const RemediationTriggeredEventType = "sh.keptn.event.remediation.triggered"

// RemediationTriggeredEventType is a CloudEvent to indicate the start of a remediation
const RemediationStartedEventType = "sh.keptn.event.remediation.started"

// RemediationStatusChangedEventType is a CloudEvent to indicate that the status of a remediation has been changed
const RemediationStatusChangedEventType = "sh.keptn.event.remediation.status.changed"

// RemediationFinishedEventType is a CloudEvent to indicate that the status of a remediation has been changed
const RemediationFinishedEventType = "sh.keptn.event.remediation.finished"

const TestTriggeredEventType = "sh.keptn.event.test.triggered"

// KeptnBase contains properties that are shared among most Keptn events
type KeptnBase struct {
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage        string  `json:"stage"`
	TestStrategy *string `json:"teststrategy,omitempty"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy *string `json:"deploymentstrategy,omitempty"`
	// Tag of the new deployed artifact
	Tag *string `json:"tag,omitempty"`
	// Image of the new deployed artifact
	Image *string `json:"image,omitempty"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

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
	// DeploymentURILocal contains the local URL
	DeploymentURILocal string `json:"deploymentURILocal,omitempty"`
	// DeploymentURIPublic contains the public URL
	DeploymentURIPublic string `json:"deploymentURIPublic,omitempty"`
}

// TestTriggeredEventData represents the data for a new test
type TestTriggeredEventData struct {
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
	// DeploymentURILocal contains the local URL
	DeploymentURILocal string `json:"deploymentURILocal,omitempty"`
	// DeploymentURIPublic contains the public URL
	DeploymentURIPublic string `json:"deploymentURIPublic,omitempty"`
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
	// ProblemURL is a back link to the original problem
	ProblemURL string `json:"ProblemURL,omitempty"`
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

type ApprovalData struct {
	Result string `json:"result"`
	Status string `json:"status"`
}

// ApprovalTriggeredEventData contains information about an approval.triggered event
type ApprovalTriggeredEventData struct {
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage        string  `json:"stage"`
	TestStrategy *string `json:"teststrategy,omitempty"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy *string `json:"deploymentstrategy,omitempty"`
	// Tag of the new deployed artifact
	Tag string `json:"tag,omitempty"`
	// Image of the new deployed artifact
	Image string `json:"image,omitempty"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`

	// Result is the result of an evaluation; possible values are: pass, warning, fail
	Result string `json:"result"`
}

// ApprovalTriggeredEventData contains information about an approval.finished event
type ApprovalFinishedEventData struct {
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage        string  `json:"stage"`
	TestStrategy *string `json:"teststrategy,omitempty"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy *string `json:"deploymentstrategy,omitempty"`
	// Tag of the new deployed artifact
	Tag string `json:"tag,omitempty"`
	// Image of the new deployed artifact
	Image string `json:"image,omitempty"`
	// Labels contains labels
	Labels   map[string]string `json:"labels"`
	Approval ApprovalData      `json:"approval"`
}

// ProblemDetails contains information about a problem
type ProblemDetails struct {
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
	// ProblemURL is a back link to the original problem
	ProblemURL     string `json:"ProblemURL,omitempty"`
	ImpactedEntity string `json:"ImpactedEntity,omitempty"`
	// Tags is a comma separated list of tags that are defined for all impacted entities.
	Tags string `json:"Tags,omitempty"`
}

// ActionInfo contains information about the action to be performed
type ActionInfo struct {
	// Name is the name of the action
	Name string `json:"name"`
	// Action determines the type of action to be executed
	Action string `json:"action"`
	// Description contains the description of the action
	Description string `json:"description"`
	// Value contains the value of the action
	Value interface{} `json:"value,omitempty"`
}

// ActionTriggeredEventData contains information about an action.triggered event
type ActionTriggeredEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Action describes the type of action
	Action ActionInfo `json:"action"`
	// Problem contains details about the problem
	Problem ProblemDetails `json:"problem"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

// ActionStartedEventData contains information about an action.started event
type ActionStartedEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

type ActionResultType string

const (
	ActionResultPass   ActionResultType = "pass"
	ActionResultFailed ActionResultType = "failed"
)

type ActionStatusType string

const (
	ActionStatusSucceeded ActionStatusType = "succeeded"
	ActionStatusErrored   ActionStatusType = "errored"
	ActionStatusUnknown   ActionStatusType = "unknown"
)

// ActionResult contains information about the execution of an action
type ActionResult struct {
	Result ActionResultType `json:"result"`
	Status ActionStatusType `json:"status"`
}

// ActionFinishedEventData contains information about the execution of an action
type ActionFinishedEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Action describes the type of action
	Action ActionResult `json:"action"`
	// Value contains additional values
	Value interface{} `json:"value,omitempty"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
}

// RemediationTriggeredEventData is a CloudEvent for triggering remediations
type RemediationTriggeredEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Problem contains details about the problem
	Problem ProblemDetails `json:"problem"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
	// Remediation describes the remediation
	Remediation interface{} `json:"remediation"`
}

// RemediationTriggeredEventType is a CloudEvent to indicate the start of a remediation
type RemediationStartedEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
	// Remediation describes the remediation
	Remediation interface{} `json:"remediation"`
}

type RemediationStatusType string

const (
	RemediationStatusSucceeded RemediationStatusType = "succeeded"
	RemediationStatusErrored   RemediationStatusType = "errored"
	RemediationStatusUnknown   RemediationStatusType = "unknown"
)

// RemediationStatus describes the status of a remediation
type RemediationResult struct {
	// ActionIndex is the index of the action
	ActionIndex int `json:"actionIndex"`
	// ActionName is the name of the action
	ActionName string `json:"actionName"`
}

// RemediationStatus describes the result and status of a remediation
type RemediationStatusChanged struct {
	// Status describes the status of the remediation
	Status RemediationStatusType `json:"status"`
	// RemediationResult indicates the result
	Result RemediationResult `json:"result"`
}

// RemediationStatusChangedEventData is a CloudEvent to indicate that the status of a remediation has been changed
type RemediationStatusChangedEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
	// Remediation describes the remediation
	Remediation RemediationStatusChanged `json:"remediation"`
}

type RemediationResultType string

const (
	RemediationResultPass   RemediationResultType = "pass"
	RemediationResultFailed RemediationResultType = "failed"
)

// RemediationFinished describes a finished remediation
type RemediationFinished struct {
	// Status describes the status
	Status RemediationStatusType `json:"status"`
	// Result describes the result
	Result RemediationResultType `json:"result"`
	// Message contains a message about the status
	Message string `json:"message"`
}

// RemediationFinishedEventData is a CloudEvent to indicate that the status of a remediation has been changed
type RemediationFinishedEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Problem contains details about the problem
	Problem ProblemDetails `json:"problem"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
	// Remediation describes the remediation
	Remediation RemediationFinished `json:"remediation"`
}

//
// Sends a ConfigurationChangeEventType = "sh.keptn.event.configuration.change"
//
func (k *Keptn) SendConfigurationChangeEvent(incomingEvent *cloudevents.Event, labels map[string]string, eventSource string) error {
	source, _ := url.Parse(eventSource)
	contentType := "application/json"

	configurationChangeData := ConfigurationChangeEventData{}

	// if we have an incoming event we pre-populate data
	if incomingEvent != nil {
		incomingEvent.DataAs(&configurationChangeData)
	}

	if k.KeptnBase.Project != "" {
		configurationChangeData.Project = k.KeptnBase.Project
	}
	if k.KeptnBase.Service != "" {
		configurationChangeData.Service = k.KeptnBase.Service
	}
	if k.KeptnBase.Stage != "" {
		configurationChangeData.Stage = k.KeptnBase.Stage
	}
	if labels != nil {
		configurationChangeData.Labels = labels
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        ConfigurationChangeEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": k.KeptnContext},
		}.AsV02(),
		Data: configurationChangeData,
	}

	log.Println(fmt.Sprintf("%s", event))

	return k.SendCloudEvent(event)
}

//
// Sends a DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
//
func (k *Keptn) SendDeploymentFinishedEvent(incomingEvent *cloudevents.Event, teststrategy, deploymentstrategy, image, tag, deploymentURILocal, deploymentURIPublic string, labels map[string]string, eventSource string) error {
	source, _ := url.Parse(eventSource)
	contentType := "application/json"

	deploymentFinishedData := DeploymentFinishedEventData{}

	// if we have an incoming event we pre-populate data
	if incomingEvent != nil {
		incomingEvent.DataAs(&deploymentFinishedData)
	}

	if k.KeptnBase.Project != "" {
		deploymentFinishedData.Project = k.KeptnBase.Project
	}
	if k.KeptnBase.Service != "" {
		deploymentFinishedData.Service = k.KeptnBase.Service
	}
	if k.KeptnBase.Stage != "" {
		deploymentFinishedData.Stage = k.KeptnBase.Stage
	}
	if teststrategy != "" {
		deploymentFinishedData.TestStrategy = teststrategy
	}
	if deploymentstrategy != "" {
		deploymentFinishedData.DeploymentStrategy = deploymentstrategy
	}
	if image != "" {
		deploymentFinishedData.Image = image
	}
	if tag != "" {
		deploymentFinishedData.Tag = tag
	}

	if labels != nil {
		deploymentFinishedData.Labels = labels
	}

	if deploymentURILocal != "" {
		deploymentFinishedData.DeploymentURILocal = deploymentURILocal
	}

	if deploymentURIPublic != "" {
		deploymentFinishedData.DeploymentURIPublic = deploymentURIPublic
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        DeploymentFinishedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": k.KeptnContext},
		}.AsV02(),
		Data: deploymentFinishedData,
	}

	log.Println(fmt.Sprintf("%s", event))

	return k.SendCloudEvent(event)

}

//
// Sends a TestsFinishedEventType = "sh.keptn.events.tests-finished"
//
func (k *Keptn) SendTestsFinishedEvent(incomingEvent *cloudevents.Event, teststrategy, deploymentstrategy string, startedAt time.Time, result string, labels map[string]string, eventSource string) error {
	source, _ := url.Parse(eventSource)
	contentType := "application/json"

	testFinishedData := TestsFinishedEventData{}

	// if we have an incoming event we pre-populate data
	if incomingEvent != nil {
		incomingEvent.DataAs(&testFinishedData)
	}

	if k.KeptnBase.Project != "" {
		testFinishedData.Project = k.KeptnBase.Project
	}
	if k.KeptnBase.Service != "" {
		testFinishedData.Service = k.KeptnBase.Service
	}
	if k.KeptnBase.Stage != "" {
		testFinishedData.Stage = k.KeptnBase.Stage
	}
	if teststrategy != "" {
		testFinishedData.TestStrategy = teststrategy
	}
	if deploymentstrategy != "" {
		testFinishedData.DeploymentStrategy = deploymentstrategy
	}

	if labels != nil {
		testFinishedData.Labels = labels
	}

	// fill in timestamps
	testFinishedData.Start = startedAt.Format(time.RFC3339)
	testFinishedData.End = time.Now().Format(time.RFC3339)

	// set test result
	testFinishedData.Result = result

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        TestsFinishedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": k.KeptnContext},
		}.AsV02(),
		Data: testFinishedData,
	}

	log.Println(fmt.Printf("%s", event))

	return k.SendCloudEvent(event)
}

//
// Sends a CloudEvent to the event broker
//
func (k *Keptn) SendCloudEvent(event cloudevents.Event) error {
	if k.useLocalFileSystem {
		log.Println(fmt.Printf("%v", event.Data))
		return nil
	}
	transport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(k.eventBrokerURL),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
		cloudeventshttp.WithHeader("x-token", os.Getenv("API_TOKEN")),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(transport)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	for i := 0; i <= MAX_SEND_RETRIES; i++ {
		_, _, err = c.Send(context.Background(), event)
		if err == nil {
			return nil
		}
		<-time.After(getExpBackoffTime(i + 1))
	}
	return errors.New("Failed to send cloudevent:, " + err.Error())
}

func getExpBackoffTime(retryNr int) time.Duration {
	f := 1.5 * float64(retryNr)
	if retryNr <= 1 {
		f = 1.5
	}
	currentInterval := float64(500*time.Millisecond) * f
	randomizationFactor := 0.5
	random := rand.Float64()

	var delta = randomizationFactor * currentInterval
	minInterval := float64(currentInterval) - delta
	maxInterval := float64(currentInterval) + delta

	return time.Duration(minInterval + (random * (maxInterval - minInterval + 1)))
}
