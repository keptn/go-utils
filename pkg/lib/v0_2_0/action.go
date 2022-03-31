package v0_2_0

const ActionTaskName = "action"

// ActionTriggeredEventData contains information about an action.triggered event
type ActionTriggeredEventData struct {
	EventData
	// Action describes the type of action
	Action ActionInfo `json:"action"`
	// Problem contains details about the problem
	Problem ProblemDetails `json:"problem"`
}

// ActionInfo contains information about the action to be performed
type ActionInfo struct {
	// Name is the name of the action
	Name string `json:"name"`
	// Action determines the type of action to be executed
	Action string `json:"action"`
	// Description contains the description of the action
	Description string `json:"description,omitempty"`
	// Value contains the value of the action
	Value interface{} `json:"value,omitempty"`
}

// ActionStartedEventData contains information about an action.started event
type ActionStartedEventData struct {
	EventData
}

// ActionFinishedEventData contains information about the execution of an action
type ActionFinishedEventData struct {
	EventData
}
