package events

import "github.com/keptn/go-utils/pkg/models"

// ServiceCreateEventType is a CloudEvent type for creating new services
const ServiceCreateEventType = "sh.keptn.event.service.create"

// InternalServiceCreateEventType is a CloudEvent type for creating new services
const InternalServiceCreateEventType = "sh.keptn.internal.event.service.create"

// ConfigurationChangeEventType is a CloudEvent type for changing the configuration
const ConfigurationChangeEventType = "sh.keptn.event.configuration.change"

// ProblemOpenEventType is a CloudEvent type to inform about an open problem
const ProblemOpenEventType = "sh.keptn.event.problem.open"

// ConfigureMonitoringEventType is a CloudEvent for configuring monitoring
const ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"

// ServiceCreateEventData represents the data for creating a new service
type ServiceCreateEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// HelmChart are the data of a Helm chart packed as tgz
	HelmChart []byte `json:"helmChart"`
}

// ProblemEventData represents the data for describing a problem
type ProblemEventData struct {
	State          string `json:"state"`
	ProblemID      string `json:"problemID"`
	ProblemTitle   string `json:"problemtitle"`
	ProblemDetails string `json:"problemdetails"`
	ImpactedEntity string `json:"impactedEntity"`
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
