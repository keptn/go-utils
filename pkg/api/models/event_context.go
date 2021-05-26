package models

// EventContext event context
type EventContext struct {

	// keptn context
	// Required: true
	KeptnContext *string `json:"keptnContext"`
}
