package models

import (
	"encoding/json"
	"github.com/keptn/go-utils/pkg/common/strutils"
)

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

// ToJSON converts object to JSON string
func (e *Error) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON converts JSON string to object
func (e *Error) FromJSON(b []byte) error {
	var res Error
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*e = res
	if e.Message == nil {
		e.Message = strutils.Stringp("")
	}
	return nil
}
