package v0_2_0

import "net/url"

const ServiceCreateTaskName = "service.create"

type Helm struct {
	Chart   string  `json:"chart"`
	RepoURL url.URL `json:"repoURL"`
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
