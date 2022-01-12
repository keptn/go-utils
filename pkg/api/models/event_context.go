package models

import "encoding/json"

// EventContext event context
type EventContext struct {

	// keptn context
	// Required: true
	KeptnContext *string `json:"keptnContext"`
}

func (ec *EventContext) ToJSON() ([]byte, error) {
	if ec == nil {
		return nil, nil
	}
	return json.Marshal(ec)
}

func (ec *EventContext) FromJSON(b []byte) error {
	var res EventContext
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*ec = res
	return nil
}
