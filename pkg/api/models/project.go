package models

import "encoding/json"

// Project project
type Project struct {

	// Creation date of the service
	CreationDate string `json:"creationDate,omitempty"`

	// Git remote URI
	GitRemoteURI string `json:"gitRemoteURI,omitempty"`

	// Git token
	GitToken string `json:"gitToken,omitempty"`

	// git private key
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`

	// git private key passphrase
	GitPrivateKeyPass string `json:"gitPrivateKeyPass,omitempty"`

	// git proxy URL
	GitProxyURL string `json:"gitProxyUrl,omitempty"`

	// git proxy scheme
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// git proxy user
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// git proxy secure
	GitProxyInsecure bool `json:"gitProxyInsecure,omitempty"`

	// git proxy password
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`

	// Git User
	GitUser string `json:"gitUser,omitempty"`

	// Project name
	ProjectName string `json:"projectName,omitempty"`

	// Shipyard version
	ShipyardVersion string `json:"shipyardVersion,omitempty"`

	// stages
	Stages []*Stage `json:"stages"`
}

// ToJSON converts object to JSON string
func (p *Project) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *Project) FromJSON(b []byte) error {
	var res Project
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
}
