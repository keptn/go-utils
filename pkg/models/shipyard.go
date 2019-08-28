package models

// Shipyard defines the name, deployment strategy and test strategy of each stage
type Shipyard struct {
	Stages []struct {
		Name               string `json:"name"`
		DeploymentStrategy string `json:"deployment_strategy"`
		TestStrategy       string `json:"test_strategy,omitempty"`
	} `json:"stages"`
}
