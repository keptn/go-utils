package models

// ExpandedProject expanded project
//
// swagger:model ExpandedProject
type ExpandedProject struct {

	// Creation date of the project
	CreationDate string `json:"creationDate,omitempty"`

	// last event context
	LastEventContext *EventContextInfo `json:"lastEventContext,omitempty"`

	// Project name
	ProjectName string `json:"projectName,omitempty"`

	// Shipyard file content
	Shipyard string `json:"shipyard,omitempty"`

	// Version of the shipyard file
	ShipyardVersion string `json:"shipyardVersion,omitempty"`

	// stages
	Stages []*ExpandedStage `json:"stages"`

	// git auth credentials
	GitCredentials *GitAuthCredentialsSecure `json:"gitCredentials,omitempty"`
}
