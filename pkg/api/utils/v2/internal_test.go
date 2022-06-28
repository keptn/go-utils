package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiSetInternalMappings(t *testing.T) {
	t.Run("TestInternalAPISet - Default API Mappings", func(t *testing.T) {
		internal, err := NewInternal(nil)
		require.Nil(t, err)
		require.NotNil(t, internal)
		assert.Equal(t, DefaultInClusterAPIMappings[MongoDBDatastore], internal.Events().(*EventHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ApiService], internal.Auth().(*AuthHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.API().(*InternalAPIHandler).shipyardControllerApiHandler.baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.ShipyardControl().(*ShipyardControllerHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.Uniform().(*UniformHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.Logs().(*LogHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.Sequences().(*SequenceControlHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.Stages().(*StageHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[SecretService], internal.Secrets().(*SecretHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ConfigurationService], internal.Resources().(*ResourceHandler).baseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.Projects().(*ProjectHandler).baseURL)
	})

	t.Run("TestInternalAPISet - Override Mappings", func(t *testing.T) {
		overrideMappings := InClusterAPIMappings{
			ConfigurationService: "special-configuration-service:8080",
			ShipyardController:   "special-shipyard-controller:8080",
			ApiService:           "speclial-api-service:8080",
			SecretService:        "special-secret-service:8080",
			MongoDBDatastore:     "special-monogodb-datastore:8080",
		}
		internal, err := NewInternal(nil, overrideMappings)
		require.Nil(t, err)
		require.NotNil(t, internal)
		assert.Equal(t, overrideMappings[MongoDBDatastore], internal.Events().(*EventHandler).baseURL)
		assert.Equal(t, overrideMappings[ApiService], internal.Auth().(*AuthHandler).baseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.API().(*InternalAPIHandler).shipyardControllerApiHandler.baseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.ShipyardControl().(*ShipyardControllerHandler).baseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.Uniform().(*UniformHandler).baseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.Logs().(*LogHandler).baseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.Sequences().(*SequenceControlHandler).baseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.Stages().(*StageHandler).baseURL)
		assert.Equal(t, overrideMappings[SecretService], internal.Secrets().(*SecretHandler).baseURL)
		assert.Equal(t, overrideMappings[ConfigurationService], internal.Resources().(*ResourceHandler).baseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.Projects().(*ProjectHandler).baseURL)
	})

}
