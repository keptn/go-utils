package utils

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/keptn/go-utils/pkg/mongodb-datastore/models"
)

// EventHandler handles event
type EventHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewEventHandler returns a new EventHandler
func NewEventHandler(baseURL string) *EventHandler {
	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &EventHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{},
		Scheme:     "http",
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

// GetEvent returns the latest event of a specific event type and from a specific Keptn context
func (e *EventHandler) GetEvent(keptnContext string, eventType string) (*models.KeptnContextExtendedCE, *models.Error) {
	return getLatestEvent(keptnContext, eventType, e.Scheme+"://"+e.getBaseURL()+"/event?keptnContext="+keptnContext+"&type="+eventType+"&pageSize=10", e)
}

func getLatestEvent(keptnContext string, eventType string, uri string, datastore Datastore) (*models.KeptnContextExtendedCE, *models.Error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := datastore.getHTTPClient().Do(req)
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

			response := models.Events{}
			err = json.Unmarshal(body, &response)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}

			// find latest event
			var latest *models.KeptnContextExtendedCE
			for _, event := range response.Events {
				if latest == nil || time.Time(latest.Time).Before(time.Time(event.Time)) {
					latest = event
				}
			}

			if latest != nil {
				return latest, nil
			}
		}

		var respMessage models.Error
		message := "No Keptn " + eventType + " event found for context: " + keptnContext
		respMessage.Message = &message
		respMessage.Code = 400
		return nil, &respMessage
	}

	var respErr models.Error
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	return nil, &respErr
}
