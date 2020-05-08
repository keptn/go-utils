package keptn

// Shipyard defines the name, deployment strategy and test strategy of each stage
type Shipyard struct {
	Stages []struct {
		Name                string `json:"name" yaml:"name"`
		DeploymentStrategy  string `json:"deployment_strategy" yaml:"deployment_strategy"`
		TestStrategy        string `json:"test_strategy,omitempty" yaml:"test_strategy"`
		RemediationStrategy string `json:"remediation_strategy,omitempty" yaml:"remediation_strategy"`
		ApprovalStrategy    *struct {
			Pass    ApprovalStrategy `json:"pass,omitempty" yaml:"pass"`
			Warning ApprovalStrategy `json:"warning,omitempty" yaml:"warning"`
		} `json:"approval_strategy,omitempty" yaml:"approval_strategy"`
	} `json:"stages" yaml:"stages"`
}
