package httputils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadFromURL(t *testing.T) {

	t.Run("File Download", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(":-)"))
		}))
		content, err := DownloadFromURL(server.URL)
		assert.Nil(t, err)
		assert.Equal(t, content, []byte(":-)"))
	})

}

func TestIsValidURL(t *testing.T) {

	testURLS := map[string]bool{
		"keptn.sh":                    false,
		"keptn..sh":                   false,
		"":                            false,
		"lakjglakbgjoejgfrlej":        false,
		"1":                           false,
		"http://keptn.sh":             true,
		"http://www.keptn.sh":         true,
		"http://keptn.sh/a/b/c":       true,
		"http://keptn.sh/a/b?c=d&e=f": true,
		"http://127.0.0.1/":           true,
	}

	t.Parallel()

	for k, v := range testURLS {
		res := IsValidURL(k)
		assert.Equal(t, v, res, "Value mismatch.\nExpected: %v\nActual: %v", res, v)
	}
}
