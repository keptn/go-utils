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

// ServiceHandler handles services
type ServiceHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewServiceHandler returns a new ServiceHandler which sends all requests directly to the configuration-service
func NewServiceHandler(baseURL string) *ServiceHandler {
	if strings.Contains(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	} else if strings.Contains(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	}
	return &ServiceHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     "http",
	}
}

// NewAuthenticatedServiceHandler returns a new ServiceHandler that authenticates at the api via the provided token
// and sends all requests directly to the configuration-service
func NewAuthenticatedServiceHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *ServiceHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, configurationServiceBaseUrl) {
		baseURL += "/" + configurationServiceBaseUrl
	}
	return &ServiceHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (s *ServiceHandler) getBaseURL() string {
	return s.BaseURL
}

func (s *ServiceHandler) getAuthToken() string {
	return s.AuthToken
}

func (s *ServiceHandler) getAuthHeader() string {
	return s.AuthHeader
}

func (s *ServiceHandler) getHTTPClient() *http.Client {
	return s.HTTPClient
}

// CreateService creates a new service
func (s *ServiceHandler) CreateServiceInStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error) {

	service := models.Service{ServiceName: serviceName}
	body, err := json.Marshal(service)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return post(s.Scheme+"://"+s.BaseURL+"/v1/project/"+project+"/stage/"+stage+"/service", body, s)
}

func (s *ServiceHandler) GetService(project, stage, service string) (*models.Service, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + "/v1/project/" + project + "/stage/" + stage + "/service/" + service)
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
		var received models.Service
		err = json.Unmarshal(body, &received)
		if err != nil {
			return nil, err
		}
		return &received, nil
	} else {
		var respErr models.Error
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Response Error Code: " + strconv.FormatInt(respErr.Code, 10) + " Message: " + *respErr.Message)
	}
}

// GetAllServices returns a list of all services.
func (s *ServiceHandler) GetAllServices(project string, stage string) ([]*models.Service, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	services := []*models.Service{}

	nextPageKey := ""

	for {
		url, err := url.Parse(s.Scheme + "://" + s.getBaseURL() + "/v1/project/" + project + "/stage/" + stage + "/service")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
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
			var received models.Services
			err = json.Unmarshal(body, &received)
			if err != nil {
				return nil, err
			}
			services = append(services, received.Services...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey
		} else {
			var respErr models.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, err
			}
			return nil, errors.New("Response Error Code: " + strconv.FormatInt(respErr.Code, 10) + " Message: " + *respErr.Message)
		}
	}

	return services, nil
}
