package models

// ExpandedStage stage with expanded service
//
// swagger:model ExpandedStage
type ExpandedStage struct {

	// last event context
	LastEventContext *EventContextInfo `json:"lastEventContext,omitempty"`

	// services
	Services []*ExpandedService `json:"services"`

	// Stage name
	StageName string `json:"stageName,omitempty"`

	// Parent Stages
	ParentStages []string `json:"parentStages,omitempty"`
}
