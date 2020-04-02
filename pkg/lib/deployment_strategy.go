package keptn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// DeploymentStrategy describes how a keptn-managed service is deployed
type DeploymentStrategy int

const (
	// Direct stores the chart which results in the
	Direct DeploymentStrategy = iota + 1

	// Duplicate generates a second chart in order to duplicate the deployments
	Duplicate
)

func (s DeploymentStrategy) String() string {
	return deploymentStrategyToString[s]
}

// GetDeploymentStrategy tries to parse the deployment strategy into the enum
// If the provided deployment strategy is unsupported, an error is returned
func GetDeploymentStrategy(deploymentStrategy string) (DeploymentStrategy, error) {
	if val, ok := deploymentStrategyToID[deploymentStrategy]; ok {
		return val, nil
	}

	return DeploymentStrategy(-1), fmt.Errorf("The deployment strategy %s is invalid", deploymentStrategy)
}

var deploymentStrategyToString = map[DeploymentStrategy]string{
	Direct:    "direct",
	Duplicate: "duplicate",
}

var deploymentStrategyToID = map[string]DeploymentStrategy{
	"direct":    Direct,
	"duplicate": Duplicate,
}

// MarshalJSON marshals the enum as a quoted json string
func (s DeploymentStrategy) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(deploymentStrategyToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *DeploymentStrategy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = deploymentStrategyToID[strings.ToLower(j)]
	return nil
}
