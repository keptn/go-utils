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

func (r *Resource) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

func (r *Resource) FromJSON(b []byte) error {
	var res Resource
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*r = res
	return nil
}
