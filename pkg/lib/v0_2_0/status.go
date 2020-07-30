package keptn

type StatusType string

const (
	StatusSucceeded StatusType = "succeeded"
	StatusErrored   StatusType = "errored"
	StatusUnknown   StatusType = "unknown"
)
