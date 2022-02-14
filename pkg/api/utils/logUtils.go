package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/httputils"
	"io/ioutil"
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
	Log(logs []models.LogEntry)
	Flush() error
	GetLogs(params models.GetLogsParams) (*models.GetLogsResponse, error)
	DeleteLogs(filter models.LogFilter) error
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
	return createLogHandler(baseURL)
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
	baseURL = httputils.TrimHTTPScheme(baseURL)

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

func createLogHandler(baseURL string) *LogHandler {
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

func (lh *LogHandler) Log(logs []models.LogEntry) {
	lh.lock.Lock()
	defer lh.lock.Unlock()
	lh.LogCache = append(lh.LogCache, logs...)
}

func (lh *LogHandler) GetLogs(params models.GetLogsParams) (*models.GetLogsResponse, error) {
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

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, lh)

	resp, err := lh.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		received := &models.GetLogsResponse{}
		if err := received.FromJSON(body); err != nil {
			return nil, err
		}
		return received, nil
	}
	errResponse := &models.Error{}
	if err := errResponse.FromJSON(body); err != nil {
		return nil, err
	}
	return nil, errors.New(errResponse.GetMessage())

}

func (lh *LogHandler) DeleteLogs(params models.LogFilter) error {
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
	if _, err := delete(u.String(), lh); err != nil {
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

func (lh *LogHandler) Flush() error {
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
	if _, err := post(lh.Scheme+"://"+lh.getBaseURL()+v1LogPath, bodyStr, lh); err != nil {
		return errors.New(err.GetMessage())
	}
	lh.LogCache = []models.LogEntry{}
	return nil
}
