package v0_2_0

const ConfigureMonitoringTaskName = "configure-monitoring"

// ConfigureMonitoringTriggeredEventData godoc
type ConfigureMonitoringTriggeredEventData struct {
	EventData

	ConfigureMonitoring ConfigureMonitoringTriggeredParams `json:"configureMonitoring"`
}

// ConfigureMonitoringTriggeredParams godoc
type ConfigureMonitoringTriggeredParams struct {
	Type string `json:"type"`
}

// ConfigureMonitoringStartedEventData godoc
type ConfigureMonitoringStartedEventData struct {
	EventData
}

// ConfigureMonitoringFinishedEventData godoc
type ConfigureMonitoringFinishedEventData struct {
	EventData
}
