package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/keptn/go-utils/pkg/api/models"
)

// APIService represents the interface for accessing the configuration service
type APIService interface {
	getBaseURL() string
	getAuthToken() string
	getAuthHeader() string
	getHTTPClient() *http.Client
}

func getClientTransport() *http.Transport {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext:     ResolveXipIoWithContext,
	}
	return tr
}

func post(uri string, data []byte, api APIService) (*models.EventContext, *models.Error) {

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
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

	if resp.StatusCode >= 200 && resp.StatusCode <= 204 {
		if len(body) > 0 {
			var eventContext models.EventContext
			err = json.Unmarshal(body, &eventContext)
			if err != nil {
				// failed to parse json
				return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
			}

			if eventContext.KeptnContext != nil {
				fmt.Println("ID of Keptn context: " + *eventContext.KeptnContext)
			} else {
				fmt.Println("ID of Keptn context is nil")
			}
			return &eventContext, nil
		}

		return nil, nil
	}

	if len(body) > 0 {
		var respErr models.Error
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			// failed to parse json
			return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
		}

		return nil, &respErr
	}

	return nil, buildErrorResponse(fmt.Sprintf("Received unexptected response: %d %s", resp.StatusCode, resp.Status))
}

func delete(uri string, api APIService) (*models.EventContext, *models.Error) {

	req, err := http.NewRequest("DELETE", uri, nil)
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
			var eventContext models.EventContext
			err = json.Unmarshal(body, &eventContext)
			if err != nil {
				// failed to parse json
				return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
			}
			return &eventContext, nil
		}

		return nil, nil
	}

	var respErr models.Error
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		// failed to parse json
		return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
	}

	return nil, &respErr
}

func buildErrorResponse(errorStr string) *models.Error {
	err := models.Error{Message: &errorStr}
	return &err
}

func addAuthHeader(req *http.Request, api APIService) {
	if api.getAuthHeader() != "" && api.getAuthToken() != "" {
		req.Header.Set(api.getAuthHeader(), api.getAuthToken())
	}
}
