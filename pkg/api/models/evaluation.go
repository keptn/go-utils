package models

import "encoding/json"

type Evaluation struct {

	// Evaluation start timestamp
	Start string `json:"start,omitempty"`

	// labels
	Labels map[string]string `json:"labels,omitempty"`

	// Evaluation timeframe
	Timeframe string `json:"timeframe,omitempty"`

	// Evaluation end timestamp
	End string `json:"end,omitempty"`
}

func (e *Evaluation) ToJSON() ([]byte, error) {
	if e == nil {
		return nil, nil
	}
	return json.Marshal(e)
}

func (e *Evaluation) FromJSON(b []byte) error {
	var res Evaluation
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*e = res
	return nil
}
