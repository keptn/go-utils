package api

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

const shipyardControllerBaseURL = "controlPlane"

type ShipyardControlV1Interface interface {
	GetOpenTriggeredEvents(filter EventFilter) ([]*models.KeptnContextExtendedCE, error)
}

// ShipyardControllerHandler handles services
type ShipyardControllerHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewShipyardControllerHandler returns a new ShipyardControllerHandler which sends all requests directly to the configuration-service
func NewShipyardControllerHandler(baseURL string) *ShipyardControllerHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &ShipyardControllerHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))},
		Scheme:     "http",
	}
}

// NewAuthenticatedShipyardControllerHandler returns a new ShipyardControllerHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
// Deprecated: use APISet instead
func NewAuthenticatedShipyardControllerHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ShipyardControllerHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return createAuthenticatedShipyardControllerHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createAuthenticatedShipyardControllerHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ShipyardControllerHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}
	return &ShipyardControllerHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (s *ShipyardControllerHandler) getBaseURL() string {
	return s.BaseURL
}

func (s *ShipyardControllerHandler) getAuthToken() string {
	return s.AuthToken
}

func (s *ShipyardControllerHandler) getAuthHeader() string {
	return s.AuthHeader
}

func (s *ShipyardControllerHandler) getHTTPClient() *http.Client {
	return s.HTTPClient
}

// GetOpenTriggeredEvents returns all open triggered events
func (s *ShipyardControllerHandler) GetOpenTriggeredEvents(filter EventFilter) ([]*models.KeptnContextExtendedCE, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	events := []*models.KeptnContextExtendedCE{}
	nextPageKey := ""

	for {
		url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + v1EventPath + "/triggered/" + filter.EventType)

		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		if filter.Project != "" {
			q.Set("project", filter.Project)
		}
		if filter.Service != "" {
			q.Set("service", filter.Service)
		}
		if filter.Stage != "" {
			q.Set("stage", filter.Stage)
		}

		url.RawQuery = q.Encode()

		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("GET", url.String(), nil)
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

		if resp.StatusCode == 200 {
			received := &models.Events{}
			if err = received.FromJSON(body); err != nil {
				return nil, err
			}
			events = append(events, received.Events...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}

			nextPageKeyInt, _ := strconv.Atoi(received.NextPageKey)

			if filter.NumberOfPages > 0 && nextPageKeyInt >= filter.NumberOfPages {
				break
			}

			nextPageKey = received.NextPageKey
		} else {
			respErr := &models.Error{}
			if err = respErr.FromJSON(body); err == nil && respErr != nil {
				return nil, errors.New(*respErr.Message)
			} else {
				return nil, fmt.Errorf("error with, status code %d", resp.StatusCode)
			}
		}
	}
	return events, nil
}
