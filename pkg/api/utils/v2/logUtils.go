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

	// FlushWithContext flushes the log cache.
	FlushWithContext(ctx context.Context) error

	// GetLogs gets logs with the specified parameters.
	GetLogs(params models.GetLogsParams) (*models.GetLogsResponse, error)

	// GetLogsWithContext gets logs with the specified parameters.
	GetLogsWithContext(ctx context.Context, params models.GetLogsParams) (*models.GetLogsResponse, error)

	// DeleteLogs deletes logs matching the specified log filter.
	DeleteLogs(filter models.LogFilter) error

	// DeleteLogsWithContext deletes logs matching the specified log filter.
	DeleteLogsWithContext(ctx context.Context, filter models.LogFilter) error

	Start(ctx context.Context)
}

type LogHandler struct {
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

func NewLogHandler(baseURL string) *LogHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &LogHandler{
		BaseURL:      baseURL,
		AuthHeader:   "",
		AuthToken:    "",
		HTTPClient:   &http.Client{Transport: getClientTransport(nil)},
		Scheme:       "http",
		LogCache:     []models.LogEntry{},
		TheClock:     clock.New(),
		SyncInterval: defaultSyncInterval,
	}
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

	return &LogHandler{
		BaseURL:      baseURL,
		AuthHeader:   authHeader,
		AuthToken:    authToken,
		HTTPClient:   httpClient,
		Scheme:       scheme,
		LogCache:     []models.LogEntry{},
		TheClock:     clock.New(),
		SyncInterval: defaultSyncInterval,
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
	lh.lock.Lock()
	defer lh.lock.Unlock()
	lh.LogCache = append(lh.LogCache, logs...)
}

// GetLogs gets logs with the specified parameters.
func (lh *LogHandler) GetLogs(params models.GetLogsParams) (*models.GetLogsResponse, error) {
	return lh.GetLogsWithContext(context.TODO(), params)
}

// GetLogsWithContext gets logs with the specified parameters.
func (lh *LogHandler) GetLogsWithContext(ctx context.Context, params models.GetLogsParams) (*models.GetLogsResponse, error) {
	u, err := url.Parse(lh.Scheme + "://" + lh.getBaseURL() + v1LogPath)
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
func (lh *LogHandler) DeleteLogs(params models.LogFilter) error {
	return lh.DeleteLogsWithContext(context.TODO(), params)
}

// DeleteLogsWithContext deletes logs matching the specified log filter.
func (lh *LogHandler) DeleteLogsWithContext(ctx context.Context, params models.LogFilter) error {
	u, err := url.Parse(lh.Scheme + "://" + lh.getBaseURL() + v1LogPath)
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

func (lh *LogHandler) Start(ctx context.Context) {
	ticker := lh.TheClock.Ticker(lh.SyncInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lh.Flush()
			}
		}
	}()
}

// Flush flushes the log cache.
func (lh *LogHandler) Flush() error {
	return lh.FlushWithContext(context.TODO())
}

// FlushWithContext flushes the log cache.
func (lh *LogHandler) FlushWithContext(ctx context.Context) error {
	lh.lock.Lock()
	defer lh.lock.Unlock()
	if len(lh.LogCache) == 0 {
		// only send a request if we actually have some logs to send
		return nil
	}
	createLogsPayload := &models.CreateLogsRequest{
		Logs: lh.LogCache,
	}
	bodyStr, err := createLogsPayload.ToJSON()
	if err != nil {
		return err
	}
	if _, err := post(ctx, lh.Scheme+"://"+lh.getBaseURL()+v1LogPath, bodyStr, lh); err != nil {
		return errors.New(err.GetMessage())
	}
	lh.LogCache = []models.LogEntry{}
	return nil
}
