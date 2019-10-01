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
	QueryObject []*ServiceIndicatorQueryObject `json:"queryObject" yaml:"queryObject"`
}

type ServiceIndicatorQueryObject struct {
	Key string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}
