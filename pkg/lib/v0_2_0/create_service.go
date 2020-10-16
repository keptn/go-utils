package v0_2_0

const CreateServiceTaskName = "service.create"

type CreateServiceData struct {
	ProjectName  string `json:"projectName"`
	GitRemoteURL string `json:"gitRemoteURL"`
}

type CreateServiceFinishedEventData struct {
	EventData
}
