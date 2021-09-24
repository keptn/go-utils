package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

// StageHandler handles stages
type StageHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewStageHandler returns a new StageHandler which sends all requests directly to the configuration-service
func NewStageHandler(baseURL string) *StageHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &StageHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{},
		Scheme:     "http",
	}
}

// NewAuthenticatedStageHandler returns a new StageHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
func NewAuthenticatedStageHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *StageHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}
	return &StageHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
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

func (s *StageHandler) getHTTPClient() *http.Client {
	return s.HTTPClient
}

// CreateStage creates a new stage with the provided name
//
// Deprecated: Use CreateStageWithContext instead
func (s *StageHandler) CreateStage(project string, stageName string) (*models.EventContext, *models.Error) {
	return s.CreateStageWithContext(context.Background(), project, stageName)
}

// CreateStageWithContext creates a new stage with the provided name
func (s *StageHandler) CreateStageWithContext(ctx context.Context, project string, stageName string) (*models.EventContext, *models.Error) {
	stage := models.Stage{StageName: stageName}
	body, err := json.Marshal(stage)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, s.Scheme+"://"+s.BaseURL+"/v1/project/"+project+"/stage", body, s)
}

// GetAllStages returns a list of all stages.
//
// Deprecated: Use GetAllStagesWithContext instead
func (s *StageHandler) GetAllStages(project string) ([]*models.Stage, error) {
	return s.GetAllStagesWithContext(context.Background(), project)
}

// GetAllStagesWithContext returns a list of all stages.
func (s *StageHandler) GetAllStagesWithContext(ctx context.Context, project string) ([]*models.Stage, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	stages := []*models.Stage{}

	nextPageKey := ""
	for {
		url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + "/v1/project/" + project + "/stage")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, s)

		resp, err := s.HTTPClient.Do(req)
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

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey
		} else {
			var respErr models.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(*respErr.Message)
		}
	}
	return stages, nil
}
