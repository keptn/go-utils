package v0_2_0

const GetSliTaskName = "get-sli"

type GetSLITriggeredEventData struct {
	EventData
	GetSLI struct {
		SLIProvider string   `json:"sliProvider"`
		Start       string   `json:"start"`
		End         string   `json:"end"`
		Indicators  []string `json:"indicators"`
	} `json:"get-sli"`
}

type GetSLIStartedEventData struct {
	EventData
}

type GetSLIFinishedEventData struct {
	EventData
	GetSLI struct {
		Start           string       `json:"start"`
		End             string       `json:"end"`
		IndicatorValues []*SLIResult `json:"indicatorValues"`
	} `json:"get-sli"`
}
