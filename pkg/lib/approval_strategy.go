package keptn

import (
	"bytes"
	"encoding/json"
	"strings"
)

// ApprovalStrategy is used in the shipyard for either automatic or manual approvals
type ApprovalStrategy int

const (
	// Automatic: A step is approved in an automatic fashion
	Automatic ApprovalStrategy = iota
	// Manual: A step is approved in a manual fashion
	Manual
)

func (a ApprovalStrategy) String() string {
	return approvalStrategyToString[a]
}

var approvalStrategyToString = map[ApprovalStrategy]string{
	Automatic: "automatic",
	Manual:    "manual",
}

var approvalStrategyToID = map[string]ApprovalStrategy{
	"automatic": Automatic,
	"manual":    Manual,
}

// MarshalJSON marshals the enum as a quoted json string
func (s ApprovalStrategy) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(approvalStrategyToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *ApprovalStrategy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = approvalStrategyToID[strings.ToLower(j)]
	return nil
}
