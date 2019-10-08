package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetInvalidDeploymentStrategy tests whether an error is returned
// if an invalid test strategy is passed to GetDeploymentStrategy
func TestGetInvalidDeploymentStrategy(t *testing.T) {

	_, err := GetDeploymentStrategy("invalidStrategy")
	assert.Error(t, err)
}
