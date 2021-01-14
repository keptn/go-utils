package v0_2_0

const TestTaskName = "test"

type TestTriggeredEventData struct {
	EventData

	Test TestTriggeredDetails `json:"test"`

	Deployment TestTriggeredDeploymentDetails `json:"deployment"`
}
type TestTriggeredDetails struct {
	// TestStrategy is the testing strategy and is defined in the shipyard
	TestStrategy string `json:"teststrategy" jsonschema:"enum=real-user,enum=functional,enum=performance,enum=healthcheck"`
}

type TestTriggeredDeploymentDetails struct {
	// DeploymentURILocal contains the local URL
	DeploymentURIsLocal []string `json:"deploymentURIsLocal"`
	// DeploymentURIPublic contains the public URL
	DeploymentURIsPublic []string `json:"deploymentURIsPublic,omitempty"`
}

type TestStartedEventData struct {
	EventData
}

type TestStatusChangedEventData struct {
	EventData
}

type TestFinishedEventData struct {
	EventData
	Test TestFinishedDetails `json:"test"`
}

type TestFinishedDetails struct {
	// Start indicates the starting timestamp of the tests
	Start string `json:"start"`
	// End indicates the end timestamp of the tests
	End string `json:"end"`
	// GitCommit indicates the version which should be deployed
	GitCommit string `json:"gitCommit"`
}
