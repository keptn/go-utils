package v0_2_0

const DeploymentTaskName = "deployment"

type DeploymentTriggeredEventData struct {
	EventData

	ConfigurationChange ConfigurationChange         `json:"configurationChange"`
	Deployment          DeploymentTriggeredProperty `json:"deployment"`
}

type DeploymentTriggeredProperty struct {
	// DeploymentStrategy defines the used deployment strategy
	DeploymentStrategy string `json:"deploymentstrategy,omitempty"`
}

type DeploymentData struct {
	// DeploymentURILocal contains the local URL
	DeploymentURIsLocal []string `json:"deploymentURIsLocal,omitempty"`
	// DeploymentURIPublic contains the public URL
	DeploymentURIsPublic []string `json:"deploymentURIsPublic,omitempty"`
	// DeploymentNames gives the names of the deployments
	DeploymentNames []string `json:"deploymentNames"`
	// GitCommit indicates the version which should be deployed
	GitCommit string `json:"gitCommit"`
}

type DeploymentStartedEventData struct {
	EventData
}

type DeploymentStatusChangedEventData struct {
	EventData
}

type DeploymentFinishedEventData struct {
	EventData
	Deployment DeploymentData `json:"deployment"`
}
