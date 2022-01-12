package models

import "encoding/json"

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

func (c *CreateProject) ToJSON() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *CreateProject) FromJSON(b []byte) error {
	var res CreateProject
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*c = res
	return nil
}
