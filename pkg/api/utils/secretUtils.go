package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const secretServiceBaseURL = "secrets"
const v1SecretPath = "/v1/secret"

type SecretsV1Interface interface {
	SecretHandlerInterface
}

//go:generate moq -pkg utils_mock -skip-ensure -out ./fake/secret_handler_mock.go . SecretHandlerInterface
type SecretHandlerInterface interface {
	// CreateSecret creates a new secret.
	CreateSecret(secret models.Secret) error

	// UpdateSecret creates a new secret.
	UpdateSecret(secret models.Secret) error

	// DeleteSecret deletes a secret.
	DeleteSecret(secretName, secretScope string) error

	// GetSecrets returns a list of created secrets.
	GetSecrets() (*models.GetSecretsResponse, error)
}

// SecretHandler handles services
type SecretHandler struct {
	secretHandler *v2.SecretHandler
	BaseURL       string
	AuthToken     string
	AuthHeader    string
	HTTPClient    *http.Client
	Scheme        string
}

// NewSecretHandler returns a new SecretHandler which sends all requests directly to the secret-service
func NewSecretHandler(baseURL string) *SecretHandler {
	return NewSecretHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewSecretHandlerWithHTTPClient returns a new SecretHandler which sends all requests directly to the secret-service using the specified http.Client
func NewSecretHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *SecretHandler {
	return &SecretHandler{
		BaseURL:       httputils.TrimHTTPScheme(baseURL),
		HTTPClient:    httpClient,
		Scheme:        "http",
		secretHandler: v2.NewSecretHandlerWithHTTPClient(baseURL, httpClient),
	}
}

// NewAuthenticatedSecretHandler returns a new SecretHandler that authenticates at the api via the provided token
// and sends all requests directly to the secret-service
// Deprecated: use APISet instead
func NewAuthenticatedSecretHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *SecretHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedSecretHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedSecretHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *SecretHandler {
	v2SecretHandler := v2.NewAuthenticatedSecretHandler(baseURL, authToken, authHeader, httpClient, scheme)

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, secretServiceBaseURL) {
		baseURL += "/" + secretServiceBaseURL
	}

	return &SecretHandler{
		BaseURL:       httputils.TrimHTTPScheme(baseURL),
		AuthHeader:    authHeader,
		AuthToken:     authToken,
		HTTPClient:    httpClient,
		Scheme:        scheme,
		secretHandler: v2SecretHandler,
	}
}

func (s *SecretHandler) getBaseURL() string {
	return s.BaseURL
}

func (s *SecretHandler) getAuthToken() string {
	return s.AuthToken
}

func (s *SecretHandler) getAuthHeader() string {
	return s.AuthHeader
}

func (s *SecretHandler) getHTTPClient() *http.Client {
	return s.HTTPClient
}

// CreateSecret creates a new secret.
func (s *SecretHandler) CreateSecret(secret models.Secret) error {
	s.ensureHandlerIsSet()
	return s.secretHandler.CreateSecret(context.TODO(), secret, v2.SecretsCreateSecretOptions{})
}

// UpdateSecret creates a new secret.
func (s *SecretHandler) UpdateSecret(secret models.Secret) error {
	s.ensureHandlerIsSet()
	return s.secretHandler.UpdateSecret(context.TODO(), secret, v2.SecretsUpdateSecretOptions{})
}

// DeleteSecret deletes a secret.
func (s *SecretHandler) DeleteSecret(secretName, secretScope string) error {
	s.ensureHandlerIsSet()
	return s.secretHandler.DeleteSecret(context.TODO(), secretName, secretScope, v2.SecretsDeleteSecretOptions{})
}

// GetSecrets returns a list of created secrets.
func (s *SecretHandler) GetSecrets() (*models.GetSecretsResponse, error) {
	s.ensureHandlerIsSet()
	return s.secretHandler.GetSecrets(context.TODO(), v2.SecretsGetSecretsOptions{})
}

func (s *SecretHandler) ensureHandlerIsSet() {
	if s.secretHandler != nil {
		return
	}

	if s.AuthToken != "" {
		s.secretHandler = v2.NewAuthenticatedSecretHandler(s.BaseURL, s.AuthToken, s.AuthHeader, s.HTTPClient, s.Scheme)
	} else {
		s.secretHandler = v2.NewSecretHandlerWithHTTPClient(s.BaseURL, s.HTTPClient)
	}
}
