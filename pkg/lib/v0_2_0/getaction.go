package v0_2_0

const GetActionTaskName = "get-action"

type GetActionData struct {
	ActionIndex int `json:"actionIndex"`
}

type GetActionTriggeredEventData struct {
	EventData
	Problem   ProblemDetails `json:"problem"`
	GetAction GetActionData  `json:"get-action"`
}

type GetActionStartedEventData struct {
	EventData
}

type GetActionFinishedEventData struct {
	EventData
	Action    ActionInfo    `json:"action"`
	GetAction GetActionData `json:"get-action"`
}
