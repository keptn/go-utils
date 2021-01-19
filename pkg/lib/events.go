package keptn

import "encoding/json"

// ProblemEventType is a CloudEvent type to inform about a problem
const ProblemEventType = "sh.keptn.events.problem"

// ProblemOpenEventType is a CloudEvent type to inform about an open problem
const ProblemOpenEventType = "sh.keptn.event.problem.open"

// ConfigureMonitoringEventType is a CloudEvent for configuring monitoring
const ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"

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
}
