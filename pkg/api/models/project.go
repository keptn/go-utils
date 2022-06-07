package models

import "encoding/json"

// Project project
type Project struct {

	// Creation date of the service
	CreationDate string `json:"creationDate,omitempty"`

	// Project name
	ProjectName string `json:"projectName,omitempty"`

	// Shipyard version
	ShipyardVersion string `json:"shipyardVersion,omitempty"`

	// stages
	Stages []*Stage `json:"stages"`

	// git auth credentials
	GitCredentials GitAuthCredentials `json:"gitCredentials"`
}

// ToJSON converts object to JSON string
func (p *Project) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *Project) FromJSON(b []byte) error {
	var res Project
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
}
