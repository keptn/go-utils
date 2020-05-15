package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

// EventHandler handles services
type EventHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// EventFilter allows to filter events based on the provided properties
type EventFilter struct {
	Project      string
	Stage        string
	Service      string
	EventType    string
	KeptnContext string
	EventID      string
}

// NewEventHandler returns a new EventHandler
func NewEventHandler(baseURL string) *EventHandler {
	scheme := "http"
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
		scheme = "http"
	}
	return &EventHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     scheme,
	}
}

// NewAuthenticatedEventHandler returns a new EventHandler that authenticates at the endpoint via the provided token
func NewAuthenticatedEventHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *EventHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &EventHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (e *EventHandler) getBaseURL() string {
	return e.BaseURL
}

func (e *EventHandler) getAuthToken() string {
	return e.AuthToken
}

func (e *EventHandler) getAuthHeader() string {
	return e.AuthHeader
}

func (e *EventHandler) getHTTPClient() *http.Client {
	return e.HTTPClient
}

// SendEvent sends an event to Keptn
func (e *EventHandler) SendEvent(event models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(event)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(e.Scheme+"://"+e.getBaseURL()+"/v1/event", bodyStr, e)
}

// GetEvents returns all events matching the properties in the passed filter object
func (e *EventHandler) GetEvents(filter *EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error) {
	raw := e.Scheme + "://" + e.getBaseURL() + "/mongodb-datastore/event?"

	u, _ := url.Parse(raw)

	query := u.Query()

	if filter.Project != "" {
		query.Set("project", filter.Project)
	}
	if filter.Stage != "" {
		query.Set("stage", filter.Stage)
	}
	if filter.Service != "" {
		query.Set("service", filter.Service)
	}
	if filter.KeptnContext != "" {
		query.Set("keptnContext", filter.KeptnContext)
	}
	if filter.EventID != "" {
		query.Set("eventID", filter.EventID)
	}

	u.RawQuery = query.Encode()

	return e.getEvents(u.String(), e)
}

// GetEvent returns an event specified by keptnContext and eventType
//
// Deprecated: this function is deprecated and should be replaced with the GetEvents function
func (e *EventHandler) GetEvent(keptnContext string, eventType string) (*models.KeptnContextExtendedCE, *models.Error) {
	return getEvent(e.Scheme+"://"+e.getBaseURL()+"/v1/event?keptnContext="+keptnContext+"&type="+eventType+"&pageSize=10", e)
}

func getEvent(uri string, api APIService) (*models.KeptnContextExtendedCE, *models.Error) {

	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Host = "api.keptn"
	addAuthHeader(req, api)

	resp, err := api.getHTTPClient().Do(req)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		if len(body) > 0 {
			var cloudEvent models.KeptnContextExtendedCE
			err = json.Unmarshal(body, &cloudEvent)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}

			return &cloudEvent, nil
		}

		return nil, nil
	}

	var respErr models.Error
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	return nil, &respErr
}

func (e *EventHandler) getEvents(uri string, api APIService) ([]*models.KeptnContextExtendedCE, *models.Error) {
	events := []*models.KeptnContextExtendedCE{}
	nextPageKey := ""

	for {
		url, err := url.Parse(uri)
		if err != nil {
			return nil, buildErrorResponse(err.Error())
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, e)

		resp, err := e.HTTPClient.Do(req)
		if err != nil {
			return nil, buildErrorResponse(err.Error())
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, buildErrorResponse(err.Error())
		}

		if resp.StatusCode == 200 {
			received := &models.Events{}
			err = json.Unmarshal(body, received)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}
			events = append(events, received.Events...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey
		} else {
			var respErr models.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}
			return nil, &respErr
		}
	}

	return events, nil
}
