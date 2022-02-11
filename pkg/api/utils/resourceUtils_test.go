package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceHandler_buildResourceURI(t *testing.T) {
	scheme := "https"
	tests := []struct {
		name     string
		project  string
		stage    string
		service  string
		resource string
		want     string
	}{
		{
			name:     "project resources",
			project:  "myproject",
			stage:    "",
			service:  "",
			resource: "",
			want:     scheme + "://" + configurationServiceBaseURL + v1ProjectPath + "/myproject/resource",
		},
		{
			name:     "project resource",
			project:  "myproject",
			stage:    "",
			service:  "",
			resource: "metadata.yaml",
			want:     scheme + "://" + configurationServiceBaseURL + v1ProjectPath + "/myproject/resource/metadata.yaml",
		},
		{
			name:     "stage resource",
			project:  "sockshop",
			stage:    "dev",
			service:  "",
			resource: "metadata.yaml",
			want:     scheme + "://" + configurationServiceBaseURL + v1ProjectPath + "/sockshop" + pathToStage + "/dev/resource/metadata.yaml",
		},
		{
			name:     "service resource",
			project:  "sockshop",
			stage:    "dev",
			service:  "helloservice",
			resource: "hello.go",
			want:     scheme + "://" + configurationServiceBaseURL + v1ProjectPath + "/sockshop" + pathToStage + "/dev/service/helloservice/resource/hello.go",
		},
		{
			name:     "service resources",
			project:  "sockshop",
			stage:    "dev",
			service:  "helloservice",
			resource: "",
			want:     scheme + "://" + configurationServiceBaseURL + v1ProjectPath + "/sockshop" + pathToStage + "/dev/service/helloservice/resource",
		},
		{
			name:     "service resource escape / ",
			project:  "sockshop",
			stage:    "dev",
			service:  "helloservice",
			resource: "helm%2Fhelloservice.tgz",
			want:     scheme + "://" + configurationServiceBaseURL + v1ProjectPath + "/sockshop" + pathToStage + "/dev/service/helloservice/resource/helm%252Fhelloservice.tgz",
		},
	}

	r := &ResourceHandler{
		BaseURL: configurationServiceBaseURL,
		Scheme:  scheme,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := *NewResourceScope().Project(tt.project).Service(tt.service).Resource(tt.resource).Stage(tt.stage)
			assert.Equalf(t, tt.want, r.buildResourceURI(scope), "buildResourceURI(%v)", scope)
		})
	}
}
