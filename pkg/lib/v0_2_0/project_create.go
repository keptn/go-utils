package v0_2_0

const ProjectCreateTaskName = "project.create"

type ProjectCreateData struct {
	ProjectName  string `json:"projectName"`
	GitRemoteURL string `json:"gitRemoteURL"`
	Shipyard     string `json:"shipyard"`
}

type ProjectCreateStartedEventData struct {
	EventData
}

type ProjectCreateFinishedEventData struct {
	EventData
	CreatedProject ProjectCreateData `json:"createdProject"`
}
