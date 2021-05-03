package v0_2_0

type GetActionTriggeredEventData struct {
	EventData
	ProblemDetails ProblemDetails `json:"problemDetails"`
	ActionIndex    int            `json:"actionIndex"`
}

type GetActionStartedEventData struct {
	EventData
}

type GetActionFinishedEventData struct {
	EventData
	Action      ActionInfo `json:"action"`
	ActionIndex int
}
