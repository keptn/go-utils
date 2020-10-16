package v0_2_0

const ServiceDeleteTaskName = "service.delete"

type ServiceDeleteStartedEventData struct {
	EventData
}

type ServiceDeleteStatusChangedEventData struct {
	EventData
}

type ServiceDeleteFinishedEventData struct {
	EventData
}
