package models

// Error error
type Error struct {

	// Error code
	Code int64 `json:"code,omitempty"`

	// Error message
	// Required: true
	Message *string `json:"message"`
}
