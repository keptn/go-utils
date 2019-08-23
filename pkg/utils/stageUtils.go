package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/keptn/go-utils/pkg/models"
)

// StageHandler handles stages
type StageHandler struct {
	BaseURL string
}

// NewStageHandler returns a new StageHandler
func NewStageHandler(baseURL string) *StageHandler {
	return &StageHandler{
		BaseURL: baseURL,
	}
}

// GetAllStages returns a list of all stages.
func (r *StageHandler) GetAllStages(project string) ([]*models.Stage, error) {

	stages := []*models.Stage{}

	nextPageKey := ""
	for {
		url, err := url.Parse("http://" + r.BaseURL + "/v1/project/" + project + "/stage")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var received models.Stages
		err = json.Unmarshal(body, &received)
		if err != nil {
			return nil, err
		}
		stages = append(stages, received.Stages...)

		if received.NextPageKey == "" {
			break
		}
		nextPageKey = received.NextPageKey
	}
	return stages, nil
}
