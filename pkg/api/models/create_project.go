package models

import "encoding/json"

// CreateProject create project
type CreateProject struct {

	// git remote URL
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`

	// git token
	GitToken string `json:"gitToken,omitempty"`

	// git private key
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`

	// git user
	GitUser string `json:"gitUser,omitempty"`

	// git proxy
	GitProxyUrl string `json:"gitProxyUrl,omitempty"`

	// git proxy
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// git proxy
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// git proxy
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`

	// name
	// Required: true
	Name *string `json:"name"`

	// shipyard
	// Required: true
	Shipyard *string `json:"shipyard"`
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
