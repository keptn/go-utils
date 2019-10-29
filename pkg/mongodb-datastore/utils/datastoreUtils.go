package utils

import (
	"net/http"

	"github.com/keptn/go-utils/pkg/mongodb-datastore/models"
)

// Datastore represents the interface for accessing Keptn's datastore
type Datastore interface {
	getBaseURL() string
	getAuthToken() string
	getAuthHeader() string
	getHTTPClient() *http.Client
}

func buildErrorResponse(errorStr string) *models.Error {
	err := models.Error{Message: &errorStr}
	return &err
}

func addAuthHeader(req *http.Request, datastore Datastore) {
	if datastore.getAuthHeader() != "" && datastore.getAuthToken() != "" {
		req.Header.Set(datastore.getAuthHeader(), datastore.getAuthToken())
	}
}
