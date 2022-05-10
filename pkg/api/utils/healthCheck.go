package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ReadinessConditionFunc func() bool

type StatusBody struct {
	Status string `json:"status"`
}

// ToJSON converts object to JSON string
func (s *StatusBody) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

type HealthHandlerOption func(h *healthHandler)

func WithReadinessConditionFunc(rc ReadinessConditionFunc) HealthHandlerOption {
	return func(h *healthHandler) {
		h.readinessConditionFunc = rc
	}
}

type healthHandler struct {
	readinessConditionFunc ReadinessConditionFunc
}

func newHealthHandler(opts ...HealthHandlerOption) *healthHandler {
	h := &healthHandler{}
	for _, o := range opts {
		o(h)
	}
	return h
}

func (h *healthHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	ready := true
	if h.readinessConditionFunc != nil {
		ready = h.readinessConditionFunc()
	}

	if !ready {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}
	status := StatusBody{Status: "OK"}

	body, err := status.ToJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("content-type", "application/json")

	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func RunHealthEndpoint(port string, opts ...HealthHandlerOption) {
	h := newHealthHandler(opts...)
	http.HandleFunc("/health", h.healthCheck)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Println(err)
	}
}

func HealthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
	}
}
