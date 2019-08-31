package events

import (
	"k8s.io/helm/pkg/proto/hapi/chart"
)

// ServiceCreateEventType is a CloudEvent type for creating new services
const ServiceCreateEventType = "sh.keptn.event.service.create"

// InternalServiceCreateEventType is a CloudEvent type for creating new services
const InternalServiceCreateEventType = "sh.keptn.internal.event.service.create"

// ConfigurationChangeEventType is a CloudEvent type for changing the configuration
const ConfigurationChangeEventType = "sh.keptn.event.configuration.change"

// ProblemOpenEventType is a CloudEvent type to inform about an open problem
const ProblemOpenEventType = "sh.keptn.event.problem.open"

// ServiceCreateEventData represents the data for creating a new service
type ServiceCreateEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// HelmChart are the data of a Helm chart packed as tgz
	HelmChart []byte `json:"helmChart"`
}

// ConfigurationChangeEventData represents the data for
type ConfigurationChangeEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// ValuesPrimary contains new Helm values for primary
	ValuesPrimary map[string]*chart.Value `json:"valuesPrimary,omitempty"`
	// ValuesCanary contains new Helm values for canary
	ValuesCanary map[string]*chart.Value `json:"valuesCanary,omitempty"`
	// Canary contains a new configuration for canary releases
	Canary *Canary `json:"canary,omitempty"`
}

// Canary describes the new configuration in a canary release
type Canary struct {
	// Value represents the traffic percentage on the canary
	Value int32 `json:"value,omitempty"`
	// Action represents the action of the canary
	Action CanaryAction `json:"action,omitempty"`
}

// ProblemEventData represents the data for describing a problem
type ProblemEventData struct {
	State          string `json:"state"`
	ProblemID      string `json:"problemID"`
	ProblemTitle   string `json:"problemtitle"`
	ProblemDetails string `json:"problemdetails"`
	ImpactedEntity string `json:"impactedEntity"`
}
