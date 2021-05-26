package models

// Service service
type Service struct {

	// Creation date of the service
	CreationDate string `json:"creationDate,omitempty"`

	// Currently deployed image
	DeployedImage string `json:"deployedImage,omitempty"`

	// last event types
	LastEventTypes map[string]EventContextInfo `json:"lastEventTypes,omitempty"`

	// open approvals
	OpenApprovals []*Approval `json:"openApprovals"`

	// Service name
	ServiceName string `json:"serviceName,omitempty"`
}
