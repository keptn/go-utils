package events

// ServiceCreateEventType is a CloudEvents type for creating new services
const ServiceCreateEventType = "sh.keptn.event.service.create"

// ServiceCreateEventData represents the data for creating a new service
type ServiceCreateEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// HelmChart is a base64 encoded Helm chart packed as tgz
	HelmChart []byte `json:"helmChart"`
}
