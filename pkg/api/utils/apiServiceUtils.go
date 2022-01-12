package api

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/keptn/go-utils/pkg/api/models"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
		Proxy:           http.ProxyFromEnvironment,
	}
	return tr
}

// Wraps the provided http.RoundTripper with one that
// starts a span and injects the span context into the outbound request headers.
func getInstrumentedClientTransport() *otelhttp.Transport {
	return otelhttp.NewTransport(getClientTransport())
}

func putWithEventContext(uri string, data []byte, api APIService) (*models.EventContext, *models.Error) {
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
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
			eventContext := &models.EventContext{}
			if err = eventContext.FromJSON(body); err != nil {
				// failed to parse json
				return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
			}

			if eventContext.KeptnContext != nil {
				fmt.Println("ID of Keptn context: " + *eventContext.KeptnContext)
			}
			return eventContext, nil
		}

		return nil, nil
	}

	if len(body) > 0 {
		respErr := &models.Error{}
		if err = respErr.FromJSON(body); err != nil {
			// failed to parse json
			return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
		}

		return nil, respErr
	}

	return nil, buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", resp.StatusCode, resp.Status))
}

func put(uri string, data []byte, api APIService) (string, *models.Error) {
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(data))
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, api)

	resp, err := api.getHTTPClient().Do(req)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 204 {
		if len(body) > 0 {
			return string(body), nil
		}

		return "", nil
	}

	if len(body) > 0 {
		respErr := &models.Error{}
		if err = respErr.FromJSON(body); err != nil {
			// failed to parse json
			return "", buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
		}

		return "", respErr
	}

	return "", buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", resp.StatusCode, resp.Status))
}

func postWithEventContext(uri string, data []byte, api APIService) (*models.EventContext, *models.Error) {
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
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
			eventContext := &models.EventContext{}
			if err = eventContext.FromJSON(body); err != nil {
				// failed to parse json
				return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
			}

			if eventContext.KeptnContext != nil {
				fmt.Println("ID of Keptn context: " + *eventContext.KeptnContext)
			}
			return eventContext, nil
		}

		return nil, nil
	}

	if len(body) > 0 {
		respErr := &models.Error{}
		if err = respErr.FromJSON(body); err != nil {
			// failed to parse json
			return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
		}

		return nil, respErr
	}

	return nil, buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", resp.StatusCode, resp.Status))
}

func post(uri string, data []byte, api APIService) (string, *models.Error) {
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, api)

	resp, err := api.getHTTPClient().Do(req)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 204 {
		if len(body) > 0 {
			return string(body), nil
		}

		return "", nil
	}

	if len(body) > 0 {
		respErr := &models.Error{}
		if err = respErr.FromJSON(body); err != nil {
			// failed to parse json
			return "", buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
		}

		return "", respErr
	}

	return "", buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", resp.StatusCode, resp.Status))
}

func deleteWithEventContext(uri string, api APIService) (*models.EventContext, *models.Error) {
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
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
			eventContext := &models.EventContext{}
			if err = eventContext.FromJSON(body); err != nil {
				// failed to parse json
				return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
			}
			return eventContext, nil
		}

		return nil, nil
	}

	respErr := &models.Error{}
	if err = respErr.FromJSON(body); err != nil {
		// failed to parse json
		return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
	}

	return nil, respErr
}

func delete(uri string, api APIService) (string, *models.Error) {
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, api)

	resp, err := api.getHTTPClient().Do(req)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if len(body) > 0 {
			return string(body), nil
		}

		return "", nil
	}

	respErr := &models.Error{}
	if err = respErr.FromJSON(body); err != nil {
		// failed to parse json
		return "", buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
	}

	return "", respErr
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
