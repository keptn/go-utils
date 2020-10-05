package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
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

func getClientTransport(clientCertPath, clientKeyPath, rootCertPath string) *http.Transport {
	var tlsConfig *tls.Config

	if clientCertPath == "" || clientKeyPath == "" || rootCertPath == "" {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		caCert, _ := ioutil.ReadFile(rootCertPath)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		cert, _ := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		tlsConfig = &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{cert},
		}
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialContext:     ResolveXipIoWithContext,
	}
	return tr
}

func put(uri string, data []byte, api APIService) (*models.EventContext, *models.Error) {

	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(data))
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
			var eventContext models.EventContext
			err = json.Unmarshal(body, &eventContext)
			if err != nil {
				// failed to parse json
				return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
			}

			if eventContext.KeptnContext != nil {
				fmt.Println("ID of Keptn context: " + *eventContext.KeptnContext)
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

func post(uri string, data []byte, api APIService) (*models.EventContext, *models.Error) {

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
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
			var eventContext models.EventContext
			err = json.Unmarshal(body, &eventContext)
			if err != nil {
				// failed to parse json
				return nil, buildErrorResponse(err.Error() + "\n" + "-----DETAILS-----" + string(body))
			}

			if eventContext.KeptnContext != nil {
				fmt.Println("ID of Keptn context: " + *eventContext.KeptnContext)
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
