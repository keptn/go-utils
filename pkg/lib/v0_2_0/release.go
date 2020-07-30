package keptn

const ReleaseTaskName = "release"

type ReleaseTriggeredEventData struct {
	EventData
}

type ReleaseStartedEventData struct {
	EventData
}

type ReleaseStatusChangedEventData struct {
	EventData
}

type ReleaseFinishedEventData struct {
	EventData
}
