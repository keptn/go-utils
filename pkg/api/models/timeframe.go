package models

import "encoding/json"

// Timeframe timeframe
type Timeframe struct {

	// Evaluation start timestamp
	From string `json:"from,omitempty"`

	// Evaluation timeframe
	Timeframe string `json:"timeframe,omitempty"`

	// Evaluation end timestamp
	To string `json:"to,omitempty"`
}

// ToJSON converts object to JSON string
func (t *Timeframe) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

// FromJSON converts JSON string to object
func (t *Timeframe) FromJSON(b []byte) error {
	var res Timeframe
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*t = res
	return nil
}
