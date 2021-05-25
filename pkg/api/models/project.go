package models

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
