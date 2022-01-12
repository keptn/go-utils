package models

import "encoding/json"

// Error error
type Error struct {

	// Error code
	Code int64 `json:"code,omitempty"`

	// Error message
	// Required: true
	Message *string `json:"message"`
}

func (e Error) GetMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}

func (e *Error) ToJSON() ([]byte, error) {
	if e == nil {
		return nil, nil
	}
	return json.Marshal(e)
}

func (e *Error) FromJSON(b []byte) error {
	var res Error
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*e = res
	return nil
}
