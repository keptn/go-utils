package v0_2_0

const ApprovalTaskName = "approval"

// ApprovalAutomatic indicates an automatic approval strategy
const ApprovalAutomatic = "automatic"

// ApprovalManual indicates a manual approval strategy
const ApprovalManual = "manual"

type ApprovalTriggeredEventData struct {
	EventData

	// Approval contains information about the approval strategy
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
