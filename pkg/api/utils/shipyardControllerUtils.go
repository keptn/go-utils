package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

const shipyardControllerBaseURL = "shipyard-controller"

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
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     "http",
	}
}

// NewAuthenticatedShipyardControllerHandler returns a new ShipyardControllerHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
func NewAuthenticatedShipyardControllerHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ShipyardControllerHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

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
		url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + "/v1/event/triggered/" + filter.EventType)

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

		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("GET", url.String(), nil)
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
			err = json.Unmarshal(body, received)
			if err != nil {
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
			var respErr models.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, errors.New(*respErr.Message)
			}
			return nil, err
		}
	}
	return events, nil
}
