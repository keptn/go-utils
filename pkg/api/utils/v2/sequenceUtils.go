package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/common/httputils"
)

const v1SequenceControlPath = "/v1/sequence/%s/%s/control"

// SequencesControlSequenceOptions are options for SequencesInterface.ControlSequence().
type SequencesControlSequenceOptions struct{}

type SequencesInterface interface {
	ControlSequence(ctx context.Context, params SequenceControlParams, opts SequencesControlSequenceOptions) error
}

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
	var errMsg []string
	if s.Project == "" {
		errMsg = append(errMsg, "project parameter not set")
	}
	if s.KeptnContext == "" {
		errMsg = append(errMsg, "keptn context parameter not set")
	}
	if s.State == "" {
		errMsg = append(errMsg, "sequence state parameter not set")
	}
	errStr := strings.Join(errMsg, ",")

	if len(errMsg) > 0 {
		return fmt.Errorf("failed to validate sequence control parameters: %s", errStr)
	}
	return nil
}

type SequenceControlBody struct {
	Stage string `json:"stage"`
	State string `json:"state"`
}

// Converts object to JSON string
func (s *SequenceControlBody) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON converts JSON string to object
func (s *SequenceControlBody) FromJSON(b []byte) error {
	var res SequenceControlBody
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*s = res
	return nil
}

// NewSequenceControlHandlerWithHTTPClient returns a new SequenceControlHandler
func NewSequenceControlHandler(baseURL string) *SequenceControlHandler {
	return NewSequenceControlHandlerWithHTTPClient(baseURL, &http.Client{Transport: wrapOtelTransport(getClientTransport(nil))})
}

// NewSequenceControlHandlerWithHTTPClient returns a new SequenceControlHandler using the specified http.Client
func NewSequenceControlHandlerWithHTTPClient(baseURL string, httpClient *http.Client) *SequenceControlHandler {
	return createSequenceControlHandler(baseURL, "", "", httpClient, "http")
}

func createAuthenticatedSequenceControlHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *SequenceControlHandler {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, shipyardControllerBaseURL) {
		baseURL += "/" + shipyardControllerBaseURL
	}

	return createSequenceControlHandler(baseURL, authToken, authHeader, httpClient, scheme)
}

func createSequenceControlHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *SequenceControlHandler {
	return &SequenceControlHandler{
		BaseURL:    httputils.TrimHTTPScheme(baseURL),
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

func (s *SequenceControlHandler) ControlSequence(ctx context.Context, params SequenceControlParams, opts SequencesControlSequenceOptions) error {
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

	_, errResponse := post(ctx, baseurl+path, payload, s)
	if errResponse != nil {
		return fmt.Errorf(errResponse.GetMessage())
	}

	return nil
}
