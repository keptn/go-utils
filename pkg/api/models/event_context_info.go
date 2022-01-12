package models

import "encoding/json"

// EventContextInfo event context info
type EventContextInfo struct {

	// ID of the event
	EventID string `json:"eventId,omitempty"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty"`

	// Time of the event
	Time string `json:"time,omitempty"`
}

func (ec *EventContextInfo) ToJSON() ([]byte, error) {
	if ec == nil {
		return nil, nil
	}
	return json.Marshal(ec)
}

func (ec *EventContextInfo) FromJSON(b []byte) error {
	var res EventContextInfo
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*ec = res
	return nil
}
