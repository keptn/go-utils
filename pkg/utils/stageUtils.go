package utils

import (
	"encoding/json"
	"errors"
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

func (s *StageHandler) getBaseURL() string {
	return s.BaseURL
}

func (s *StageHandler) getAuthToken() string {
	return s.AuthToken
}

func (s *StageHandler) getAuthHeader() string {
	return s.AuthHeader
}

// CreateStage creates a new stage with the provided name
func (s *StageHandler) CreateStage(project string, stageName string) (*models.Error, error) {

	stage := models.Stage{StageName: stageName}
	body, err := json.Marshal(stage)
	if err != nil {
		return nil, err
	}
	return post("http://"+s.BaseURL+"/v1/project/"+project+"/stage", body, s)
}

// GetAllStages returns a list of all stages.
func (s *StageHandler) GetAllStages(project string) ([]*models.Stage, error) {

	stages := []*models.Stage{}

	nextPageKey := ""
	for {
		url, err := url.Parse("http://" + s.getBaseURL() + "/v1/project/" + project + "/stage")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, s)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
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
		} else {
			var respErr models.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, err
			}
			return nil, errors.New("Response Error Code: " + string(respErr.Code) + " Message: " + *respErr.Message)
		}
	}
	return stages, nil
}
