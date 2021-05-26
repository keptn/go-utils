package models

// Approval approval
type Approval struct {

	// ID of the event
	EventID string `json:"eventId,omitempty"`

	// image
	Image string `json:"image,omitempty"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty"`

	// tag
	Tag string `json:"tag,omitempty"`

	// Time of the event
	Time string `json:"time,omitempty"`
}
