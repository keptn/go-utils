package models

import "encoding/json"

// Metadata metadata
type Metadata struct {

	// automaticprovisioning
	Automaticprovisioning bool `json:"automaticprovisioning,omitempty"`

	// bridgeversion
	Bridgeversion string `json:"bridgeversion,omitempty"`

	// keptnlabel
	Keptnlabel string `json:"keptnlabel,omitempty"`

	// keptnservices
	Keptnservices interface{} `json:"keptnservices,omitempty"`

	// keptnversion
	Keptnversion string `json:"keptnversion,omitempty"`

	// namespace
	Namespace string `json:"namespace,omitempty"`

	// shipyardversion
	Shipyardversion string `json:"shipyardversion,omitempty"`
}

// ToJSON converts object to JSON string
func (m *Metadata) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON converts JSON string to object
func (m *Metadata) FromJSON(b []byte) error {
	var res Metadata
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
