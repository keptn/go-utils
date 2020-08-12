package v0_2_0

type StatusType string

const (
	StatusSucceeded StatusType = "succeeded"
	StatusErrored   StatusType = "errored"
	StatusUnknown   StatusType = "unknown"
)
