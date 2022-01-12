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

func (v *Version) ToJSON() ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

func (v *Version) FromJSON(b []byte) error {
	var res Version
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*v = res
	return nil
}
