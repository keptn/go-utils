package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/httputils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const uniformRegistrationBaseURL = "uniform/registration"
const v1UniformPath = "/v1/uniform/registration"

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
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     "http",
	}
}

// NewAuthenticatedSequenceControlHandler returns a new UniformHandler that authenticates at the api via the provided token
// Deprecated: use APISet instead
func NewAuthenticatedUniformHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *UniformHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()
	return createAuthenticatedUniformHandler(baseURL, authToken, authHeader, httpClient, scheme)
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

func (u *UniformHandler) Ping(integrationID string) (*models.Integration, error) {
	if integrationID == "" {
		return nil, errors.New("could not ping an invalid IntegrationID")
	}

	resp, err := put(u.Scheme+"://"+u.getBaseURL()+v1UniformPath+"/"+integrationID+"/ping", nil, u)
	if err != nil {
		return nil, errors.New(err.GetMessage())
	}

	response := &models.Integration{}
	if err := json.Unmarshal([]byte(resp), response); err != nil {
		return nil, err
	}

	return response, nil
}

func (u *UniformHandler) RegisterIntegration(integration models.Integration) (string, error) {
	bodyStr, err := integration.ToJSON()
	if err != nil {
		return "", err
	}

	resp, errResponse := post(u.Scheme+"://"+u.getBaseURL()+v1UniformPath, bodyStr, u)
	if errResponse != nil {
		return "", fmt.Errorf(errResponse.GetMessage())
	}

	registerIntegrationResponse := &models.RegisterIntegrationResponse{}
	if err := registerIntegrationResponse.FromJSON([]byte(resp)); err != nil {
		return "", err
	}

	return registerIntegrationResponse.ID, nil
}

func (u *UniformHandler) CreateSubscription(integrationID string, subscription models.EventSubscription) (string, error) {
	bodyStr, err := subscription.ToJSON()
	if err != nil {
		return "", err
	}
	resp, errResponse := post(u.Scheme+"://"+u.getBaseURL()+v1UniformPath+"/"+integrationID+"/subscription", bodyStr, u)
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

func (u *UniformHandler) UnregisterIntegration(integrationID string) error {
	_, err := delete(u.Scheme+"://"+u.getBaseURL()+v1UniformPath+"/"+integrationID, u)
	if err != nil {
		return fmt.Errorf(err.GetMessage())
	}
	return nil
}

func (u *UniformHandler) GetRegistrations() ([]*models.Integration, error) {
	url, err := url.Parse(u.Scheme + "://" + u.getBaseURL() + v1UniformPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, u)

	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var received []*models.Integration
		err := json.Unmarshal(body, &received)
		if err != nil {
			return nil, err
		}
		return received, nil
	}

	return nil, nil
}
