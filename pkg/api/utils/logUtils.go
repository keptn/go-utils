package api

import (
	"encoding/json"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const v1LogPath = "/v1/log"

type GetLogsParams struct {
	LogFilter
	PageSize    int
	NextPageKey int
}

type LogFilter struct {
	IntegrationID string
	FromTime      string
	BeforeTime    string
}

type ILogHandler interface {
	Log(logEntry models.LogEntry) (string, *models.Error)
	GetLogs(params GetLogsParams) ([]models.LogEntry, *models.Error)
	DeleteLogs(filter LogFilter) (string, *models.Error)
}

type LogHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

func NewLogHandler(baseURL string) *LogHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &LogHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     "http",
	}
}

func NewAuthenticatedLogHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *LogHandler {
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

	return &LogHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
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

func (lh *LogHandler) Log(logEntry models.LogEntry) (string, *models.Error) {
	bodyStr, err := json.Marshal(logEntry)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return post(lh.Scheme+"://"+lh.getBaseURL()+v1LogPath, bodyStr, lh)
}

func (lh *LogHandler) GetLogs(params GetLogsParams) ([]models.LogEntry, *models.Error) {
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
		query.Set("beforeTim", params.BeforeTime)
	}

	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, lh)

	resp, err := lh.HTTPClient.Do(req)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var received []models.LogEntry
		err := json.Unmarshal(body, &received)
		if err != nil {
			return nil, buildErrorResponse(err.Error())
		}
		return received, nil
	}

	return nil, nil
}

func (lh *LogHandler) DeleteLogs(params LogFilter) (string, *models.Error) {
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
		query.Set("beforeTim", params.BeforeTime)
	}
	return delete(u.String(), lh)
}
