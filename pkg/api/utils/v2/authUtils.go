package v2

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

// AuthAuthenticateOptions are options for AuthInterface.Authenticate().
type AuthAuthenticateOptions struct{}

type AuthInterface interface {
	// Authenticate authenticates the client request against the server.
	Authenticate(ctx context.Context, opts AuthAuthenticateOptions) (*models.EventContext, *models.Error)
}

// AuthHandler handles projects
type AuthHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewAuthHandler returns a new AuthHandler
func NewAuthHandler(baseURL string) *AuthHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &AuthHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))},
		Scheme:     "http",
	}
}

func createAuthenticatedAuthHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *AuthHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &AuthHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
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
func (a *AuthHandler) Authenticate(ctx context.Context, opts AuthAuthenticateOptions) (*models.EventContext, *models.Error) {
	return postWithEventContext(ctx, a.Scheme+"://"+a.getBaseURL()+"/v1/auth", nil, a)
}
