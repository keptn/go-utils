package v0_2_0

const GetSLITaskName = "get-sli"

type GetSLITriggeredEventData struct {
	EventData
	GetSLI struct {
		// SLIProvider defines the name of the monitoring solution that provides the SLIs
		SLIProvider string `json:"sliProvider"`
		// Start defines the start timestamp
		Start string `json:"start"`
		// End defines the end timestamp
		End string `json:"end"`
		// Indicators defines the SLI names
		Indicators []string `json:"indicators"`
		// CustomFilters defines filters on the SLIs
		CustomFilters []*SLIFilter `json:"customFilters"`
	} `json:"get-sli"`
}

type GetSLIStartedEventData struct {
	EventData
}

type GetSLIFinishedEventData struct {
	EventData
	GetSLI struct {
		// Start defines the start timestamp
		Start string `json:"start"`
		// End defines the end timestamp
		End string `json:"end"`
		// IndicatorValues defines the fetched SLI values
		IndicatorValues []*SLIResult `json:"indicatorValues"`
	} `json:"get-sli"`
}

type SLIFilter struct {
	// Key defines the key of the SLI filter
	Key string `json:"key"`
	// Value defines the value of the SLI filter
	Value string `json:"value"`
}
