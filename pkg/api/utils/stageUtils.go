package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type StagesV1Interface interface {
	// CreateStage creates a new stage with the provided name.
	CreateStage(project string, stageName string) (*models.EventContext, *models.Error)

	// GetAllStages returns a list of all stages.
	GetAllStages(project string) ([]*models.Stage, error)
}

// StageHandler handles stages
type StageHandler struct {
	stageHandler v2.StageHandler
	BaseURL      string
	AuthToken    string
	AuthHeader   string
	HTTPClient   *http.Client
	Scheme       string
}

// NewStageHandler returns a new StageHandler which sends all requests directly to the configuration-service
func NewStageHandler(baseURL string) *StageHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}

	httpClient := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	return &StageHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: httpClient,
		Scheme:     "http",

		stageHandler: v2.StageHandler{
			BaseURL:    baseURL,
			AuthHeader: "",
			AuthToken:  "",
			HTTPClient: httpClient,
			Scheme:     "http",
		},
	}
}

// NewAuthenticatedStageHandler returns a new StageHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
// Deprecated: use APISet instead
func NewAuthenticatedStageHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *StageHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedStageHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedStageHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *StageHandler {
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

		stageHandler: v2.StageHandler{
			BaseURL:    baseURL,
			AuthHeader: authHeader,
			AuthToken:  authToken,
			HTTPClient: httpClient,
			Scheme:     scheme,
		},
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

// CreateStage creates a new stage with the provided name.
func (s *StageHandler) CreateStage(project string, stageName string) (*models.EventContext, *models.Error) {
	return s.stageHandler.CreateStage(context.TODO(), project, stageName, v2.StagesCreateStageOptions{})
}

// GetAllStages returns a list of all stages.
func (s *StageHandler) GetAllStages(project string) ([]*models.Stage, error) {
	return s.stageHandler.GetAllStages(context.TODO(), project, v2.StagesGetAllStagesOptions{})
}
