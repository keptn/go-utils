package v0_2_0

// ProblemDetails contains information about a problem
type ProblemDetails struct {
	// ProblemTitle is the display number of the reported problem.
	ProblemTitle string `json:"problemTitle"`
	// RootCause is the root cause of the problem
	RootCause string `json:"rootCause"`
}
