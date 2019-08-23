package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

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

// GetStage returns a list of stages.
// If page size is >0 it is used as query parameter. Also if nextPageKey is no empty string
// it is used as query parameter.
func (r *StageHandler) GetStage(project string, pageSize int, nextPageKey string) (*models.Stages, error) {
	url, err := url.Parse("http://" + r.BaseURL + "/v1/project/" + project + "/stage")
	if err != nil {
		return nil, err
	}
	q := url.Query()
	if pageSize > 0 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
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
	var stages models.Stages
	err = json.Unmarshal(body, &stages)
	return &stages, err
}
