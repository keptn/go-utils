package models

// Secret secret
type Secret struct {

	// data
	// Required: true
	Data map[string]string `json:"data"`

	// The name of the secret
	// Required: true
	Name *string `json:"name"`

	// The scope of the secret
	// Required: true
	Scope *string `json:"scope"`
}
