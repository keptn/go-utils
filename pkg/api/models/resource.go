package models

// Resource resource
type Resource struct {

	// Metadata
	Metadata *Version `json:"metadata,omitempty"`

	// Resource content
	ResourceContent string `json:"resourceContent,omitempty"`

	// Resource URI
	// Required: true
	ResourceURI *string `json:"resourceURI"`
}
