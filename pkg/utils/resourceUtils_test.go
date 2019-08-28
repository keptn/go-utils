package utils

import (
	"testing"

	"github.com/keptn/go-utils/pkg/models"
)

func TestResourceHandler_CreateProjectResources(t *testing.T) {
	type fields struct {
		BaseURL string
	}
	type args struct {
		project   string
		resources []*models.Resource
	}
	uri := "shipyard-hello2.yaml"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "create resource",
			fields: fields{
				BaseURL: "localhost:8080",
			},
			args: args{
				project: "rockshop",
				resources: []*models.Resource{
					&models.Resource{
						ResourceURI:     &uri,
						ResourceContent: "foo",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceHandler{
				BaseURL: tt.fields.BaseURL,
			}
			got, err := r.CreateProjectResources(tt.args.project, tt.args.resources)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceHandler.CreateProjectResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) <= 0 {
				t.Errorf("Got empty response")
			}
		})
	}
}

func TestResourceHandler_UpdateProjectResource(t *testing.T) {
	type fields struct {
		BaseURL string
	}
	type args struct {
		project  string
		resource *models.Resource
	}
	uri := "shipyard-tests.yaml"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "update resource",
			fields: fields{
				BaseURL: "localhost:8080",
			},
			wantErr: false,
			args: args{
				project: "rockshop",
				resource: &models.Resource{
					ResourceURI:     &uri,
					ResourceContent: "this is a test!",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceHandler{
				BaseURL: tt.fields.BaseURL,
			}
			got, err := r.UpdateProjectResource(tt.args.project, tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceHandler.UpdateProjectResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResourceHandler.UpdateProjectResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
