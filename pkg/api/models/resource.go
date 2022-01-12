package models

import "encoding/json"

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

// ToJSON converts object to JSON string
func (r *Resource) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

// FromJSON converts JSON string to object
func (r *Resource) FromJSON(b []byte) error {
	var res Resource
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*r = res
	return nil
}
