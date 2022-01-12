package models

import "encoding/json"

// EventContext event context
type EventContext struct {

	// keptn context
	// Required: true
	KeptnContext *string `json:"keptnContext"`
}

// ToJSON converts object to JSON string
func (ec *EventContext) ToJSON() ([]byte, error) {
	return json.Marshal(ec)
}

// FromJSON converts JSON string to object
func (ec *EventContext) FromJSON(b []byte) error {
	var res EventContext
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*ec = res
	return nil
}
