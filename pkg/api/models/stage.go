package models

import "encoding/json"

// Stage stage
type Stage struct {

	// services
	Services []*Service `json:"services"`

	// Stage name
	StageName string `json:"stageName,omitempty"`
}

// ToJSON converts object to JSON string
func (s *Stage) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON converts JSON string to object
func (s *Stage) FromJSON(b []byte) error {
	var res Stage
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*s = res
	return nil
}
