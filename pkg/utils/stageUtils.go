package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/keptn/go-utils/pkg/models"
)

// StageHandler handles stages
type StageHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
}

// NewStageHandler returns a new StageHandler
func NewStageHandler(baseURL string) *StageHandler {
	return &StageHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
	}
}

// NewAuthenticatedStageHandler returns a new StageHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedStageHandler(baseURL string, authToken string, authHeader string) *StageHandler {
	return &StageHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
	}
}

// CreateStage creates a new stage with the provided name
func (r *StageHandler) CreateStage(project string, stageName string) (*models.Error, error) {

	stage := models.Stage{StageName: stageName}
	resourceStr, err := json.Marshal(stage)
	if err != nil {
		return nil, err
	}
	return r.post("http://"+r.BaseURL+"/v1/project/"+project+"/stage", resourceStr)
}

func (r *StageHandler) post(uri string, data []byte) (*models.Error, error) {

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if r.AuthHeader != "" && r.AuthToken != "" {
		req.Header.Set(r.AuthHeader, r.AuthToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respErr models.Error
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return nil, err
	}

	return &respErr, nil
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
