package models

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
