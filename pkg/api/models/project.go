package models

import "encoding/json"

// Project project
type Project struct {

	// Creation date of the service
	CreationDate string `json:"creationDate,omitempty"`

	// Git remote URI
	GitRemoteURI string `json:"gitRemoteURI,omitempty"`

	// Git token
	GitToken string `json:"gitToken,omitempty"`

	// Git User
	GitUser string `json:"gitUser,omitempty"`

	// Project name
	ProjectName string `json:"projectName,omitempty"`

	// Shipyard version
	ShipyardVersion string `json:"shipyardVersion,omitempty"`

	// stages
	Stages []*Stage `json:"stages"`
}

func (p *Project) ToJSON() ([]byte, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

func (p *Project) FromJSON(b []byte) error {
	var res Project
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
}
