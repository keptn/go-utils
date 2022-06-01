package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
)

const shipyardControllerBaseURL = "controlPlane"

type ShipyardControlV1Interface interface {
	// GetOpenTriggeredEvents returns all open triggered events.
	GetOpenTriggeredEvents(filter EventFilter) ([]*models.KeptnContextExtendedCE, error)
}

// ShipyardControllerHandler handles services
type ShipyardControllerHandler struct {
	shipyardControllerHandler v2.ShipyardControllerHandler
	BaseURL                   string
	AuthToken                 string
	AuthHeader                string
	HTTPClient                *http.Client
	Scheme                    string
}

// NewShipyardControllerHandler returns a new ShipyardControllerHandler which sends all requests directly to the configuration-service
func NewShipyardControllerHandler(baseURL string) *ShipyardControllerHandler {
	return NewShipyardControllerHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewShipyardControllerHandlerWithHTTPClient returns a new ShipyardControllerHandler which sends all requests directly to the configuration-service using the specified http.Client
func NewShipyardControllerHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *ShipyardControllerHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}

	return createShipyardControllerHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedShipyardControllerHandler returns a new ShipyardControllerHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
// Deprecated: use APISet instead
func NewAuthenticatedShipyardControllerHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ShipyardControllerHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedShipyardControllerHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedShipyardControllerHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ShipyardControllerHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return createShipyardControllerHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createShipyardControllerHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ShipyardControllerHandler {
	return &ShipyardControllerHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,

		shipyardControllerHandler: v2.ShipyardControllerHandler{
			BaseURL:    baseURL,
			AuthHeader: authHeader,
			AuthToken:  authToken,
			HTTPClient: httpClient,
			Scheme:     scheme,
		},
	}
}

func (s *ShipyardControllerHandler) getBaseURL() string {
	return s.BaseURL
}

func (s *ShipyardControllerHandler) getAuthToken() string {
	return s.AuthToken
}

func (s *ShipyardControllerHandler) getAuthHeader() string {
	return s.AuthHeader
}

func (s *ShipyardControllerHandler) getHTTPClient() *http.Client {
	return s.HTTPClient
}

// GetOpenTriggeredEvents returns all open triggered events.
func (s *ShipyardControllerHandler) GetOpenTriggeredEvents(filter EventFilter) ([]*models.KeptnContextExtendedCE, error) {
	return s.shipyardControllerHandler.GetOpenTriggeredEvents(context.TODO(), *toV2EventFilter(&filter), v2.ShipyardControlGetOpenTriggeredEventsOptions{})
}
