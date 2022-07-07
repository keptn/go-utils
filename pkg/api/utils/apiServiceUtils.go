package api

import (
	"crypto/tls"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// APIService represents the interface for accessing the configuration service
type APIService interface {
	getBaseURL() string
	getAuthToken() string
	getAuthHeader() string
	getHTTPClient() *http.Client
}

// createInstrumentedClientTransport tries to add support for opentelemetry
// to the given http.Client. If httpClient is nil, a fresh http.Client
// with opentelemetry support is created
func createInstrumentedClientTransport(httpClient *http.Client) *http.Client {
	if httpClient == nil {
		return &http.Client{
			Transport: wrapOtelTransport(getClientTransport(nil)),
		}
	}
	httpClient.Transport = wrapOtelTransport(getClientTransport(httpClient.Transport))
	return httpClient
}

// Wraps the provided http.RoundTripper with one that
// starts a span and injects the span context into the outbound request headers.
func wrapOtelTransport(base http.RoundTripper) *otelhttp.Transport {
	return otelhttp.NewTransport(base)
}

// getClientTransport returns a client transport which
// skips verifying server certificates and is able to
// read proxy configuration from environment variables
//
// If the given http.RoundTripper is nil then a new http.Transport
// is created, otherwise the given http.RoundTripper is analysed whether it
// is of type *http.Transport. If so, the respective settings for
// disabling server certificate verification as well as proxy server support are set
// If not, the given http.RoundTripper is passed through untouched
func getClientTransport(rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
		}
		return tr
	}
	if tr, isDefaultTransport := rt.(*http.Transport); isDefaultTransport {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		tr.Proxy = http.ProxyFromEnvironment
		return tr
	}
	return rt
}
