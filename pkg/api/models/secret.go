package models

import "encoding/json"

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

type GetSecretResponseItem struct {
	SecretMetadata `yaml:",inline"`
	Keys           []string `json:"keys" yaml:"keys"`
}

type GetSecretsResponse struct {
	Secrets []GetSecretResponseItem `json:"secrets" yaml:"secrets"`
}

func (s *Secret) ToJSON() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *Secret) FromJSON(b []byte) error {
	var res Secret
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*s = res
	return nil
}

func (s *GetSecretsResponse) ToJSON() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *GetSecretsResponse) FromJSON(b []byte) error {
	var res GetSecretsResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*s = res
	return nil
}
