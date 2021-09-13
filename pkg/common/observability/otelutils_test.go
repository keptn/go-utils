package observability

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

type tpTestCase struct {
	name                  string
	collectorEndpoint     string
	wantTraceProviderType string
}

func TestCreateTraceProvider(t *testing.T) {

	tests := []tpTestCase{
		{
			name:                  "empty collector endpoint",
			collectorEndpoint:     "",
			wantTraceProviderType: "trace.noopTracerProvider",
		},
		{
			name:                  "invalid-url",
			collectorEndpoint:     "some-url",
			wantTraceProviderType: "trace.noopTracerProvider",
		},
		{
			name:                  "k8s dns service name",
			collectorEndpoint:     "otel-collector.observability:4317",
			wantTraceProviderType: "*trace.TracerProvider",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(EnvVarOTelCollectorEndpoint, tc.collectorEndpoint)

			shutdown := InitOTelTraceProvider("my-service")
			tp := otel.GetTracerProvider()

			assert.NotNil(t, tp)
			assert.Equal(t, tc.wantTraceProviderType, reflect.TypeOf(tp).String())

			shutdown()
			os.Setenv(EnvVarOTelCollectorEndpoint, "")
		})
	}
}
