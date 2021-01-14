package v0_2_0

import "encoding/json"

// ProblemDetails contains information about a problem
type ProblemDetails struct {
	// State is the state of the problem; possible values are: OPEN, RESOLVED
	State string `json:"State,omitempty jsonschema:"enum=open,enum=resolved"`
	// ProblemID is a unique system identifier of the reported problem
	ProblemID string `json:"ProblemID"`
	// ProblemTitle is the display number of the reported problem.
	ProblemTitle string `json:"ProblemTitle"`
	// ProblemDetails are all problem event details including root cause
	ProblemDetails json.RawMessage `json:"ProblemDetails"`
	// PID is a unique system identifier of the reported problem.
	PID string `json:"PID"`
	// ImpcatedEntity is an identifier of the impacted entity
	// ProblemURL is a back link to the original problem
	ProblemURL     string `json:"ProblemURL,omitempty"`
	ImpactedEntity string `json:"ImpactedEntity,omitempty"`
	// Tags is a comma separated list of tags that are defined for all impacted entities.
	Tags string `json:"Tags,omitempty"`
}
