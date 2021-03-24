package api

import (
	"encoding/json"
	"github.com/keptn/go-utils/pkg/api/models"
	"net/http"
	"strings"
)

const secretServiceBaseURL = "secrets"
const v1SecretPath = "/v1/secrets"

//go:generate moq -pkg utils_mock -skip-ensure -out ./fake/secret_handler_mock.go . SecretHandlerInterface
type SecretHandlerInterface interface {
	CreateSecret(secret models.Secret) (string, *models.Error)
	UpdateSecret(secret models.Secret) (string, *models.Error)
	DeleteSecret(secretName, secretScope string) (string, *models.Error)
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
func (s *SecretHandler) CreateSecret(secret models.Secret) (string, *models.Error) {
	body, err := json.Marshal(secret)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return post(s.Scheme+"://"+s.BaseURL+v1SecretPath, body, s)
}

// UpdateSecret creates a new secret
func (s *SecretHandler) UpdateSecret(secret models.Secret) (string, *models.Error) {
	body, err := json.Marshal(secret)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return put(s.Scheme+"://"+s.BaseURL+v1SecretPath, body, s)
}

// DeleteSecret deletes a secret
func (s *SecretHandler) DeleteSecret(secretName, secretScope string) (string, *models.Error) {
	return delete(s.Scheme+"://"+s.BaseURL+v1SecretPath+"?name="+secretName+"&scope="+secretScope, s)
}
