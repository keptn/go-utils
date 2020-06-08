package keptn

// Remediations contains remediation definitions for a service
type Remediations struct {
	Remediations []*Remediation `json:"remediations" yaml:"remediations"`
}

// Remediation represents a remediation
type Remediation struct {
	Name    string               `json:"name" yaml:"name"`
	Actions []*RemediationAction `json:"actions" yaml:"actions"`
}

// RemediationAction represents a remediation action
type RemediationAction struct {
	Action string `json:"action" yaml:"action"`
	Value  string `json:"value" yaml:"value"`
}

///// v0.2.0 Remediation Spec ///////

// RemediationV02 describes a remediation specification according to Keptn spec 0.7.0
type RemediationV02 struct {
	Version  string                 `json:"version" yaml:"version"`
	Kind     string                 `json:"kind" yaml:"kind"`
	Metadata RemediationV02Metadata `json:"metadata" yaml:"metadata"`
	Spec     RemediationV02Spec     `json:"spec" yaml:"spec"`
}

// RemediationV02Metadata describes Remediation metadata
type RemediationV02Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type RemediationV02ActionsOnOpen struct {
	Name        string      `json:"name" yaml:"name"`
	Action      string      `json:"action" yaml:"action"`
	Description string      `json:"description" yaml:"description"`
	Value       interface{} `json:"value" yaml:"value"`
}
type RemediationV02Remediations struct {
	ProblemType   string                        `json:"problemType" yaml:"problemType"`
	ActionsOnOpen []RemediationV02ActionsOnOpen `json:"actionsOnOpen" yaml:"actionsOnOpen"`
}
type RemediationV02Spec struct {
	Remediations []RemediationV02Remediations `json:"remediations" yaml:"remediations"`
}
