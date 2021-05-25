package keptn

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
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

// TestGetEvent tests whether the GetEvent function returns the latest Keptn event
// of type sh.keptn.events.evaluation-done and from Keptn context: 8929e5e5-3826-488f-9257-708bfa974909
func TestGetEventStatusOK(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Expect GET request")
		assert.Equal(t, r.URL.EscapedPath(), "/event", "Expect /event endpoint")
		w.WriteHeader(http.StatusOK)
		eventList := `{
			"events":[
				{"contenttype":"application/json",
					"data":{"deploymentstrategy":"blue_green_service",
					"evaluationdetails":{"result": "pass"},
					"evaluationpassed":true,
					"project":"sockshop",
					"service":"carts",
					"stage":"production",
					"teststrategy":""},
				"id":"aaa50752-ab33-493b-8b28-3548f7960f80",
				"source":"pitometer-service",
				"specversion":"0.2",
				"time":"2019-10-21T14:12:48.000Z",
				"type":"sh.keptn.events.evaluation-done",
				"shkeptncontext":"8929e5e5-3826-488f-9257-708bfa974909"},
				
				{"contenttype":"application/json",
					"data":{"deploymentstrategy":"blue_green_service",
					"evaluationdetails":{"result": "pass"},
					"evaluationpassed":true,
					"project":"sockshop",
					"service":"carts",
					"stage":"staging",
					"teststrategy":"performance"},
				"id":"573610d2-3643-4513-9a8e-df7c6614356f",
				"source":"pitometer-service",
				"specversion":"0.2",
				"time":"2019-10-21T14:10:05.000Z",
				"type":"sh.keptn.events.evaluation-done",
				"shkeptncontext":"8929e5e5-3826-488f-9257-708bfa974909"},
				
				{"contenttype":"application/json",
					"data":{"deploymentstrategy":"direct",
					"evaluationdetails":{"result": "pass"},
					"evaluationpassed":true,
					"project":"sockshop",
					"service":"carts",
					"stage":"dev",
					"teststrategy":"functional"},
				"id":"a46be431-b45b-4f18-bf74-73fc7d2da062",
				"source":"pitometer-service",
				"specversion":"0.2",
				"time":"2019-10-21T14:04:25.000Z",
				"type":"sh.keptn.events.evaluation-done",
				"shkeptncontext":"8929e5e5-3826-488f-9257-708bfa974909"}],
				
				"pageSize":10,
				"totalCount":3
			}`

		w.Write([]byte(eventList))
	})

	httpClient, teardown := testingHTTPClient(handler)
	defer teardown()

	eventHandler := NewEventHandler("https://localhost")
	eventHandler.HTTPClient = httpClient
	cloudEvent, errObj := eventHandler.GetEvent("8929e5e5-3826-488f-9257-708bfa974909", "sh.keptn.events.evaluation-done")

	if cloudEvent == nil {
		t.Error("no Keptn event returned")
	}

	// check whether the last event is returned
	if cloudEvent.Time != "2019-10-21T14:12:48.000Z" {
		t.Error("did not receive the latest event")
	}

	if errObj != nil {
		t.Errorf("an error occurred %v", errObj.Message)
	}
}

// TestGetEvent tests whether the GetEvent function returns no event found when no event is available
func TestGetEventStatusOKNoEvent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Expect GET request")
		assert.Equal(t, r.URL.EscapedPath(), "/event", "Expect /event endpoint")
		w.WriteHeader(http.StatusOK)
		eventList := `{
			"events":[],
				"pageSize":10,
				"totalCount":0
			}`

		w.Write([]byte(eventList))
	})

	httpClient, teardown := testingHTTPClient(handler)
	defer teardown()

	eventHandler := NewEventHandler("https://localhost")
	eventHandler.HTTPClient = httpClient
	cloudEvent, errObj := eventHandler.GetEvent("8929e5e5-3826-488f-9257-708bfa974909", "sh.keptn.events.evaluation-done")

	if cloudEvent != nil {
		t.Error("do not expect a Keptn Cloud event")
	}

	if errObj == nil {
		t.Errorf("an error occurred %v", errObj.Message)
	}

	if *errObj.Message != "No Keptn sh.keptn.events.evaluation-done event found for context: 8929e5e5-3826-488f-9257-708bfa974909" {
		t.Error("response message has changed")
	}
}
