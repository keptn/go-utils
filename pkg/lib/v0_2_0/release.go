package v0_2_0

const ReleaseTaskName = "release"

type ReleaseTriggeredEventData struct {
	EventData
	Deployment DeploymentFinishedData `json:"deployment"`
}

type ReleaseStartedEventData struct {
	EventData
}

type ReleaseStatusChangedEventData struct {
	EventData
}

type ReleaseFinishedEventData struct {
	EventData
	Release ReleaseData `json:"release"`
}

type ReleaseData struct {
	// GitCommit indicates the version which should be deployed
	GitCommit string `json:"gitCommit"` //TODO this one ?
}
