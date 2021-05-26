package models

import (
	"encoding/json"
	"time"
)

// KeptnContextExtendedCE keptn context extended CloudEvent
type KeptnContextExtendedCE struct {

	// contenttype
	Contenttype string `json:"contenttype,omitempty"`

	// data
	// Required: true
	Data interface{} `json:"data"`

	// extensions
	Extensions interface{} `json:"extensions,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// shkeptncontext
	Shkeptncontext string `json:"shkeptncontext,omitempty"`

	// shkeptnspecversion
	Shkeptnspecversion string `json:"shkeptnspecversion,omitempty"`

	// source
	// Required: true
	Source *string `json:"source"`

	// specversion
	Specversion string `json:"specversion,omitempty"`

	// time
	// Format: date-time
	Time time.Time `json:"time,omitempty"`

	// triggeredid
	Triggeredid string `json:"triggeredid,omitempty"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// DataAs attempts to populate the provided data object with the event payload.
// data should be a pointer type.
func (ce *KeptnContextExtendedCE) DataAs(out interface{}) error {
	bytes, err := json.Marshal(ce.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}
