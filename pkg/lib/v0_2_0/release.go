package v0_2_0

const ReleaseTaskName = "release"

type ReleaseTriggeredEventData struct {
	EventData
	Deployment struct {
		// DeploymentStrategy defines the used deployment strategy
		DeploymentStrategy string `json:"deploymentstrategy,omitempty"`
	} `json:"deployment"`
}

type ReleaseStartedEventData struct {
	EventData
}

type ReleaseStatusChangedEventData struct {
	EventData
}

type ReleaseFinishedEventData struct {
	EventData
	Release struct {
		// GitCommit indicates the version which should be deployed
		GitCommit string `json:"gitCommit"`
	} `json:"Release"`
}
