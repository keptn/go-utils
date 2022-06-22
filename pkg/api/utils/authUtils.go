package api

import (
	"context"
	"net/http"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

type AuthV1Interface interface {
	// Authenticate authenticates the client request against the server.
	Authenticate() (*models.EventContext, *models.Error)
}

// AuthHandler handles projects
type AuthHandler struct {
	authHandler *v2.AuthHandler
	BaseURL     string
	AuthToken   string
	AuthHeader  string
	HTTPClient  *http.Client
	Scheme      string
}

// NewAuthHandler returns a new AuthHandler
func NewAuthHandler(baseURL string) *AuthHandler {
	return NewAuthHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewAuthHandlerWithHTTPClient returns a new AuthHandler that uses the specified http.Client
func NewAuthHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *AuthHandler {
	return &AuthHandler{
		BaseURL:     httputils.TrimHTTPScheme(baseURL),
		HTTPClient:  httpClient,
		Scheme:      "http",
		authHandler: v2.NewAuthHandlerWithHTTPClient(baseURL, httpClient),
	}
}

// NewAuthenticatedAuthHandler returns a new AuthHandler that authenticates at the endpoint via the provided token
// Deprecated: use APISet instead
func NewAuthenticatedAuthHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *AuthHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))

	return createAuthenticatedAuthHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedAuthHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *AuthHandler {
	return &AuthHandler{
		BaseURL:     httputils.TrimHTTPScheme(baseURL),
		AuthHeader:  authHeader,
		AuthToken:   authToken,
		HTTPClient:  httpClient,
		Scheme:      scheme,
		authHandler: v2.NewAuthenticatedAuthHandler(baseURL, authToken, authHeader, httpClient, scheme),
	}
}

func (a *AuthHandler) getBaseURL() string {
	return a.BaseURL
}

func (a *AuthHandler) getAuthToken() string {
	return a.AuthToken
}

func (a *AuthHandler) getAuthHeader() string {
	return a.AuthHeader
}

func (a *AuthHandler) getHTTPClient() *http.Client {
	return a.HTTPClient
}

// Authenticate authenticates the client request against the server.
func (a *AuthHandler) Authenticate() (*models.EventContext, *models.Error) {
	a.ensureHandlerIsSet()
	return a.authHandler.Authenticate(context.TODO(), v2.AuthAuthenticateOptions{})
}

func (a *AuthHandler) ensureHandlerIsSet() {
	if a.authHandler != nil {
		return
	}

	if a.AuthToken != "" {
		a.authHandler = v2.NewAuthenticatedAuthHandler(a.BaseURL, a.AuthToken, a.AuthHeader, a.HTTPClient, a.Scheme)
	} else {
		a.authHandler = v2.NewAuthHandlerWithHTTPClient(a.BaseURL, a.HTTPClient)
	}
}
