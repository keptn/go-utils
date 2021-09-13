package models

import (
	"encoding/json"
	"errors"
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

	// traceparent
	TraceParent string `json:"traceparent,omitempty"`

	// tracestate
	TraceState string `json:"tracestate,omitempty"`
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

// Validate checks whether the required properties 'time', 'type', 'id' and 'source' are defined and non-empty
func (ce *KeptnContextExtendedCE) Validate() error {
	if ce.Time.IsZero() {
		return errors.New("time must be specified")
	}
	if ce.Type == nil || *ce.Type == "" {
		return errors.New("type must be specified")
	}
	if ce.ID == "" {
		return errors.New("id must be specified")
	}
	if ce.Source == nil || *ce.Source == "" {
		return errors.New("source must be specified")
	}
	return nil
}
