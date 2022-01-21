package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestApiSetWithInvalidURL(t *testing.T) {
	apiSet, err := New("://http.lol")
	assert.Nil(t, apiSet)
	assert.Error(t, err)
}

func TestApiSetCreatesHandlers(t *testing.T) {
	apiSet, err := New("http://base-url.com")
	assert.NoError(t, err)
	assert.Equal(t, "http://base-url.com", apiSet.Endpoint().String())
	assert.Equal(t, "http", apiSet.scheme)
	assert.NotNil(t, apiSet.UniformV1())
	assert.NotNil(t, apiSet.Endpoint())
	assert.NotNil(t, apiSet.ShipyardControlHandlerV1())
	assert.NotNil(t, apiSet.StagesV1())
	assert.NotNil(t, apiSet.ServicesV1())
	assert.NotNil(t, apiSet.SequencesV1())
	assert.NotNil(t, apiSet.SecretsV1())
	assert.NotNil(t, apiSet.ProjectsV1())
	assert.NotNil(t, apiSet.APIV1())
	assert.NotNil(t, apiSet.EventsV1())
	assert.NotNil(t, apiSet.AuthV1())
	assert.NotNil(t, apiSet.ResourcesV1())
	assert.NotNil(t, apiSet.LogsV1())
}

func TestAPISetWithOptions(t *testing.T) {
	apiSet, err := New("https://base-url.com", WithAuthToken("a-token"), WithHTTPClient(&http.Client{}))
	assert.NoError(t, err)
	assert.Equal(t, "a-token", apiSet.Token())
	assert.Equal(t, "x-token", apiSet.authHeader)
	assert.Equal(t, "https", apiSet.scheme)
	assert.NotNil(t, apiSet.httpClient)
}
