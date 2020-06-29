package keptn

// Remediations contains remediation definitions for a service
// Deprecated: A new spec for Remediations is available
type Remediations struct {
	Remediations []*Remediation `json:"remediations" yaml:"remediations"`
}

// Remediation represents a remediation
// Deprecated: A new spec for Remediation is available
type Remediation struct {
	Name    string               `json:"name" yaml:"name"`
	Actions []*RemediationAction `json:"actions" yaml:"actions"`
}

// RemediationAction represents a remediation action
// Deprecated: A new spec for RemediationAction is available
type RemediationAction struct {
	Action string `json:"action" yaml:"action"`
	Value  string `json:"value" yaml:"value"`
}
