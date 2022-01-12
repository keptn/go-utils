package models

import "encoding/json"

// Stage stage
type Stage struct {

	// services
	Services []*Service `json:"services"`

	// Stage name
	StageName string `json:"stageName,omitempty"`
}

func (s *Stage) ToJSON() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *Stage) FromJSON(b []byte) error {
	var res Stage
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*s = res
	return nil
}
