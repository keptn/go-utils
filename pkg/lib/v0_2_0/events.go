package keptn

// GetTriggeredEventType returns for the given task the name of the triggered event type
func GetTriggeredEventType(task string) string {
	return "sh.keptn.event." + task + ".triggered"
}

// GetStartedEventType returns for the given task the name of the started event type
func GetStartedEventType(task string) string {
	return "sh.keptn.event." + task + ".started"
}

// GetStatusChangedEventType returns for the given task the name of the status.changed event type
func GetStatusChangedEventType(task string) string {
	return "sh.keptn.event." + task + ".status.changed"
}

// GetFinishedEventType returns for the given task the name of the finished event type
func GetFinishedEventType(task string) string {
	return "sh.keptn.event." + task + ".finished"
}

// EventData contains mandatory fields of all Keptn CloudEvents
type EventData struct {
	Project string            `json:"project"`
	Stage   string            `json:"stage"`
	Service string            `json:"service"`
	Labels  map[string]string `json:"labels"`

	Status  StatusType `json:"status"`
	Result  ResultType `json:"result"`
	Message string     `json:"message"`
}
