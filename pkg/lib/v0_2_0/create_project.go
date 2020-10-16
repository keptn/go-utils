package v0_2_0

const CreateProjectTaskName = "project.create"

type CreateProjectData struct {
	ProjectName  string `json:"projectName"`
	GitRemoteURL string `json:"gitRemoteURL"`
	Shipyard     string `json:"shipyard"`
}

type CreateProjectFinishedEventData struct {
	EventData
	Project CreateProjectData `json:"project"`
}
