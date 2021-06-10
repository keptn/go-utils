package models

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
