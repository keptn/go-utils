package utils

import (
	"net/http"

	"github.com/keptn/go-utils/pkg/datastore/models"
)

// APIService represents the interface for accessing the configuration service
type APIService interface {
	getBaseURL() string
	getAuthToken() string
	getAuthHeader() string
	getHTTPClient() *http.Client
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
