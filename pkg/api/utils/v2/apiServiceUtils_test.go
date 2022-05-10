package v2

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Test_createInstrumentedClientTransport(t *testing.T) {
	client := createInstrumentedClientTransport(nil)
	assert.NotNil(t, client)
	assert.NotNil(t, client)
	_, isOtelTransport := client.Transport.(*otelhttp.Transport)
	assert.True(t, isOtelTransport)

	client = createInstrumentedClientTransport(&http.Client{})
	assert.NotNil(t, client)
	_, isOtelTransport = client.Transport.(*otelhttp.Transport)
	assert.True(t, isOtelTransport)

	client = createInstrumentedClientTransport(&http.Client{Transport: &http.Transport{}})
	assert.NotNil(t, client)
	_, isOtelTransport = client.Transport.(*otelhttp.Transport)
	assert.True(t, isOtelTransport)
}
