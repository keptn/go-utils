package utils

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

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

func (p *EventHandler) getBaseURL() string {
	return p.BaseURL
}

func (p *EventHandler) getAuthToken() string {
	return p.AuthToken
}

func (p *EventHandler) getAuthHeader() string {
	return p.AuthHeader
}

func (p *EventHandler) getHTTPClient() *http.Client {
	return p.HTTPClient
}

// GetEvent returns an event specified by keptnContext and eventType
func (p *EventHandler) GetEvent(keptnContext string, eventType string) (*models.KeptnContextExtendedCE, *models.Error) {
	return get(p.Scheme+"://"+p.getBaseURL()+"/event?keptnContext="+keptnContext+"type="+eventType+"&pageSize=10", p)
}

func get(uri string, datastore Datastore) (*models.KeptnContextExtendedCE, *models.Error) {

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
