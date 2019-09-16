package models

// ServiceIndicators contains the definition of service indicators
type ServiceIndicators struct {
	Indicators []*ServiceIndicator `json:"indicators" yaml:"indicators"`
}

// ServiceIndicator describes a service indicator
type ServiceIndicator struct {
	Metric string `json:"metric" yaml:"metric"`
	Source string `json:"source" yaml:"source"`
	Query  string `json:"query" yaml:"query"`
}
