package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAbortSequence(t *testing.T) {

}

func TestSequenceControlHandler_AbortSequence(t *testing.T) {
	tests := []struct {
		name    string
		Handler http.HandlerFunc
		params  SequenceControlParams
		wantErr bool
	}{
		{"test abort sequence - empty params",
			nil,
			SequenceControlParams{},
			true},
		{"test abort sequence - missing project",
			nil,
			SequenceControlParams{
				Project:      "",
				KeptnContext: "c1",
				Stage:        "s1",
				State:        "s1",
			}, true},
		{"test abort sequence - missing context",
			nil,
			SequenceControlParams{
				Project:      "p1",
				KeptnContext: "",
				Stage:        "s1",
				State:        "s1",
			}, true},
		{"test abort sequence - missing state",
			nil,
			SequenceControlParams{
				Project:      "p1",
				KeptnContext: "c1",
				Stage:        "s1",
				State:        "",
			}, true},
		{"test abort sequence - valid params",
			func(writer http.ResponseWriter, request *http.Request) {
				assert.Equal(t, "/v1/sequence/p1/c1/control", request.RequestURI)
				payload, _ := io.ReadAll(request.Body)
				request.Body.Close()

				params := &SequenceControlBody{}
				json.Unmarshal(payload, params)
				assert.Equal(t, "stg1", params.Stage)
				assert.Equal(t, "stt1", params.State)
			},
			SequenceControlParams{
				Project:      "p1",
				KeptnContext: "c1",
				Stage:        "stg1",
				State:        "stt1",
			}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.Handler)
			defer ts.Close()
			s := NewSequenceControlHandler(ts.URL)
			if err := s.AbortSequence(tt.params); (err != nil) != tt.wantErr {
				t.Errorf("AbortSequence() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
