package api

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	http2 "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_Proxy(t *testing.T) {

	//data := url.Values{}
	//data.Set("name", "foo")
	//data.Set("surname", "bar")
	//request, err := http.NewRequest(http.MethodGet, "http://localhost:8000?abc=2", strings.NewReader(data.Encode()))
	//if err != nil {
	//	return
	//}
	//
	//h := createProxyHandler(ProxyHost{Host: "localhost:8000", Scheme: "http"}, &http.Client{})
	//h.Proxy(&http2.TestResponseWriter{}, request)
}

func Test_ProxyInternal(t *testing.T) {
	type input struct {
		url               string
		apiMappings       InClusterAPIMappings
		pathProxyMappings map[string]string
	}
	type want struct {
		url string
	}

	tests := []struct {
		input input
		want  want
	}{
		{
			input: input{
				url: "http://localhost:8080/configuration-service/v1/something",
				apiMappings: InClusterAPIMappings{
					ShipyardController:   "shipyard-controller:8080",
					ConfigurationService: "configuration-service:8080",
				},
				pathProxyMappings: map[string]string{
					"/configuration-service": "configuration-service:8080",
					"/controlPlane":          "shipyard-controller:8080",
				},
			},
			want: want{
				url: "http://configuration-service:8080/v1/something",
			},
		},
		{
			input: input{
				url: "http://localhost:8080/controlPlane/v1/something",
				apiMappings: InClusterAPIMappings{
					ShipyardController:   "shipyard-controller:8080",
					ConfigurationService: "configuration-service:8080",
				},
				pathProxyMappings: map[string]string{
					"/configuration-service": "configuration-service:8080",
					"/controlPlane":          "shipyard-controller:8080",
				},
			},
			want: want{
				url: "http://shipyard-controller:8080/v1/something",
			},
		},
	}
	for _, tc := range tests {
		client := newTestClient(func(req *http.Request) *http.Response {
			fmt.Println(req.URL.String())
			assert.Equal(t, tc.want.url, req.URL.String())
			return &http.Response{}
		})
		req, _ := http.NewRequest(http.MethodGet, tc.input.url, strings.NewReader(url.Values{}.Encode()))
		internal, err := NewInternal(tc.input.apiMappings, tc.input.pathProxyMappings, client)
		require.Nil(t, err)
		internal.ProxyV1().Proxy(&http2.TestResponseWriter{}, req)
	}
}

func TestInternalAPISet(t *testing.T) {
	apiMappings := InClusterAPIMappings{
		ShipyardController:   "shipyard-controller:8080",
		ConfigurationService: "configuration-service:8080",
		SecretService:        "secret-service:8080",
	}
	pathProxyMappings := map[string]string{
		"/configuration-service": "configuration-service:8080",
		"/controlPlane":          "shipyard-controller:8080",
	}

	t.Run("", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			assert.Equal(t, "http://secret-service:8080/v1/secret", req.URL.String())
			return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`OK`))}
		})
		internal, err := NewInternal(apiMappings, pathProxyMappings, client)
		require.Nil(t, err)
		internal.SecretsV1().GetSecrets()
	})
	t.Run("", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			assert.Equal(t, "http://shipyard-controller:8080/v1/secret", req.URL.String())
			return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`OK`))}
		})
		internal, err := NewInternal(apiMappings, pathProxyMappings, client)
		require.Nil(t, err)
		internal.ShipyardControlV1().GetOpenTriggeredEvents(EventFilter{})
	})
	t.Run("", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			assert.Equal(t, "http://configuration-service:8080/v1/secret", req.URL.String())
			return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`OK`))}
		})
		internal, err := NewInternal(apiMappings, pathProxyMappings, client)
		require.Nil(t, err)
		internal.ResourcesV1().CreateResources("", "", "", nil)
	})

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
