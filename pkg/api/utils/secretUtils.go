package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

const secretServiceBaseURL = "secrets"
const v1SecretPath = "/v1/secret"

//go:generate moq -pkg utils_mock -skip-ensure -out ./fake/secret_handler_mock.go . SecretHandlerInterface
type SecretHandlerInterface interface {
	CreateSecret(secret models.Secret) error
	UpdateSecret(secret models.Secret) error
	DeleteSecret(secretName, secretScope string) error
	GetSecrets() (*models.GetSecretsResponse, error)
}

// SecretHandler handles services
type SecretHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewSecretHandler returns a new SecretHandler which sends all requests directly to the secret-service
func NewSecretHandler(baseURL string) *SecretHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &SecretHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     "http",
	}
}

// NewAuthenticatedSecretHandler returns a new SecretHandler that authenticates at the api via the provided token
// and sends all requests directly to the secret-service
func NewAuthenticatedSecretHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *SecretHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

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

// CreateSecret creates a new secret
//
// Deprecated: Use CreateSecretWithContext instead
func (s *SecretHandler) CreateSecret(secret models.Secret) error {
	return s.CreateSecretWithContext(context.Background(), secret)
}

// CreateSecretWithContext creates a new secret
func (s *SecretHandler) CreateSecretWithContext(ctx context.Context, secret models.Secret) error {
	body, err := json.Marshal(secret)
	if err != nil {
		return err
	}
	_, errObj := post(ctx, s.Scheme+"://"+s.BaseURL+v1SecretPath, body, s)
	if errObj != nil {
		return errors.New(errObj.GetMessage())
	}
	return nil
}

// UpdateSecret creates a new secret
//
// Deprecated: Use UpdateSecretWithContext instead
func (s *SecretHandler) UpdateSecret(secret models.Secret) error {
	return s.UpdateSecretWithContext(context.Background(), secret)
}

// UpdateSecretWithContext creates a new secret
func (s *SecretHandler) UpdateSecretWithContext(ctx context.Context, secret models.Secret) error {
	body, err := json.Marshal(secret)
	if err != nil {
		return err
	}
	_, errObj := put(ctx, s.Scheme+"://"+s.BaseURL+v1SecretPath, body, s)
	if errObj != nil {
		return errors.New(errObj.GetMessage())
	}
	return nil
}

// DeleteSecret deletes a secret
//
// Deprecated: Use DeleteSecretWithContext instead
func (s *SecretHandler) DeleteSecret(secretName, secretScope string) error {
	return s.DeleteSecretWithContext(context.Background(), secretName, secretScope)
}

// DeleteSecretWithContext deletes a secret
func (s *SecretHandler) DeleteSecretWithContext(ctx context.Context, secretName, secretScope string) error {
	_, err := delete(ctx, s.Scheme+"://"+s.BaseURL+v1SecretPath+"?name="+secretName+"&scope="+secretScope, s)
	if err != nil {
		return errors.New(err.GetMessage())
	}
	return nil
}

// GetSecrets returns a list of created secrets
//
// Deprecated: Use GetSecretsWithContext instead
func (s *SecretHandler) GetSecrets() (*models.GetSecretsResponse, error) {
	return s.GetSecretsWithContext(context.Background())
}

// GetSecretsWithContext returns a list of created secrets
func (s *SecretHandler) GetSecretsWithContext(ctx context.Context) (*models.GetSecretsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.Scheme+"://"+s.BaseURL+v1SecretPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, s)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errObj := &models.Error{}
		if err := json.Unmarshal(body, errObj); err != nil {
			return nil, err
		}
		return nil, errors.New(*errObj.Message)
	}
	result := &models.GetSecretsResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}
