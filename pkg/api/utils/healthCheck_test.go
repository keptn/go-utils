package api

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestRunHealthEndpoint(t *testing.T) {
	go RunHealthEndpoint("8080")

	require.Eventually(t, func() bool {
		get, err := http.Get("http://localhost:8080/health")
		if err != nil {
			return false
		}
		if get.StatusCode != http.StatusOK {
			return false
		}
		return true
	}, 2*time.Second, 50*time.Millisecond)
}

func TestRunHealthEndpoint_WithReadinessCondition(t *testing.T) {
	ready := false
	go RunHealthEndpoint("8080", WithReadinessConditionFunc(func() bool {
		return ready
	}))

	require.Eventually(t, func() bool {
		get, err := http.Get("http://localhost:8080/health")
		if err != nil {
			return false
		}
		if get.StatusCode != http.StatusPreconditionFailed {
			return false
		}
		return true
	}, 2*time.Second, 50*time.Millisecond)

	ready = true

	require.Eventually(t, func() bool {
		get, err := http.Get("http://localhost:8080/health")
		if err != nil {
			return false
		}
		if get.StatusCode != http.StatusOK {
			return false
		}
		return true
	}, 2*time.Second, 50*time.Millisecond)
}
