package v2

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const v1LogPath = "/v1/log"

var defaultSyncInterval = 1 * time.Minute

// LogsLogOptions are options for LogsInterface.Log().
type LogsLogOptions struct{}

// LogsFlushOptions are options for LogsInterface.Flush().
type LogsFlushOptions struct{}

// LogsGetLogsOptions are options for LogsInterface.GetLogs().
type LogsGetLogsOptions struct{}

// LogsDeleteLogsOptions are options for LogsInterface.DeleteLogs().
type LogsDeleteLogsOptions struct{}

// LogsStartOptions are options for LogsInterface.Start().
type LogsStartOptions struct{}

//go:generate moq -pkg utils_mock -skip-ensure -out ./fake/log_handler_mock.go . LogsInterface
type LogsInterface interface {
	// Log appends the specified logs to the log cache.
	Log(logs []models.LogEntry, opts LogsLogOptions)

	// Flush flushes the log cache.
	Flush(ctx context.Context, opts LogsFlushOptions) error

	// GetLogs gets logs with the specified parameters.
	GetLogs(ctx context.Context, params models.GetLogsParams, opts LogsGetLogsOptions) (*models.GetLogsResponse, error)

	// DeleteLogs deletes logs matching the specified log filter.
	DeleteLogs(ctx context.Context, filter models.LogFilter, opts LogsDeleteLogsOptions) error

	Start(ctx context.Context, opts LogsStartOptions)
}

type LogHandler struct {
	baseURL      string
	authToken    string
	authHeader   string
	httpClient   *http.Client
	scheme       string
	logCache     []models.LogEntry
	theClock     clock.Clock
	syncInterval time.Duration
	lock         sync.Mutex
}

// NewLogHandler returns a new LogHandler
func NewLogHandler(baseURL string) *LogHandler {
	return NewLogHandlerWithHTTPClient(baseURL, &http.Client{Transport: getClientTransport(nil)})
}

// NewLogHandlerWithHTTPClient returns a new LogHandler that uses the specified http.Client
func NewLogHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *LogHandler {
	return createLogHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedLogHandler returns a new LogHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedLogHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *LogHandler {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return createLogHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createLogHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *LogHandler {
	return &LogHandler{
		baseURL:      httputils.TrimHTTPScheme(baseURL),
		authHeader:   authHeader,
		authToken:    authToken,
		httpClient:   httpClient,
		scheme:       scheme,
		logCache:     []models.LogEntry{},
		theClock:     clock.New(),
		syncInterval: defaultSyncInterval,
	}
}

func (lh *LogHandler) getBaseURL() string {
	return lh.baseURL
}

func (lh *LogHandler) getAuthToken() string {
	return lh.authToken
}

func (lh *LogHandler) getAuthHeader() string {
	return lh.authHeader
}

func (lh *LogHandler) getHTTPClient() *http.Client {
	return lh.httpClient
}

// Log appends the specified logs to the log cache.
func (lh *LogHandler) Log(logs []models.LogEntry, opts LogsLogOptions) {
	lh.lock.Lock()
	defer lh.lock.Unlock()
	lh.logCache = append(lh.logCache, logs...)
}

// GetLogs gets logs with the specified parameters.
func (lh *LogHandler) GetLogs(ctx context.Context, params models.GetLogsParams, opts LogsGetLogsOptions) (*models.GetLogsResponse, error) {
	u, err := url.Parse(lh.scheme + "://" + lh.getBaseURL() + v1LogPath)
	if err != nil {
		log.Fatal("error parsing url")
	}

	query := u.Query()

	if params.IntegrationID != "" {
		query.Set("integrationId", params.IntegrationID)
	}
	if params.PageSize != 0 {
		query.Set("pageSize", fmt.Sprintf("%d", params.PageSize))
	}
	if params.FromTime != "" {
		query.Set("fromTime", params.FromTime)
	}
	if params.BeforeTime != "" {
		query.Set("beforeTime", params.BeforeTime)
	}

	u.RawQuery = query.Encode()

	body, mErr := getAndExpectOK(ctx, u.String(), lh)
	if mErr != nil {
		return nil, mErr.ToError()
	}

	received := &models.GetLogsResponse{}
	if err := received.FromJSON(body); err != nil {
		return nil, err
	}

	return received, nil
}

// DeleteLogs deletes logs matching the specified log filter.
func (lh *LogHandler) DeleteLogs(ctx context.Context, params models.LogFilter, opts LogsDeleteLogsOptions) error {
	u, err := url.Parse(lh.scheme + "://" + lh.getBaseURL() + v1LogPath)
	if err != nil {
		log.Fatal("error parsing url")
	}

	query := u.Query()

	if params.IntegrationID != "" {
		query.Set("integrationId", params.IntegrationID)
	}
	if params.FromTime != "" {
		query.Set("fromTime", params.FromTime)
	}
	if params.BeforeTime != "" {
		query.Set("beforeTime", params.BeforeTime)
	}
	if _, err := delete(ctx, u.String(), lh); err != nil {
		return errors.New(err.GetMessage())
	}
	return nil
}

func (lh *LogHandler) Start(ctx context.Context, opts LogsStartOptions) {
	ticker := lh.theClock.Ticker(lh.syncInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lh.Flush(ctx, LogsFlushOptions{})
			}
		}
	}()
}

// Flush flushes the log cache.
func (lh *LogHandler) Flush(ctx context.Context, opts LogsFlushOptions) error {
	lh.lock.Lock()
	defer lh.lock.Unlock()
	if len(lh.logCache) == 0 {
		// only send a request if we actually have some logs to send
		return nil
	}
	createLogsPayload := &models.CreateLogsRequest{
		Logs: lh.logCache,
	}
	bodyStr, err := createLogsPayload.ToJSON()
	if err != nil {
		return err
	}
	if _, err := post(ctx, lh.scheme+"://"+lh.getBaseURL()+v1LogPath, bodyStr, lh); err != nil {
		return errors.New(err.GetMessage())
	}
	lh.logCache = []models.LogEntry{}
	return nil
}
