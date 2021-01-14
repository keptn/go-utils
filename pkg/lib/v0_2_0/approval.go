package v0_2_0

const ApprovalTaskName = "approval"

// ApprovalAutomatic indicates an automatic approval strategy
const ApprovalAutomatic = "automatic"

// ApprovalManual indicates a manual approval strategy
const ApprovalManual = "manual"

type ApprovalTriggeredEventData struct {
	EventData
	// Approval contains information about the approval strategy
	Approval Approval `json:"approval"`
}

type Approval struct {
	Pass    string `json:"pass" jsonschema:"enum=automatic,enum=manual"`
	Warning string `json:"warning" jsonschema:"enum=automatic,enum=manual"`
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
