package models

import "encoding/json"

// CreateProject create project
type CreateProject struct {
	// name
	// Required: true
	Name *string `json:"name"`

	// shipyard
	// Required: true
	Shipyard *string `json:"shipyard"`

	// git auth credentials
	GitCredentials GitAuthCredentials `json:"gitCredentials"`
}

// ToJSON converts object to JSON string
func (c *CreateProject) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// FromJSON converts JSON string to object
func (c *CreateProject) FromJSON(b []byte) error {
	var res CreateProject
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*c = res
	return nil
}
