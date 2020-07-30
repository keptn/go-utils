package keptn

const DeploymentTaskName = "deployment"

type DeploymentTriggeredEventData struct {
	EventData

	ConfigurationChange struct {
		GitCommit string `json:"gitCommit"`
	} `json:"configurationChange"`

	Deployment DeploymentStartedData `json:"deployment"`
}

type DeploymentStartedData struct {
	// DeploymentURILocal contains the local URL
	DeploymentURIsLocal []string `json:"deploymentURIsLocal,omitempty"`
	// DeploymentURIPublic contains the public URL
	DeploymentURIsPublic []string `json:"deploymentURIsPublic,omitempty"`
	// DeploymentNames gives the names of the deployments
	DeploymentNames []string `json:"deploymentNames"`
	// GitCommit indicates the version which was deployed
	GitCommit string `json:"gitCommit"`
}

type DeploymentFinishedData struct {
	DeploymentStartedData
	// GitCommit indicates the version which was deployed
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
	Deployment DeploymentFinishedData `json:"deployment"`
}
