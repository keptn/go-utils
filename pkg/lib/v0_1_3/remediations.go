package v0_1_3

///// v0.1.3 Remediation Spec ///////

// Remediation describes a remediation specification according to Keptn spec 0.1.3
type Remediation struct {
	ApiVersion string              `json:"apiVersion" yaml:"version"`
	Kind       string              `json:"kind" yaml:"kind"`
	Metadata   RemediationMetadata `json:"metadata" yaml:"metadata"`
	Spec       RemediationSpec     `json:"spec" yaml:"spec"`
}

// RemediationMetadata describes Remediation metadata
type RemediationMetadata struct {
	Name string `json:"name" yaml:"name"`
}

// RemediationActionsOnOpen describes an action which is executed when a problem.open occurred
type RemediationActionsOnOpen struct {
	Name        string      `json:"name" yaml:"name"`
	Action      string      `json:"action" yaml:"action"`
	Description string      `json:"description" yaml:"description"`
	Value       interface{} `json:"value" yaml:"value"`
}

// RemediationMap maps a problem to a list of actions which are executed when a problem.open occurred
type RemediationMap struct {
	ProblemType   string                     `json:"problemType" yaml:"problemType"`
	ActionsOnOpen []RemediationActionsOnOpen `json:"actionsOnOpen" yaml:"actionsOnOpen"`
}

// RemediationSpec contains a list of remediations
type RemediationSpec struct {
	Remediations []RemediationMap `json:"remediations" yaml:"remediations"`
}
