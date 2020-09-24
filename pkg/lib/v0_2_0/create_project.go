package v0_2_0

const CreateProjectTaskName = "create-project"

type CreateProjectTriggeredEventData struct {
	EventData
	// Shipyard is a base64 encoded string of the shipyard file
	Shipyard string `json:"shipyard"`
	// GitUser is the name of a git user of an upstream repository
	GitUser string `json:"gitUser,omitempty"`
	// GitToken is the authentication token for the git user
	GitToken string `json:"gitToken,omitempty"`
	// GitRemoteURL is the remote url of a repository
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`
}

type CreateProjectStartedEventData struct {
	EventData
}

type CreateProjectStatusChangedEventData struct {
	EventData
}

type CreateProjectFinishedEventData struct {
	EventData
}
