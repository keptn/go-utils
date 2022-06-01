package api

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
)

const v1LogPath = "/v1/log"

var defaultSyncInterval = 1 * time.Minute

type LogsV1Interface interface {
	ILogHandler
}

//go:generate moq -pkg utils_mock -skip-ensure -out ./fake/log_handler_mock.go . ILogHandler
type ILogHandler interface {

	// Log appends the specified logs to the log cache.
	Log(logs []models.LogEntry)

	// Flush flushes the log cache.
	Flush() error

	// GetLogs gets logs with the specified parameters.
	GetLogs(params models.GetLogsParams) (*models.GetLogsResponse, error)

	// DeleteLogs deletes logs matching the specified log filter.
	DeleteLogs(filter models.LogFilter) error

	Start(ctx context.Context)
}

type LogHandler struct {
	logHandler   v2.LogHandler
	BaseURL      string
	AuthToken    string
	AuthHeader   string
	HTTPClient   *http.Client
	Scheme       string
	LogCache     []models.LogEntry
	TheClock     clock.Clock
	SyncInterval time.Duration
	lock         sync.Mutex
}

// NewLogHandler returns a new LogHandler
func NewLogHandler(baseURL string) *LogHandler {
	return NewLogHandlerWithHTTPClient(baseURL, &http.Client{Transport: getClientTransport(nil)})
}

// NewLogHandlerWithHTTPClient returns a new LogHandler that uses the specified http.Client
func NewLogHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *LogHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}

	return createLogHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedLogHandler returns a new EventHandler that authenticates at the endpoint via the provided token
// Deprecated: use APISet instead
func NewAuthenticatedLogHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *LogHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport(httpClient.Transport)
	return createAuthenticatedLogHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedLogHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *LogHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return createLogHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createLogHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *LogHandler {
	return &LogHandler{
		BaseURL:      baseURL,
		AuthHeader:   authHeader,
		AuthToken:    authToken,
		HTTPClient:   httpClient,
		Scheme:       scheme,
		LogCache:     []models.LogEntry{},
		TheClock:     clock.New(),
		SyncInterval: defaultSyncInterval,

		logHandler: v2.LogHandler{
			BaseURL:      baseURL,
			AuthHeader:   authHeader,
			AuthToken:    authToken,
			HTTPClient:   httpClient,
			Scheme:       scheme,
			LogCache:     []models.LogEntry{},
			TheClock:     clock.New(),
			SyncInterval: defaultSyncInterval,
		},
	}
}

func (lh *LogHandler) getBaseURL() string {
	return lh.BaseURL
}

func (lh *LogHandler) getAuthToken() string {
	return lh.AuthToken
}

func (lh *LogHandler) getAuthHeader() string {
	return lh.AuthHeader
}

func (lh *LogHandler) getHTTPClient() *http.Client {
	return lh.HTTPClient
}

// Log appends the specified logs to the log cache.
func (lh *LogHandler) Log(logs []models.LogEntry) {
	lh.logHandler.Log(logs, v2.LogsLogOptions{})
}

// GetLogs gets logs with the specified parameters.
func (lh *LogHandler) GetLogs(params models.GetLogsParams) (*models.GetLogsResponse, error) {
	return lh.logHandler.GetLogs(context.TODO(), params, v2.LogsGetLogsOptions{})
}

// DeleteLogs deletes logs matching the specified log filter.
func (lh *LogHandler) DeleteLogs(params models.LogFilter) error {
	return lh.logHandler.DeleteLogs(context.TODO(), params, v2.LogsDeleteLogsOptions{})
}

func (lh *LogHandler) Start(ctx context.Context) {
	lh.logHandler.Start(ctx, v2.LogsStartOptions{})
}

// Flush flushes the log cache.
func (lh *LogHandler) Flush() error {
	return lh.logHandler.Flush(context.TODO(), v2.LogsFlushOptions{})
}
