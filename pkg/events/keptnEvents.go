package events

// ServiceCreateEventType is a CloudEvent type for creating new services
const ServiceCreateEventType = "sh.keptn.event.service.create"

// InternalServiceCreateEventType is a CloudEvent type for creating new services
const InternalServiceCreateEventType = "sh.keptn.internal.event.service.create"

// ConfigurationChangeEventType is a CloudEvent type for changing the configuration
const ConfigurationChangeEventType = "sh.keptn.event.configuration.change"

// ServiceCreateEventData represents the data for creating a new service
type ServiceCreateEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// HelmChart are the data of a Helm chart packed as tgz
	HelmChart []byte `json:"helmChart"`
}
