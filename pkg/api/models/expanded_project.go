package models

import "encoding/json"

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

	// is upstream auto provisioned
	IsUpstreamAutoProvisioned bool `json:"isUpstreamAutoProvisioned"`
}

// ToJSON converts object to JSON string
func (a *ExpandedProject) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

// FromJSON converts JSON string to object
func (a *ExpandedProject) FromJSON(b []byte) error {
	var res ExpandedProject
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*a = res
	return nil
}
