package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
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
	stageHandler *v2.StageHandler
	BaseURL      string
	AuthToken    string
	AuthHeader   string
	HTTPClient   *http.Client
	Scheme       string
}

// NewStageHandler returns a new StageHandler which sends all requests directly to the configuration-service
func NewStageHandler(baseURL string) *StageHandler {
	return NewStageHandlerWithHTTPClient(baseURL, &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)})
}

// NewStageHandlerWithHTTPClient returns a new StageHandler which sends all requests directly to the configuration-service using the specified http.Client
func NewStageHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *StageHandler {
	return &StageHandler{
		BaseURL:      httputils.TrimHTTPScheme(baseURL),
		HTTPClient:   httpClient,
		Scheme:       "http",
		stageHandler: v2.NewStageHandlerWithHTTPClient(baseURL, httpClient),
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
	v2StageHandler := v2.NewAuthenticatedStageHandler(baseURL, authToken, authHeader, httpClient, scheme)

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return &StageHandler{
		BaseURL:      httputils.TrimHTTPScheme(baseURL),
		AuthHeader:   authHeader,
		AuthToken:    authToken,
		HTTPClient:   httpClient,
		Scheme:       scheme,
		stageHandler: v2StageHandler,
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
	s.ensureHandlerIsSet()
	return s.stageHandler.CreateStage(context.TODO(), project, stageName, v2.StagesCreateStageOptions{})
}

// GetAllStages returns a list of all stages.
func (s *StageHandler) GetAllStages(project string) ([]*models.Stage, error) {
	s.ensureHandlerIsSet()
	return s.stageHandler.GetAllStages(context.TODO(), project, v2.StagesGetAllStagesOptions{})
}

func (s *StageHandler) ensureHandlerIsSet() {
	if s.stageHandler != nil {
		return
	}

	if s.AuthToken != "" {
		s.stageHandler = v2.NewAuthenticatedStageHandler(s.BaseURL, s.AuthToken, s.AuthHeader, s.HTTPClient, s.Scheme)
	} else {
		s.stageHandler = v2.NewStageHandlerWithHTTPClient(s.BaseURL, s.HTTPClient)
	}
}
