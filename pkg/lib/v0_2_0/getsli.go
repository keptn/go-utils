package v0_2_0

const GetSLITaskName = "get-sli"

type GetSLITriggeredEventData struct {
	EventData
	GetSLI GetSLI `json:"get-sli"`
}

type GetSLI struct {
	// SLIProvider defines the name of the monitoring solution that provides the SLIs
	SLIProvider string `json:"sliProvider"`
	// Start defines the start timestamp
	Start string `json:"start"`
	// End defines the end timestamp
	End string `json:"end"`
	// Indicators defines the SLI names
	Indicators []string `json:"indicators,omitempty"`
	// CustomFilters defines filters on the SLIs
	CustomFilters []*SLIFilter `json:"customFilters,omitempty"`
}

type GetSLIStartedEventData struct {
	EventData
}

type GetSLIFinishedEventData struct {
	EventData
	GetSLI GetSLIFinished `json:"get-sli"`
}

type GetSLIFinished struct {
	// Start defines the start timestamp
	Start string `json:"start"`
	// End defines the end timestamp
	End string `json:"end"`
	// IndicatorValues defines the fetched SLI values
	IndicatorValues []*SLIResult `json:"indicatorValues,omitempty"`
}
type SLIFilter struct {
	// Key defines the key of the SLI filter
	Key string `json:"key"`
	// Value defines the value of the SLI filter
	Value string `json:"value"`
}
