package models

import "encoding/json"

// Approval approval
type Approval struct {

	// ID of the event
	EventID string `json:"eventId,omitempty"`

	// image
	Image string `json:"image,omitempty"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty"`

	// tag
	Tag string `json:"tag,omitempty"`

	// Time of the event
	Time string `json:"time,omitempty"`
}

func (a *Approval) ToJSON() ([]byte, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *Approval) FromJSON(b []byte) error {
	var res Approval
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*a = res
	return nil
}
