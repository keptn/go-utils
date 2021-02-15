package v0_2_0

const RollbackTaskName = "rollback"

type RollbackTriggeredEventData struct {
	EventData
}

type RollbackStartedEventData struct {
	EventData
}

type RollbackFinishedEventData struct {
	EventData
}

type RollbackData struct {
}
