package v0_2_0

const RemediationTaskName = "remediation"

// RemediationTriggeredEventData is a CloudEvent for triggering remediations
type RemediationTriggeredEventData struct {
	EventData

	// Problem contains details about the problem
	Problem ProblemDetails `json:"problem"`
}

// RemediationTriggeredEventType is a CloudEvent to indicate the start of a remediation
type RemediationStartedEventData struct {
	EventData
}

// RemediationStatus describes the result and status of a remediation
type RemediationStatusChangedEventData struct {
	EventData

	// Remediation indicates the result
	Remediation Remediation `json:"remediation"`
}

// RemediationStatus describes the status of a remediation
type Remediation struct {
	// ActionIndex is the index of the action
	ActionIndex int `json:"actionIndex"`
	// ActionName is the name of the action
	ActionName string `json:"actionName"`
}

// RemediationFinishedEventData describes a finished remediation
type RemediationFinishedEventData struct {
	EventData
}
