package v0_2_0

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/keptn/go-utils/pkg/lib/keptn"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

func TestKeptn_SendCloudEvent(t *testing.T) {
	failOnFirstTry := true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			if failOnFirstTry {
				failOnFirstTry = false
				w.WriteHeader(500)
				w.Write([]byte(`{}`))
			}
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	eventNew := cloudevents.NewEvent()
	eventNew.SetSource("https://test-source")
	eventNew.SetID(uuid.New().String())
	eventNew.SetType("sh.keptn.events.test")
	eventNew.SetExtension("shkeptncontext", "test-context")
	eventNew.SetData(cloudevents.ApplicationJSON, map[string]string{"project": "sockshop"})

	k := Keptn{
		KeptnBase: keptn.KeptnBase{
			EventBrokerURL: ts.URL,
		},
	}

	if err := k.SendCloudEvent(eventNew); err != nil {
		t.Errorf("SendCloudEvent() error = %v", err)
	}
}
