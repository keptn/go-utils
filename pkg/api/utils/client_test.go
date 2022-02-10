package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.NotNil(t, apiSet.ShipyardControlV1())
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

func TestApiSetInternalMappings(t *testing.T) {
	internal, err := NewInternal(nil)
	require.Nil(t, err)
	require.NotNil(t, internal)
	assert.Equal(t, DefaultInClusterAPIMappings[MongoDBDatastore], internal.EventsV1().(*EventHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ApiService], internal.AuthV1().(*AuthHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.APIV1().(*APIHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.ShipyardControlV1().(*ShipyardControllerHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.UniformV1().(*UniformHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.LogsV1().(*LogHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.SequencesV1().(*SequenceControlHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.StagesV1().(*StageHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[SecretService], internal.SecretsV1().(*SecretHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ConfigurationService], internal.ResourcesV1().(*ResourceHandler).BaseURL)
	assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.ProjectsV1().(*ProjectHandler).BaseURL)
}

type roundTripFn func(req *http.Request) *http.Response

func (f roundTripFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn roundTripFn) *http.Client {
	return &http.Client{
		Transport: roundTripFn(fn),
	}
}
