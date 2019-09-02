package models

// ServiceObjectives describes objectives for a service
type ServiceObjectives struct {
	Pass       int                 `json:"pass" yaml:"pass"`
	Warning    int                 `json:"warning" yaml:"warning"`
	Objectives []*ServiceObjective `json:"objectives" yaml:"objectives"`
}

// ServiceObjective describes a service objective
type ServiceObjective struct {
	Name      string  `json:"name" yaml:"name"`
	Threshold float32 `json:"threshold" yaml:"threshold"`
	Timeframe string  `json:"timeframe" yaml:"timeframe"`
	Score     float32 `json:"score" yaml:"score"`
}
