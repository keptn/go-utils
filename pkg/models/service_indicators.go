package models

// ServiceIndicators contains the definition of service indicators
type ServiceIndicators struct {
	Indicators []*ServiceIndicator `json:"indicators" yaml:"indicators"`
}

// ServiceIndicator describes a service indicator
type ServiceIndicator struct {
	Name   string `json:"name" yaml:"name"`
	Source string `json:"source" yaml:"source"`
	Query  string `json:"query" yaml:"query"`
}
