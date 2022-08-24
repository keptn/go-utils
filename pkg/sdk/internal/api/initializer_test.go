package api

import (
	"fmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnapiv2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	"github.com/keptn/go-utils/pkg/sdk/connector/logger"
	"github.com/keptn/go-utils/pkg/sdk/internal/config"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type fakeHTTPClientFactory struct {
	GetFn func() (*http.Client, error)
}

func (f *fakeHTTPClientFactory) Get() (*http.Client, error) {
	return f.GetFn()
}

func Test_Initialize(t *testing.T) {
	t.Run("Remote use case - invalid keptn api endpoint", func(t *testing.T) {
		env := config.EnvConfig{KeptnAPIEndpoint: "://mynotsogoodendpoint"}
		result, err := Initialize(env, CreateClientGetter(env), logger.NewDefaultLogger())
		require.Error(t, err)
		require.Nil(t, result)
	})
	t.Run("Remote use case - no http address as keptn api endpoint", func(t *testing.T) {
		env := config.EnvConfig{KeptnAPIEndpoint: "ssh://mynotsogoodendpoint"}
		result, err := Initialize(env, CreateClientGetter(env), logger.NewDefaultLogger())
		require.Error(t, err)
		require.Nil(t, result)

	})
	t.Run("Remote use case - remote api set is used", func(t *testing.T) {
		env := config.EnvConfig{KeptnAPIEndpoint: "http://endpoint"}
		result, err := Initialize(env, CreateClientGetter(env), logger.NewDefaultLogger())
		require.NoError(t, err)
		require.NotNil(t, result.ControlPlane)
		require.NotNil(t, result.EventSenderCallback)
		require.NotNil(t, result.KeptnAPI)
		require.IsType(t, &keptnapi.APISet{}, result.KeptnAPI)
		require.NotNil(t, result.KeptnAPIV2)
		require.IsType(t, &keptnapiv2.APISet{}, result.KeptnAPIV2)
	})
	t.Run("Internal Use case - internal api set is used", func(t *testing.T) {
		env := config.EnvConfig{}
		result, err := Initialize(env, CreateClientGetter(env), logger.NewDefaultLogger())
		require.NoError(t, err)
		require.NotNil(t, result.ControlPlane)
		require.NotNil(t, result.EventSenderCallback)
		require.NotNil(t, result.KeptnAPI)
		require.IsType(t, &keptnapi.InternalAPISet{}, result.KeptnAPI)
		require.NotNil(t, result.KeptnAPIV2)
		require.IsType(t, &keptnapiv2.InternalAPISet{}, result.KeptnAPIV2)
	})
	t.Run("HTTP client creation fails", func(t *testing.T) {
		env := config.EnvConfig{KeptnAPIEndpoint: "http://endpoint"}
		result, err := Initialize(env, &fakeHTTPClientFactory{GetFn: func() (*http.Client, error) { return nil, fmt.Errorf("err") }}, logger.NewDefaultLogger())
		require.Error(t, err)
		require.Nil(t, result)
	})
	t.Run("HTTP client creation returns nil client", func(t *testing.T) {
		env := config.EnvConfig{KeptnAPIEndpoint: "http://endpoint"}
		result, err := Initialize(env, &fakeHTTPClientFactory{GetFn: func() (*http.Client, error) { return nil, nil }}, logger.NewDefaultLogger())
		require.NoError(t, err)
		require.NotNil(t, result.ControlPlane)
		require.NotNil(t, result.EventSenderCallback)
		require.NotNil(t, result.KeptnAPI)
		require.IsType(t, &keptnapi.APISet{}, result.KeptnAPI)
		require.NotNil(t, result.KeptnAPIV2)
		require.IsType(t, &keptnapiv2.APISet{}, result.KeptnAPIV2)
	})
}
