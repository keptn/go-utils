package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApiSetInternalMappings(t *testing.T) {
	t.Run("TestInternalAPISet - Default API Mappings", func(t *testing.T) {
		internal, err := NewInternal(nil)
		require.Nil(t, err)
		require.NotNil(t, internal)
		assert.Equal(t, DefaultInClusterAPIMappings[MongoDBDatastore], internal.EventsV1().(*EventHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ApiService], internal.AuthV1().(*AuthHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.APIV1().(*InternalAPIHandler).shipyardControllerApiHandler.BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.ShipyardControlV1().(*ShipyardControllerHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.UniformV1().(*UniformHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.LogsV1().(*LogHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.SequencesV1().(*SequenceControlHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.StagesV1().(*StageHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[SecretService], internal.SecretsV1().(*SecretHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ConfigurationService], internal.ResourcesV1().(*ResourceHandler).BaseURL)
		assert.Equal(t, DefaultInClusterAPIMappings[ShipyardController], internal.ProjectsV1().(*ProjectHandler).BaseURL)
	})

	t.Run("TestInternalAPISet - Override Mappings", func(t *testing.T) {
		overrideMappings := InClusterAPIMappings{
			ConfigurationService: "special-resource-service:8080",
			ShipyardController:   "special-shipyard-controller:8080",
			ApiService:           "speclial-api-service:8080",
			SecretService:        "special-secret-service:8080",
			MongoDBDatastore:     "special-monogodb-datastore:8080",
		}
		internal, err := NewInternal(nil, overrideMappings)
		require.Nil(t, err)
		require.NotNil(t, internal)
		assert.Equal(t, overrideMappings[MongoDBDatastore], internal.EventsV1().(*EventHandler).BaseURL)
		assert.Equal(t, overrideMappings[ApiService], internal.AuthV1().(*AuthHandler).BaseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.APIV1().(*InternalAPIHandler).shipyardControllerApiHandler.BaseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.ShipyardControlV1().(*ShipyardControllerHandler).BaseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.UniformV1().(*UniformHandler).BaseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.LogsV1().(*LogHandler).BaseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.SequencesV1().(*SequenceControlHandler).BaseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.StagesV1().(*StageHandler).BaseURL)
		assert.Equal(t, overrideMappings[SecretService], internal.SecretsV1().(*SecretHandler).BaseURL)
		assert.Equal(t, overrideMappings[ConfigurationService], internal.ResourcesV1().(*ResourceHandler).BaseURL)
		assert.Equal(t, overrideMappings[ShipyardController], internal.ProjectsV1().(*ProjectHandler).BaseURL)
	})

}
