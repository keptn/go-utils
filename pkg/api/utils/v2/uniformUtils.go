package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/httputils"
)

const uniformRegistrationBaseURL = "uniform/registration"
const v1UniformPath = "/v1/uniform/registration"

// UniformPingOptions are options for UniformInterface.Ping().
type UniformPingOptions struct{}

// UniformRegisterIntegrationOptions are options for UniformInterface.RegisterIntegration().
type UniformRegisterIntegrationOptions struct{}

// UniformCreateSubscriptionOptions are options for UniformInterface.CreateSubscription().
type UniformCreateSubscriptionOptions struct{}

// UniformUnregisterIntegrationOptions are options for UniformInterface.UnregisterIntegration().
type UniformUnregisterIntegrationOptions struct{}

// UniformGetRegistrationsOptions are options for UniformInterface.GetRegistrations().
type UniformGetRegistrationsOptions struct{}

type UniformInterface interface {
	Ping(ctx context.Context, integrationID string, opts UniformPingOptions) (*models.Integration, error)
	RegisterIntegration(ctx context.Context, integration models.Integration, opts UniformRegisterIntegrationOptions) (string, error)
	CreateSubscription(ctx context.Context, integrationID string, subscription models.EventSubscription, opts UniformCreateSubscriptionOptions) (string, error)
	UnregisterIntegration(ctx context.Context, integrationID string, opts UniformUnregisterIntegrationOptions) error
	GetRegistrations(ctx context.Context, opts UniformGetRegistrationsOptions) ([]*models.Integration, error)
}

type UniformHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

func NewUniformHandler(baseURL string) *UniformHandler {
	baseURL = httputils.TrimHTTPScheme(baseURL)
	return &UniformHandler{
		BaseURL:    baseURL,
		AuthToken:  "",
		AuthHeader: "",
		HTTPClient: &http.Client{Transport: getClientTransport(nil)},
		Scheme:     "http",
	}
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

func (u *UniformHandler) Ping(ctx context.Context, integrationID string, opts UniformPingOptions) (*models.Integration, error) {
	if integrationID == "" {
		return nil, errors.New("could not ping an invalid IntegrationID")
	}

	resp, err := put(ctx, u.Scheme+"://"+u.getBaseURL()+v1UniformPath+"/"+integrationID+"/ping", nil, u)
	if err != nil {
		return nil, errors.New(err.GetMessage())
	}

	response := &models.Integration{}
	if err := response.FromJSON([]byte(resp)); err != nil {
		return nil, err
	}

	return response, nil
}

func (u *UniformHandler) RegisterIntegration(ctx context.Context, integration models.Integration, opts UniformRegisterIntegrationOptions) (string, error) {
	bodyStr, err := integration.ToJSON()
	if err != nil {
		return "", err
	}

	resp, errResponse := post(ctx, u.Scheme+"://"+u.getBaseURL()+v1UniformPath, bodyStr, u)
	if errResponse != nil {
		return "", fmt.Errorf(errResponse.GetMessage())
	}

	registerIntegrationResponse := &models.RegisterIntegrationResponse{}
	if err := registerIntegrationResponse.FromJSON([]byte(resp)); err != nil {
		return "", err
	}

	return registerIntegrationResponse.ID, nil
}

func (u *UniformHandler) CreateSubscription(ctx context.Context, integrationID string, subscription models.EventSubscription, opts UniformCreateSubscriptionOptions) (string, error) {
	bodyStr, err := subscription.ToJSON()
	if err != nil {
		return "", err
	}
	resp, errResponse := post(ctx, u.Scheme+"://"+u.getBaseURL()+v1UniformPath+"/"+integrationID+"/subscription", bodyStr, u)
	if errResponse != nil {
		return "", fmt.Errorf(errResponse.GetMessage())
	}
	_ = resp

	createSubscriptionResponse := &models.CreateSubscriptionResponse{}
	if err := createSubscriptionResponse.FromJSON([]byte(resp)); err != nil {
		return "", err
	}

	return createSubscriptionResponse.ID, nil
}

func (u *UniformHandler) UnregisterIntegration(ctx context.Context, integrationID string, opts UniformUnregisterIntegrationOptions) error {
	_, err := delete(ctx, u.Scheme+"://"+u.getBaseURL()+v1UniformPath+"/"+integrationID, u)
	if err != nil {
		return fmt.Errorf(err.GetMessage())
	}
	return nil
}

func (u *UniformHandler) GetRegistrations(ctx context.Context, opts UniformGetRegistrationsOptions) ([]*models.Integration, error) {
	url, err := url.Parse(u.Scheme + "://" + u.getBaseURL() + v1UniformPath)
	if err != nil {
		return nil, err
	}

	body, mErr := getAndExpectOK(ctx, url.String(), u)
	if mErr != nil {
		return nil, mErr.ToError()
	}

	var received []*models.Integration
	err = json.Unmarshal(body, &received)
	if err != nil {
		return nil, err
	}
	return received, nil
}
