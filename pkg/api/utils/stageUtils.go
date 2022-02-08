package api

import (
	"crypto/tls"
	"errors"
	"github.com/keptn/go-utils/pkg/common/httputils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type StagesV1Interface interface {
	CreateStage(project string, stageName string) (*models.EventContext, *models.Error)
	GetAllStages(project string) ([]*models.Stage, error)
}

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
	return createStageHandler(baseURL)
}

// NewAuthenticatedStageHandler returns a new StageHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
// Deprecated: use APISet instead
func NewAuthenticatedStageHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *StageHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedStageHandler(baseURL, authToken, authHeader, httpClient, scheme, false)
}

func createAuthenticatedStageHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string, internal bool) *StageHandler {
	baseURL = httputils.TrimHTTPScheme(baseURL)
	if internal {
		return &StageHandler{
			BaseURL:    baseURL,
			AuthHeader: "",
			AuthToken:  "",
			HTTPClient: &http.Client{Transport: otelhttp.NewTransport(httpClient.Transport)},
			Scheme:     "http",
		}
	}
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

func createStageHandler(baseURL string) *StageHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &StageHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
		Scheme:     "http",
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
func (s *StageHandler) CreateStage(project string, stageName string) (*models.EventContext, *models.Error) {

	stage := models.Stage{StageName: stageName}
	body, err := stage.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(s.Scheme+"://"+s.BaseURL+pathToProject+project+pathToStage, body, s)
}

// GetAllStages returns a list of all stages.
func (s *StageHandler) GetAllStages(project string) ([]*models.Stage, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	stages := []*models.Stage{}

	nextPageKey := ""
	for {
		url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + pathToProject + project + pathToStage)
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, s)

		resp, err := s.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			received := &models.Stages{}
			if err = received.FromJSON(body); err != nil {
				return nil, err
			}
			stages = append(stages, received.Stages...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey
		} else {
			respErr := &models.Error{}
			if err = respErr.FromJSON(body); err != nil {
				return nil, err
			}
			return nil, errors.New(*respErr.Message)
		}
	}
	return stages, nil
}
