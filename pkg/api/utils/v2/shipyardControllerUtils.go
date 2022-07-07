package v2

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const shipyardControllerBaseURL = "controlPlane"

// ShipyardControlGetOpenTriggeredEventsOptions are options for ShipyardControlInterface.GetOpenTriggeredEvents().
type ShipyardControlGetOpenTriggeredEventsOptions struct{}

type ShipyardControlInterface interface {
	// GetOpenTriggeredEvents returns all open triggered events.
	GetOpenTriggeredEvents(ctx context.Context, filter EventFilter, opts ShipyardControlGetOpenTriggeredEventsOptions) ([]*models.KeptnContextExtendedCE, error)
}

type ShipyardControllerHandler struct {
	baseURL    string
	authToken  string
	authHeader string
	httpClient *http.Client
	scheme     string
}

// NewShipyardControllerHandler returns a new ShipyardControllerHandler which sends all requests directly to the resource-service
func NewShipyardControllerHandler(baseURL string) *ShipyardControllerHandler {
	return NewShipyardControllerHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewShipyardControllerHandlerWithHTTPClient returns a new ShipyardControllerHandler which sends all requests directly to the resource-service using the specified http.Client
func NewShipyardControllerHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *ShipyardControllerHandler {
	return createShipyardControllerHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedShipyardControllerHandler returns a new ShipyardControllerHandler that authenticates at the api via the provided token
// and sends all requests directly to the resource-service
func NewAuthenticatedShipyardControllerHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ShipyardControllerHandler {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return createShipyardControllerHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createShipyardControllerHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ShipyardControllerHandler {
	return &ShipyardControllerHandler{
		baseURL:    httputils.TrimHTTPScheme(baseURL),
		authHeader: authHeader,
		authToken:  authToken,
		httpClient: httpClient,
		scheme:     scheme,
	}
}

func (s *ShipyardControllerHandler) getBaseURL() string {
	return s.baseURL
}

func (s *ShipyardControllerHandler) getAuthToken() string {
	return s.authToken
}

func (s *ShipyardControllerHandler) getAuthHeader() string {
	return s.authHeader
}

func (s *ShipyardControllerHandler) getHTTPClient() *http.Client {
	return s.httpClient
}

// GetOpenTriggeredEvents returns all open triggered events.
func (s *ShipyardControllerHandler) GetOpenTriggeredEvents(ctx context.Context, filter EventFilter, opts ShipyardControlGetOpenTriggeredEventsOptions) ([]*models.KeptnContextExtendedCE, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	events := []*models.KeptnContextExtendedCE{}
	nextPageKey := ""

	for {
		url, err := url.Parse(s.scheme + "://" + s.getBaseURL() + v1EventPath + "/triggered/" + filter.EventType)

		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		if filter.Project != "" {
			q.Set("project", filter.Project)
		}
		if filter.Service != "" {
			q.Set("service", filter.Service)
		}
		if filter.Stage != "" {
			q.Set("stage", filter.Stage)
		}

		url.RawQuery = q.Encode()

		if err != nil {
			return nil, err
		}

		body, mErr := getAndExpectOK(ctx, url.String(), s)
		if mErr != nil {
			return nil, mErr.ToError()
		}

		received := &models.Events{}
		if err = received.FromJSON(body); err != nil {
			return nil, err
		}
		events = append(events, received.Events...)

		if received.NextPageKey == "" || received.NextPageKey == "0" {
			break
		}

		nextPageKeyInt, _ := strconv.Atoi(received.NextPageKey)

		if filter.NumberOfPages > 0 && nextPageKeyInt >= filter.NumberOfPages {
			break
		}

		nextPageKey = received.NextPageKey
	}
	return events, nil
}
