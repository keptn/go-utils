package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeploymentStrategy(t *testing.T) {

	_, err := GetDeploymentStrategy("invalidStrategy")
	assert.Error(t, err)
}
