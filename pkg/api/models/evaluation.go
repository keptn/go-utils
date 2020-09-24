package models

type Evaluation struct {

	// Evaluation start timestamp
	From string `json:"from,omitempty"`

	// labels
	Labels map[string]string `json:"labels,omitempty"`

	// Evaluation timeframe
	Timeframe string `json:"timeframe,omitempty"`

	// Evaluation end timestamp
	To string `json:"to,omitempty"`
}
