package models

// Secret secret
type Secret struct {

	// data
	// Required: true
	Data map[string]string `json:"data"`

	SecretMetadata
}

type SecretMetadata struct {
	// The name of the secret
	// Required: true
	Name *string `json:"name" yaml:"name"`

	// The scope of the secret
	// Required: true
	Scope *string `json:"scope,omitempty" yaml:"scope,omitempty"`
}

type GetSecretsResponse struct {
	Secrets []SecretMetadata `json:"secrets"`
}
