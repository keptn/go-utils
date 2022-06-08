package v2

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/httputils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// StagesCreateStageOptions are options for StagesInterface.CreateStage().
type StagesCreateStageOptions struct{}

// StagesGetAllStagesOptions are options for StagesInterface.GetAllStages().
type StagesGetAllStagesOptions struct{}

type StagesInterface interface {

	// CreateStage creates a new stage with the provided name.
	CreateStage(ctx context.Context, project string, stageName string, opts StagesCreateStageOptions) (*models.EventContext, *models.Error)

	// GetAllStages returns a list of all stages.
	GetAllStages(ctx context.Context, project string, opts StagesGetAllStagesOptions) ([]*models.Stage, error)
}

// StageHandler handles stages
type StageHandler struct {
	baseURL    string
	authToken  string
	authHeader string
	httpClient *http.Client
	scheme     string
}

// NewStageHandler returns a new StageHandler which sends all requests directly to the configuration-service
func NewStageHandler(baseURL string) *StageHandler {
	return NewStageHandlerWithHTTPClient(baseURL, &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)})
}

// NewStageHandlerWithHTTPClient returns a new StageHandler which sends all requests directly to the configuration-service using the specified http.Client
func NewStageHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *StageHandler {
	return createStageHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedStageHandler returns a new StageHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
func NewAuthenticatedStageHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *StageHandler {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return createStageHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createStageHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *StageHandler {
	return &StageHandler{
		baseURL:    httputils.TrimHTTPScheme(baseURL),
		authHeader: authHeader,
		authToken:  authToken,
		httpClient: httpClient,
		scheme:     scheme,
	}
}

func (s *StageHandler) getBaseURL() string {
	return s.baseURL
}

func (s *StageHandler) getAuthToken() string {
	return s.authToken
}

func (s *StageHandler) getAuthHeader() string {
	return s.authHeader
}

func (s *StageHandler) getHTTPClient() *http.Client {
	return s.httpClient
}

// CreateStage creates a new stage with the provided name.
func (s *StageHandler) CreateStage(ctx context.Context, project string, stageName string, opts StagesCreateStageOptions) (*models.EventContext, *models.Error) {
	stage := models.Stage{StageName: stageName}
	body, err := stage.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, s.scheme+"://"+s.baseURL+v1ProjectPath+"/"+project+pathToStage, body, s)
}

// GetAllStages returns a list of all stages.
func (s *StageHandler) GetAllStages(ctx context.Context, project string, opts StagesGetAllStagesOptions) ([]*models.Stage, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	stages := []*models.Stage{}

	nextPageKey := ""
	for {
		url, err := url.Parse(s.scheme + "://" + s.getBaseURL() + v1ProjectPath + "/" + project + pathToStage)
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}

		body, mErr := getAndExpectOK(ctx, url.String(), s)
		if mErr != nil {
			return nil, mErr.ToError()
		}

		received := &models.Stages{}
		if err = received.FromJSON(body); err != nil {
			return nil, err
		}
		stages = append(stages, received.Stages...)

		if received.NextPageKey == "" || received.NextPageKey == "0" {
			break
		}
		nextPageKey = received.NextPageKey
	}
	return stages, nil
}
