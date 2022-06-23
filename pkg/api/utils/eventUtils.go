package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

type EventsV1Interface interface {
	// GetEvents returns all events matching the properties in the passed filter object.
	GetEvents(filter *EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error)

	// GetEventsWithRetry tries to retrieve events matching the passed filter.
	GetEventsWithRetry(filter *EventFilter, maxRetries int, retrySleepTime time.Duration) ([]*models.KeptnContextExtendedCE, error)
}

// EventHandler handles services
type EventHandler struct {
	eventHandler *v2.EventHandler
	BaseURL      string
	AuthToken    string
	AuthHeader   string
	HTTPClient   *http.Client
	Scheme       string
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
	return NewEventHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewEventHandlerWithHTTPClient returns a new EventHandler that uses the specified http.Client
func NewEventHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *EventHandler {
	return &EventHandler{
		BaseURL:      httputils.TrimHTTPScheme(baseURL),
		HTTPClient:   httpClient,
		Scheme:       "http",
		eventHandler: v2.NewEventHandlerWithHTTPClient(baseURL, httpClient),
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
	v2EventHandler := v2.NewAuthenticatedEventHandler(baseURL, authToken, authHeader, httpClient, scheme)

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, mongodbDatastoreServiceBaseUrl) {
		baseURL += "/" + mongodbDatastoreServiceBaseUrl
	}

	return &EventHandler{
		BaseURL:      httputils.TrimHTTPScheme(baseURL),
		AuthHeader:   authHeader,
		AuthToken:    authToken,
		HTTPClient:   httpClient,
		Scheme:       scheme,
		eventHandler: v2EventHandler,
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

// GetEvents returns all events matching the properties in the passed filter object.
func (e *EventHandler) GetEvents(filter *EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error) {
	e.ensureHandlerIsSet()
	return e.eventHandler.GetEvents(context.TODO(), toV2EventFilter(filter), v2.EventsGetEventsOptions{})
}

// GetEventsWithRetry tries to retrieve events matching the passed filter.
func (e *EventHandler) GetEventsWithRetry(filter *EventFilter, maxRetries int, retrySleepTime time.Duration) ([]*models.KeptnContextExtendedCE, error) {
	e.ensureHandlerIsSet()
	return e.eventHandler.GetEventsWithRetry(context.TODO(), toV2EventFilter(filter), maxRetries, retrySleepTime, v2.EventsGetEventsWithRetryOptions{})
}

func toV2EventFilter(filter *EventFilter) *v2.EventFilter {
	return &v2.EventFilter{
		Project:       filter.Project,
		Stage:         filter.Stage,
		Service:       filter.Service,
		EventType:     filter.EventType,
		KeptnContext:  filter.KeptnContext,
		EventID:       filter.EventID,
		PageSize:      filter.PageSize,
		NumberOfPages: filter.NumberOfPages,
		FromTime:      filter.FromTime,
	}
}

func (e *EventHandler) ensureHandlerIsSet() {
	if e.eventHandler != nil {
		return
	}

	if e.AuthToken != "" {
		e.eventHandler = v2.NewAuthenticatedEventHandler(e.BaseURL, e.AuthToken, e.AuthHeader, e.HTTPClient, e.Scheme)
	} else {
		e.eventHandler = v2.NewEventHandlerWithHTTPClient(e.BaseURL, e.HTTPClient)
	}
}
