package v0_2_0

const ProjectDeleteTaskName = "project.delete"

type ProjectDeleteData struct {
}

type ProjectDeleteStartedEventData struct {
	EventData
}

type ProjectDeleteFinishedEventData struct {
	EventData
}
