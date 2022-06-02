package v2

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const secretServiceBaseURL = "secrets"
const v1SecretPath = "/v1/secret"

// SecretsCreateSecretOptions are options for SecretsInterface.CreateSecret().
type SecretsCreateSecretOptions struct{}

// SecretsUpdateSecretOptions are options for SecretsInterface.UpdateSecret().
type SecretsUpdateSecretOptions struct{}

// SecretsDeleteSecretOptions are options for SecretsInterface.DeleteSecret().
type SecretsDeleteSecretOptions struct{}

// SecretsGetSecretsOptions are options for SecretsInterface.GetSecrets().
type SecretsGetSecretsOptions struct{}

//go:generate moq -pkg utils_mock -skip-ensure -out ./fake/secret_handler_mock.go . SecretsInterface
type SecretsInterface interface {
	// CreateSecret creates a new secret.
	CreateSecret(ctx context.Context, secret models.Secret, opts SecretsCreateSecretOptions) error

	// UpdateSecret creates a new secret.
	UpdateSecret(ctx context.Context, secret models.Secret, opts SecretsUpdateSecretOptions) error

	// DeleteSecret deletes a secret.
	DeleteSecret(ctx context.Context, secretName, secretScope string, opts SecretsDeleteSecretOptions) error

	// GetSecrets returns a list of created secrets.
	GetSecrets(ctx context.Context, opts SecretsGetSecretsOptions) (*models.GetSecretsResponse, error)
}

// SecretHandler handles secrets
type SecretHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewSecretHandler returns a new SecretHandler which sends all requests directly to the secret-service
func NewSecretHandler(baseURL string) *SecretHandler {
	return NewSecretHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewSecretHandlerWithHTTPClient returns a new SecretHandler which sends all requests directly to the secret-service using the specified http.Client
func NewSecretHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *SecretHandler {
	return createSecretHandler(baseURL, "", "", httpClient, "http")
}

func createAuthenticatedSecretHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *SecretHandler {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, secretServiceBaseURL) {
		baseURL += "/" + secretServiceBaseURL
	}

	return createSecretHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createSecretHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *SecretHandler {
	return &SecretHandler{
		BaseURL:    httputils.TrimHTTPScheme(baseURL),
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
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
func (s *SecretHandler) CreateSecret(ctx context.Context, secret models.Secret, opts SecretsCreateSecretOptions) error {
	body, err := secret.ToJSON()
	if err != nil {
		return err
	}
	_, errObj := post(ctx, s.Scheme+"://"+s.BaseURL+v1SecretPath, body, s)
	if errObj != nil {
		return errors.New(errObj.GetMessage())
	}
	return nil
}

// UpdateSecret creates a new secret.
func (s *SecretHandler) UpdateSecret(ctx context.Context, secret models.Secret, opts SecretsUpdateSecretOptions) error {
	body, err := secret.ToJSON()
	if err != nil {
		return err
	}
	_, errObj := put(ctx, s.Scheme+"://"+s.BaseURL+v1SecretPath, body, s)
	if errObj != nil {
		return errors.New(errObj.GetMessage())
	}
	return nil
}

// DeleteSecret deletes a secret.
func (s *SecretHandler) DeleteSecret(ctx context.Context, secretName, secretScope string, opts SecretsDeleteSecretOptions) error {
	_, err := delete(ctx, s.Scheme+"://"+s.BaseURL+v1SecretPath+"?name="+secretName+"&scope="+secretScope, s)
	if err != nil {
		return errors.New(err.GetMessage())
	}
	return nil
}

// GetSecrets returns a list of created secrets.
func (s *SecretHandler) GetSecrets(ctx context.Context, opts SecretsGetSecretsOptions) (*models.GetSecretsResponse, error) {
	body, mErr := getAndExpectOK(ctx, s.Scheme+"://"+s.BaseURL+v1SecretPath, s)
	if mErr != nil {
		return nil, mErr.ToError()
	}

	result := &models.GetSecretsResponse{}
	if err := result.FromJSON(body); err != nil {
		return nil, err
	}
	return result, nil
}
