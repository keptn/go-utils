package models

// Remediations contains remediation definitions for a service
type Remediations struct {
	Remediations []*Remediation `json:"remediations" yaml:"remediations"`
}

// Remediation represents a remediation
type Remediation struct {
	Name    string               `json:"name" yaml:"name"`
	Actions *[]RemediationAction `json:"actions" yaml:"actions"`
}

// RemediationAction represents a remediation action
type RemediationAction struct {
	Action string `json:"action" yaml:"action"`
	Value  string `json:"value" yaml:"value"`
}
