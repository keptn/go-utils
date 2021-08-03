package api

import (
	"encoding/json"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/httputils"
	"net/http"
	"strings"
)

const v1SequenceControlPath = "/v1/sequence/%s/%s/control"

type SequenceControlHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

type SequenceControlParams struct {
	Project      string `json:"project"`
	KeptnContext string `json:"keptnContext"`
	Stage        string `json:"stage"`
	State        string `json:"state"`
}

func (s *SequenceControlParams) Validate() error {
	var missingFieldsErr []string
	if s.Project == "" {
		missingFieldsErr = append(missingFieldsErr, "project parameter not set")
	}
	if s.KeptnContext == "" {
		missingFieldsErr = append(missingFieldsErr, "keptn context parameter not set")
	}
	if s.State == "" {
		missingFieldsErr = append(missingFieldsErr, "sequence state parameter not set")
	}
	errStr := strings.Join(missingFieldsErr, ",")

	if len(missingFieldsErr) > 0 {
		return fmt.Errorf("failed to validate sequence control parameters: %s", errStr)
	}
	return nil
}

type SequenceControlBody struct {
	Stage string `json:"stage"`
	State string `json:"state"`
}

func (s *SequenceControlBody) ToJSON() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

type PauseSequenceParams struct {
}

type ResumeSequenceParams struct {
}

func NewSequenceControlHandler(baseURL string) *SequenceControlHandler {
	baseURL = httputils.TrimHTTPScheme(baseURL)
	return &SequenceControlHandler{
		BaseURL:    baseURL,
		AuthHeader: "",
		AuthToken:  "",
		HTTPClient: &http.Client{Transport: getClientTransport()},
		Scheme:     "http",
	}
}

func NewAuthenticatedSequenceControlHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *SequenceControlHandler {
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

	return &SequenceControlHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (s *SequenceControlHandler) getBaseURL() string {
	return s.BaseURL
}

func (s *SequenceControlHandler) getAuthToken() string {
	return s.AuthToken
}

func (s *SequenceControlHandler) getAuthHeader() string {
	return s.AuthHeader
}

func (s *SequenceControlHandler) getHTTPClient() *http.Client {
	return s.HTTPClient
}

func (s *SequenceControlHandler) AbortSequence(params SequenceControlParams) error {
	err := params.Validate()
	if err != nil {
		return err
	}

	baseurl := fmt.Sprintf("%s://%s", s.Scheme, s.getBaseURL())
	path := fmt.Sprintf(v1SequenceControlPath, params.Project, params.KeptnContext)

	body := SequenceControlBody{
		Stage: params.Stage,
		State: params.State,
	}

	payload, err := body.ToJSON()
	if err != nil {
		return err
	}

	_, errResponse := post(baseurl+path, payload, s)
	if errResponse != nil {
		return fmt.Errorf(errResponse.GetMessage())
	}

	return nil
}

func (s *SequenceControlHandler) PauseSequence(params PauseSequenceParams) error {
	return nil
}

func (s *SequenceControlHandler) ResumeSequence(params ResumeSequenceParams) error {
	return nil
}
