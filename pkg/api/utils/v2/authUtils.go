package v2

import (
	"context"
	"net/http"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

// AuthAuthenticateOptions are options for AuthInterface.Authenticate().
type AuthAuthenticateOptions struct{}

type AuthInterface interface {
	// Authenticate authenticates the client request against the server.
	Authenticate(ctx context.Context, opts AuthAuthenticateOptions) (*models.EventContext, *models.Error)
}

type AuthHandler struct {
	baseURL    string
	authToken  string
	authHeader string
	httpClient *http.Client
	scheme     string
}

// NewAuthHandler returns a new AuthHandler
func NewAuthHandler(baseURL string) *AuthHandler {
	return NewAuthHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewAuthHandlerWithHTTPClient returns a new AuthHandler using the specified http.Client
func NewAuthHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *AuthHandler {
	return createAuthHandler(baseURL, "", "", httpClient, "http")
}

// NewAuthenticatedAuthHandler returns a new AuthHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedAuthHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *AuthHandler {
	return createAuthHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *AuthHandler {
	return &AuthHandler{
		baseURL:    httputils.TrimHTTPScheme(baseURL),
		authHeader: authHeader,
		authToken:  authToken,
		httpClient: httpClient,
		scheme:     scheme,
	}
}

func (a *AuthHandler) getBaseURL() string {
	return a.baseURL
}

func (a *AuthHandler) getAuthToken() string {
	return a.authToken
}

func (a *AuthHandler) getAuthHeader() string {
	return a.authHeader
}

func (a *AuthHandler) getHTTPClient() *http.Client {
	return a.httpClient
}

// Authenticate authenticates the client request against the server.
func (a *AuthHandler) Authenticate(ctx context.Context, opts AuthAuthenticateOptions) (*models.EventContext, *models.Error) {
	return postWithEventContext(ctx, a.scheme+"://"+a.getBaseURL()+"/v1/auth", nil, a)
}
