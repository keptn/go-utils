package v0_2_0

const ServiceCreateTaskName = "service.create"

type Helm struct {
	Chart string `json:"chart"`
}

type ServiceCreateStartedEventData struct {
	EventData
}

type ServiceCreateStatusChangedEventData struct {
	EventData
}

type ServiceCreateFinishedEventData struct {
	EventData
	Helm Helm `json:"helm"`
}
