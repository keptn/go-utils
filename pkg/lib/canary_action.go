package keptn

import (
	"bytes"
	"encoding/json"
	"strings"
)

// CanaryAction describes the excution of a canary release
type CanaryAction int

const (
	// Set is used for setting a new traffic weight for the canary
	Set CanaryAction = iota
	// Promote is used for promoting the canary
	Promote
	// Discard is used for discarding the canary
	Discard
)

func (s CanaryAction) String() string {
	return canaryActionToString[s]
}

var canaryActionToString = map[CanaryAction]string{
	Set:     "set",
	Promote: "promote",
	Discard: "discard",
}

var canaryActionToID = map[string]CanaryAction{
	"set":     Set,
	"promote": Promote,
	"discard": Discard,
}

// MarshalJSON marshals the enum as a quoted json string
func (s CanaryAction) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(canaryActionToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *CanaryAction) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = canaryActionToID[strings.ToLower(j)]
	return nil
}
