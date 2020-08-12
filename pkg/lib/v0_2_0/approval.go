package v0_2_0

const ApprovalTaskName = "approval"

type ApprovalTriggeredEventData struct {
	EventData
}

type ApprovalStartedEventData struct {
	EventData
}

type ApprovalStatusChangedEventData struct {
	EventData
}

type ApprovalFinishedEventData struct {
	EventData
}
