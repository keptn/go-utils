package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiSetWithInvalidURL(t *testing.T) {
	apiSet, err := NewApiSet("://http.lol", "a-token", "x-token", nil, "http")
	assert.Nil(t, apiSet)
	assert.Error(t, err)
}

func TestApiSetCreatesHandlers(t *testing.T) {
	apiSet, err := NewApiSet("http://base-url.com", "a-token", "x-token", nil, "http")
	assert.NoError(t, err)
	assert.Equal(t, "a-token", apiSet.Token())
	assert.Equal(t, "http://base-url.com", apiSet.Endpoint().String())
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
