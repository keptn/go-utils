package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPIHandler_getAPIServicePath(t *testing.T) {
	type fields struct {
		BaseURL string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "remove controlPlane path suffix",
			fields: fields{
				BaseURL: "my-api.sh/api/controlPlane",
			},
			want: "my-api.sh/api",
		},
		{
			name: "don't modify anything for internal API calls",
			fields: fields{
				BaseURL: "api-service",
			},
			want: "api-service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				BaseURL: tt.fields.BaseURL,
			}
			assert.Equalf(t, tt.want, a.getAPIServicePath(), "getAPIServicePath()")
		})
	}
}
