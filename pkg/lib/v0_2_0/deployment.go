package v0_2_0

const DeploymentTaskName = "deployment"

type DeploymentTriggeredEventData struct {
	EventData
	ConfigurationChange ConfigurationChange     `json:"configurationChange"`
	Deployment          DeploymentTriggeredData `json:"deployment"`
}

// DeploymentTriggeredData contains the data associated with a .deployment.triggered event
type DeploymentTriggeredData struct {
	// DeploymentURILocal contains the local URL
	DeploymentURIsLocal []string `json:"deploymentURIsLocal"`
	// DeploymentURIPublic contains the public URL
	DeploymentURIsPublic []string `json:"deploymentURIsPublic,omitempty"`
}

type DeploymentStartedEventData struct {
	EventData
}

type DeploymentStatusChangedEventData struct {
	EventData
}

type DeploymentFinishedEventData struct {
	EventData
	Deployment DeploymentFinishedData `json:"deployment"`
}

// DeploymentFinishedData contains the data associated with a .deployment.finished event
type DeploymentFinishedData struct {
	// DeploymentStrategy defines the used deployment strategy
	DeploymentStrategy string `json:"deploymentstrategy" jsonschema:"enum=direct,enum=blue_green_service,enum=user_managed"`
	// DeploymentURILocal contains the local URL
	DeploymentURIsLocal []string `json:"deploymentURIsLocal"`
	// DeploymentURIPublic contains the public URL
	DeploymentURIsPublic []string `json:"deploymentURIsPublic,omitempty"`
	// DeploymentNames gives the names of the deployments
	DeploymentNames []string `json:"deploymentNames"`
	// GitCommit indicates the version which should be deployed
	GitCommit string `json:"gitCommit"`
}
