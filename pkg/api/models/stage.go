package models

// Stage stage
type Stage struct {

	// services
	Services []*Service `json:"services"`

	// Stage name
	StageName string `json:"stageName,omitempty"`
}
