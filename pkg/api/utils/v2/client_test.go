package v2

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NotNil(t, apiSet.Uniform())
	assert.NotNil(t, apiSet.Endpoint())
	assert.NotNil(t, apiSet.ShipyardControl())
	assert.NotNil(t, apiSet.Stages())
	assert.NotNil(t, apiSet.Services())
	assert.NotNil(t, apiSet.Sequences())
	assert.NotNil(t, apiSet.Secrets())
	assert.NotNil(t, apiSet.Projects())
	assert.NotNil(t, apiSet.API())
	assert.NotNil(t, apiSet.Events())
	assert.NotNil(t, apiSet.Auth())
	assert.NotNil(t, apiSet.Resources())
	assert.NotNil(t, apiSet.Logs())
}

func TestAPISetDefaultValues(t *testing.T) {
	apiSet, err := New("base-url.com")
	assert.Nil(t, err)
	assert.NotNil(t, apiSet)
	assert.Equal(t, "http", apiSet.scheme)
	assert.Equal(t, "", apiSet.authHeader)
	assert.Equal(t, "", apiSet.apiToken)
	assert.NotNil(t, apiSet.httpClient)

	apiSet, err = New("https://base-url.com")
	assert.Nil(t, err)
	assert.NotNil(t, apiSet)
	assert.Equal(t, "https", apiSet.scheme)
	assert.Equal(t, "", apiSet.authHeader)
	assert.Equal(t, "", apiSet.apiToken)
	assert.NotNil(t, apiSet.httpClient)
}

func TestAPISetWithOptions(t *testing.T) {
	apiSet, err := New("base-url.com", WithAuthToken("a-token"), WithHTTPClient(&http.Client{}), WithScheme("https"))
	assert.NoError(t, err)
	assert.Equal(t, "a-token", apiSet.Token())
	assert.Equal(t, "x-token", apiSet.authHeader)
	assert.Equal(t, "https", apiSet.scheme)
	assert.NotNil(t, apiSet.httpClient)
}
