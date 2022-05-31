package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
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
	secretHandler v2.SecretHandler
	BaseURL       string
	AuthToken     string
	AuthHeader    string
	HTTPClient    *http.Client
	Scheme        string
}

// NewSecretHandler returns a new SecretHandler which sends all requests directly to the secret-service
func NewSecretHandler(baseURL string) *SecretHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}

	httpClient := &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))}

	return &SecretHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: httpClient,
		Scheme:     "http",

		secretHandler: v2.SecretHandler{
			BaseURL:    baseURL,
			AuthHeader: "",
			AuthToken:  "",
			HTTPClient: httpClient,
			Scheme:     "http",
		},
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
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")

	baseURL = strings.TrimRight(baseURL, "/")

	if !strings.HasSuffix(baseURL, secretServiceBaseURL) {
		baseURL += "/" + secretServiceBaseURL
	}

	return &SecretHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,

		secretHandler: v2.SecretHandler{
			BaseURL:    baseURL,
			AuthHeader: authHeader,
			AuthToken:  authToken,
			HTTPClient: httpClient,
			Scheme:     scheme,
		},
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
	return s.secretHandler.CreateSecret(context.TODO(), secret, v2.SecretsCreateSecretOptions{})
}

// UpdateSecret creates a new secret.
func (s *SecretHandler) UpdateSecret(secret models.Secret) error {
	return s.secretHandler.UpdateSecret(context.TODO(), secret, v2.SecretsUpdateSecretOptions{})
}

// DeleteSecret deletes a secret.
func (s *SecretHandler) DeleteSecret(secretName, secretScope string) error {
	return s.secretHandler.DeleteSecret(context.TODO(), secretName, secretScope, v2.SecretsDeleteSecretOptions{})
}

// GetSecrets returns a list of created secrets.
func (s *SecretHandler) GetSecrets() (*models.GetSecretsResponse, error) {
	return s.secretHandler.GetSecrets(context.TODO(), v2.SecretsGetSecretsOptions{})
}
