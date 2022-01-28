package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
)

type EventsV1Interface interface {
	GetEvents(filter *EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error)
	GetEventsWithRetry(filter *EventFilter, maxRetries int, retrySleepTime time.Duration) ([]*models.KeptnContextExtendedCE, error)
}

// EventHandler handles services
type EventHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// EventFilter allows to filter events based on the provided properties
type EventFilter struct {
	Project       string
	Stage         string
	Service       string
	EventType     string
	KeptnContext  string
	EventID       string
	PageSize      string
	NumberOfPages int
	FromTime      string
}

// NewEventHandler returns a new EventHandler
func NewEventHandler(baseURL string) *EventHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &EventHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))},
		Scheme:     "http",
	}
}

const mongodbDatastoreServiceBaseUrl = "mongodb-datastore"

// NewAuthenticatedEventHandler returns a new EventHandler that authenticates at the endpoint via the provided token
// Deprecated: use APISet instead
func NewAuthenticatedEventHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *EventHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedEventHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedEventHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *EventHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, mongodbDatastoreServiceBaseUrl) {
		baseURL += "/" + mongodbDatastoreServiceBaseUrl
	}

	return &EventHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (e *EventHandler) getBaseURL() string {
	return e.BaseURL
}

func (e *EventHandler) getAuthToken() string {
	return e.AuthToken
}

func (e *EventHandler) getAuthHeader() string {
	return e.AuthHeader
}

func (e *EventHandler) getHTTPClient() *http.Client {
	return e.HTTPClient
}

// GetEvents returns all events matching the properties in the passed filter object
func (e *EventHandler) GetEvents(filter *EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error) {

	u, err := url.Parse(e.Scheme + "://" + e.getBaseURL() + "/event?")
	if err != nil {
		log.Fatal("error parsing url")
	}

	query := u.Query()

	if filter.Project != "" {
		query.Set("project", filter.Project)
	}
	if filter.Stage != "" {
		query.Set("stage", filter.Stage)
	}
	if filter.Service != "" {
		query.Set("service", filter.Service)
	}
	if filter.KeptnContext != "" {
		query.Set("keptnContext", filter.KeptnContext)
	}
	if filter.EventID != "" {
		query.Set("eventID", filter.EventID)
	}
	if filter.EventType != "" {
		query.Set("type", filter.EventType)
	}
	if filter.PageSize != "" {
		query.Set("pageSize", filter.PageSize)
	}
	if filter.FromTime != "" {
		query.Set("fromTime", filter.FromTime)
	}

	u.RawQuery = query.Encode()

	return e.getEvents(u.String(), filter.NumberOfPages)
}

// GetEventsWithRetry tries to retrieve events matching the passed filter
func (e *EventHandler) GetEventsWithRetry(filter *EventFilter, maxRetries int, retrySleepTime time.Duration) ([]*models.KeptnContextExtendedCE, error) {
	for i := 0; i < maxRetries; i = i + 1 {
		events, errObj := e.GetEvents(filter)
		if errObj == nil && len(events) > 0 {
			return events, nil
		}
		<-time.After(retrySleepTime)
	}
	return nil, fmt.Errorf("could not find matching event after %d x %s", maxRetries, retrySleepTime.String())
}

func (e *EventHandler) getEvents(uri string, numberOfPages int) ([]*models.KeptnContextExtendedCE, *models.Error) {
	events := []*models.KeptnContextExtendedCE{}
	nextPageKey := ""

	for {
		url, err := url.Parse(uri)
		if err != nil {
			return nil, buildErrorResponse(err.Error())
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			return nil, buildErrorResponse(err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, e)

		resp, err := e.HTTPClient.Do(req)
		if err != nil {
			return nil, buildErrorResponse(err.Error())
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, buildErrorResponse(err.Error())
		}

		if resp.StatusCode == 200 {
			received := &models.Events{}
			if err = received.FromJSON(body); err != nil {
				return nil, buildErrorResponse(err.Error())
			}
			events = append(events, received.Events...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}

			nextPageKeyInt, _ := strconv.Atoi(received.NextPageKey)

			if numberOfPages > 0 && nextPageKeyInt >= numberOfPages {
				break
			}

			nextPageKey = received.NextPageKey
		} else {
			respErr := &models.Error{}
			if err = respErr.FromJSON(body); err != nil {
				return nil, buildErrorResponse(err.Error())
			}
			return nil, respErr
		}
	}

	return events, nil
}
