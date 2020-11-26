package v0_2_0

const GetSLITaskName = "get-sli"

type GetSLITriggeredEventData struct {
	EventData
	GetSLI struct {
		SLIProvider   string       `json:"sliProvider"`
		Start         string       `json:"start"`
		End           string       `json:"end"`
		Indicators    []string     `json:"indicators"`
		CustomFilters []*SLIFilter `json:"customFilters"`
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

type SLIFilter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
