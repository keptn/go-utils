package v2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type StatusBody struct {
	Status string `json:"status"`
}

// ToJSON converts object to JSON string
func (s *StatusBody) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
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

func RunHealthEndpoint(port string) {

	http.HandleFunc("/health", healthHandler)
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
