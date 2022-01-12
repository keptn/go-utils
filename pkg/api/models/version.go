package models

import "encoding/json"

// Version version
type Version struct {

	// Branch in repository containing the resource
	Branch string `json:"branch,omitempty"`

	// Upstream respository containing the resource
	UpstreamURL string `json:"upstreamURL,omitempty"`

	// Version identifier
	Version string `json:"version,omitempty"`
}

// ToJSON converts object to JSON string
func (v *Version) ToJSON() ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

// FromJSON converts JSON string to object
func (v *Version) FromJSON(b []byte) error {
	var res Version
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*v = res
	return nil
}
