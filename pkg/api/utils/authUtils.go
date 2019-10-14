package utils

import (
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

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
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &AuthHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{},
		Scheme:     "http",
	}
}

// NewAuthenticatedAuthHandler returns a new AuthHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedAuthHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *AuthHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
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

func (p *AuthHandler) getBaseURL() string {
	return p.BaseURL
}

func (p *AuthHandler) getAuthToken() string {
	return p.AuthToken
}

func (p *AuthHandler) getAuthHeader() string {
	return p.AuthHeader
}

func (p *AuthHandler) getHTTPClient() *http.Client {
	return p.HTTPClient
}

// Authenticate creates a new project
func (p *AuthHandler) Authenticate() (*models.ChannelInfo, *models.Error) {
	return post(p.Scheme+"://"+p.getBaseURL()+"/v1/auth", nil, p)
}
