package v0_2_0

const GetActionTaskName = "get-action"

type GetActionTriggeredEventData struct {
	EventData
	Problem     ProblemDetails `json:"problem"`
	ActionIndex int            `json:"actionIndex"`
}

type GetActionStartedEventData struct {
	EventData
}

type GetActionFinishedEventData struct {
	EventData
	Action      ActionInfo `json:"action"`
	ActionIndex int
}
