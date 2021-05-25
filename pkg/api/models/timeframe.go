package models

// Timeframe timeframe
type Timeframe struct {

	// Evaluation start timestamp
	From string `json:"from,omitempty"`

	// Evaluation timeframe
	Timeframe string `json:"timeframe,omitempty"`

	// Evaluation end timestamp
	To string `json:"to,omitempty"`
}
