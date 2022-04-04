package api

import (
	"bytes"
	"context"
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

// createInstrumentedClientTransport tries to add support for opentelemetry
// to the given http.Client. If httpClient is nil, a fresh http.Client
// with opentelemetry support is created
func createInstrumentedClientTransport(httpClient *http.Client) *http.Client {
	if httpClient == nil {
		return &http.Client{
			Transport: wrapOtelTransport(getClientTransport(nil)),
		}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return httpClient
}

// Wraps the provided http.RoundTripper with one that
// starts a span and injects the span context into the outbound request headers.
func wrapOtelTransport(base http.RoundTripper) *otelhttp.Transport {
	return otelhttp.NewTransport(base)
}

// getClientTransport returns a client transport which
// skips verifying server certificates and is able to
// read proxy configuration from environment variables
//
// If the given http.RoundTripper is nil then a new http.Transport
// is created, otherwise the given http.RoundTripper is analysed whether it
// is of type *http.Transport. If so, the respective settings for
// disabling server certificate verification as well as proxy server support are set
// If not, the given http.RoundTripper is passed through untouched
func getClientTransport(rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
		}
		return tr
	}
	if tr, isDefaultTransport := rt.(*http.Transport); isDefaultTransport {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		tr.Proxy = http.ProxyFromEnvironment
		return tr
	}
	return rt
}

func getAndExpectOK(ctx context.Context, uri string, api APIService) ([]byte, *models.Error) {
	body, statusCode, status, err := get(ctx, uri, api)
	if err != nil {
		return nil, err
	}

	if statusCode == 200 {
		return body, nil
	}

	if len(body) > 0 {
		return nil, handleErrStatusCode(statusCode, body)
	}

	return nil, buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", statusCode, status))
}

func getAndExpectSuccess(ctx context.Context, uri string, api APIService) ([]byte, *models.Error) {
	body, statusCode, status, err := get(ctx, uri, api)
	if err != nil {
		return nil, err
	}

	if statusCode >= 200 && statusCode < 300 {
		return body, nil
	}

	if len(body) > 0 {
		return nil, handleErrStatusCode(statusCode, body)
	}

	return nil, buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", statusCode, status))
}

func get(ctx context.Context, uri string, api APIService) ([]byte, int, string, *models.Error) {
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, 0, "", buildErrorResponse(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, api)

	resp, err := api.getHTTPClient().Do(req)
	if err != nil {
		return nil, 0, "", buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", buildErrorResponse(err.Error())
	}

	return body, resp.StatusCode, resp.Status, nil
}

func putWithEventContext(ctx context.Context, uri string, data []byte, api APIService) (*models.EventContext, *models.Error) {
	req, err := http.NewRequestWithContext(ctx, "PUT", uri, bytes.NewBuffer(data))
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
		if len(body) == 0 {
			return nil, nil
		}

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

	if len(body) > 0 {
		return nil, handleErrStatusCode(resp.StatusCode, body)
	}

	return nil, buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", resp.StatusCode, resp.Status))
}

func put(ctx context.Context, uri string, data []byte, api APIService) (string, *models.Error) {
	req, err := http.NewRequestWithContext(ctx, "PUT", uri, bytes.NewBuffer(data))
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
		return string(body), nil
	}

	if len(body) > 0 {
		return "", handleErrStatusCode(resp.StatusCode, body)
	}

	return "", buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", resp.StatusCode, resp.Status))
}

func postWithEventContext(ctx context.Context, uri string, data []byte, api APIService) (*models.EventContext, *models.Error) {
	req, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewBuffer(data))
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
		if len(body) == 0 {
			return nil, nil
		}

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

	if len(body) > 0 {
		return nil, handleErrStatusCode(resp.StatusCode, body)
	}

	return nil, buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", resp.StatusCode, resp.Status))
}

func post(ctx context.Context, uri string, data []byte, api APIService) (string, *models.Error) {
	req, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewBuffer(data))
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
		return string(body), nil
	}

	if len(body) > 0 {
		return "", handleErrStatusCode(resp.StatusCode, body)
	}

	return "", buildErrorResponse(fmt.Sprintf("Received unexpected response: %d %s", resp.StatusCode, resp.Status))
}

func deleteWithEventContext(ctx context.Context, uri string, api APIService) (*models.EventContext, *models.Error) {
	req, err := http.NewRequestWithContext(ctx, "DELETE", uri, nil)
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
		if len(body) == 0 {
			return nil, nil
		}

		eventContext := &models.EventContext{}
		if err = eventContext.FromJSON(body); err != nil {
			// failed to parse json
			return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
		}
		return eventContext, nil
	}

	return nil, handleErrStatusCode(resp.StatusCode, body)
}

func delete(ctx context.Context, uri string, api APIService) (string, *models.Error) {
	req, err := http.NewRequestWithContext(ctx, "DELETE", uri, nil)
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
		return string(body), nil
	}

	return "", handleErrStatusCode(resp.StatusCode, body)
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
