package v0_2_0

const ApprovalTaskName = "approval"

type ApprovalTriggeredEventData struct {
	EventData

	Approval struct {
		Pass    string `json:"pass"`
		Warning string `json:"warning"`
	} `json:"approval"`
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
