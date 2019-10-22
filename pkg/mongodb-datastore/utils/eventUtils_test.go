package utils

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/magiconair/properties/assert"
)

// Helper function to build a test client with a httptest server
func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	server := httptest.NewServer(handler)

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}

	return client, server.Close
}

func TestGetEvent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Expect GET request")
		assert.Equal(t, r.URL.EscapedPath(), "/event", "Expect /event endpoint")
		w.WriteHeader(http.StatusOK) // 200 - SatusOk

		response := `
		
		`

		w.Write()
	})

	httpClient, teardown := testingHTTPClient(handler)
	defer teardown()

	eventHandler := NewEventHandler("https://localhost")
	eventHandler.HTTPClient = httpClient

	cloudEvent, errObj := eventHandler.GetEvent("8929e5e5-3826-488f-9257-708bfa974909", keptnevents.EvaluationDoneEventType)

	if cloudEvent == nil {
		t.Error("No CloudEvent returned")
	}

	if errObj != nil {
		t.Errorf("An error occured %v", errObj.Message)
	}
}
