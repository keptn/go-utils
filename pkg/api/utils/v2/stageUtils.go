package v2

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
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
		HTTPClient: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
		Scheme:     "http",
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
func (s *StageHandler) CreateStage(ctx context.Context, project string, stageName string, opts StagesCreateStageOptions) (*models.EventContext, *models.Error) {
	stage := models.Stage{StageName: stageName}
	body, err := stage.ToJSON()
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(ctx, s.Scheme+"://"+s.BaseURL+v1ProjectPath+"/"+project+pathToStage, body, s)
}

// GetAllStages returns a list of all stages.
func (s *StageHandler) GetAllStages(ctx context.Context, project string, opts StagesGetAllStagesOptions) ([]*models.Stage, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	stages := []*models.Stage{}

	nextPageKey := ""
	for {
		url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + v1ProjectPath + "/" + project + pathToStage)
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