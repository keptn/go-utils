package models

import "encoding/json"

// Metadata metadata
type Metadata struct {

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
}

func (m *Metadata) ToJSON() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Metadata) FromJSON(b []byte) error {
	var res Metadata
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
