package models

// CreateProject create project
type CreateProject struct {

	// git remote URL
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`

	// git token
	GitToken string `json:"gitToken,omitempty"`

	// git user
	GitUser string `json:"gitUser,omitempty"`

	// name
	// Required: true
	Name *string `json:"name"`

	// shipyard
	// Required: true
	Shipyard *string `json:"shipyard"`
}
