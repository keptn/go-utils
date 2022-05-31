package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const uniformRegistrationBaseURL = "uniform/registration"
const v1UniformPath = "/v1/uniform/registration"

type UniformV1Interface interface {
	Ping(integrationID string) (*models.Integration, error)
	RegisterIntegration(integration models.Integration) (string, error)
	CreateSubscription(integrationID string, subscription models.EventSubscription) (string, error)
	UnregisterIntegration(integrationID string) error
	GetRegistrations() ([]*models.Integration, error)
}

type UniformHandler struct {
	uniformHandler v2.UniformHandler
	BaseURL        string
	AuthToken      string
	AuthHeader     string
	HTTPClient     *http.Client
	Scheme         string
}

func NewUniformHandler(baseURL string) *UniformHandler {
	baseURL = httputils.TrimHTTPScheme(baseURL)

	httpClient := &http.Client{Transport: getClientTransport(nil)}

	return &UniformHandler{
		BaseURL:    baseURL,
		AuthToken:  "",
		AuthHeader: "",
		HTTPClient: httpClient,
		Scheme:     "http",

		uniformHandler: v2.UniformHandler{
			BaseURL:    baseURL,
			AuthToken:  "",
			AuthHeader: "",
			HTTPClient: httpClient,
			Scheme:     "http",
		},
	}
}

// NewAuthenticatedUniformHandler returns a new UniformHandler that authenticates at the api via the provided token
// Deprecated: use APISet instead
func NewAuthenticatedUniformHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *UniformHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport(httpClient.Transport)
	return createAuthenticatedUniformHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedUniformHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *UniformHandler {
	baseURL = httputils.TrimHTTPScheme(baseURL)
	baseURL = strings.TrimRight(baseURL, "/")

	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return &UniformHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,

		uniformHandler: v2.UniformHandler{
			BaseURL:    baseURL,
			AuthHeader: authHeader,
			AuthToken:  authToken,
			HTTPClient: httpClient,
			Scheme:     scheme,
		},
	}
}

func (u *UniformHandler) getBaseURL() string {
	return u.BaseURL
}

func (u *UniformHandler) getAuthToken() string {
	return u.AuthToken
}

func (u *UniformHandler) getAuthHeader() string {
	return u.AuthHeader
}

func (u *UniformHandler) getHTTPClient() *http.Client {
	return u.HTTPClient
}

func (u *UniformHandler) Ping(integrationID string) (*models.Integration, error) {
	return u.uniformHandler.Ping(context.TODO(), integrationID, v2.UniformPingOptions{})
}

func (u *UniformHandler) RegisterIntegration(integration models.Integration) (string, error) {
	return u.uniformHandler.RegisterIntegration(context.TODO(), integration, v2.UniformRegisterIntegrationOptions{})
}

func (u *UniformHandler) CreateSubscription(integrationID string, subscription models.EventSubscription) (string, error) {
	return u.uniformHandler.CreateSubscription(context.TODO(), integrationID, subscription, v2.UniformCreateSubscriptionOptions{})
}

func (u *UniformHandler) UnregisterIntegration(integrationID string) error {
	return u.uniformHandler.UnregisterIntegration(context.TODO(), integrationID, v2.UniformUnregisterIntegrationOptions{})
}

func (u *UniformHandler) GetRegistrations() ([]*models.Integration, error) {
	return u.uniformHandler.GetRegistrations(context.TODO(), v2.UniformGetRegistrationsOptions{})
}
