package keptn

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	v2 "github.com/keptn/go-utils/pkg/api/utils/v2"
	fake "github.com/keptn/go-utils/pkg/api/utils/v2/fake"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetResource_APIReturnsError(t *testing.T) {
	resourcesHandler := &fake.ResourcesInterfaceMock{
		GetResourceFunc: func(ctx context.Context, scope v2.ResourceScope, opts v2.ResourcesGetResourceOptions) (*models.Resource, error) {
			return nil, fmt.Errorf("error")
		},
	}
	resourceClient := NewResourceHelper(resourcesHandler)
	resource0, err0 := resourceClient.GetResource(context.TODO(), "my-project", "my-stage", "my-service", "test/resource.yaml")
	assert.Equal(t, "", resource0)
	assert.Error(t, err0)

	resource1, err1 := resourceClient.GetResource(context.TODO(), "my-project", "my-stage", "", "test/resource.yaml")
	assert.Equal(t, "", resource1)
	assert.Error(t, err1)

	resource2, err2 := resourceClient.GetResource(context.TODO(), "my-project", "", "", "test/resource.yaml")
	assert.Equal(t, "", resource2)
	assert.Error(t, err2)

}

func TestGetResource_ResourceNotFoundError(t *testing.T) {
	resourcesHandler := &fake.ResourcesInterfaceMock{
		GetResourceFunc: func(ctx context.Context, scope v2.ResourceScope, opts v2.ResourcesGetResourceOptions) (*models.Resource, error) {
			return nil, v2.ResourceNotFoundError
		},
	}
	resourceClient := NewResourceHelper(resourcesHandler)
	resource, err := resourceClient.GetResource(context.TODO(), "my-project", "my-stage", "my-service", "test/resource.yaml")
	assert.Equal(t, "", resource)
	var errExp *ResourceNotFoundError
	assert.ErrorAs(t, err, &errExp)
}

func TestGetResource(t *testing.T) {
	resourcesHandler := &fake.ResourcesInterfaceMock{
		GetResourceFunc: func(ctx context.Context, scope v2.ResourceScope, opts v2.ResourcesGetResourceOptions) (*models.Resource, error) {
			return &models.Resource{
				Metadata: &models.Version{
					Branch:      "branch",
					UpstreamURL: "http://some.url.com",
					Version:     "1.0",
				},
				ResourceContent: "some-content",
				ResourceURI:     strutils.Stringp("test/resource.yaml"),
			}, nil
		},
	}

	//service resource
	resourceClient := NewResourceHelper(resourcesHandler)
	resource0, err0 := resourceClient.GetResource(context.TODO(), "my-project", "my-stage", "my-service", "test/resource.yaml")
	assert.Equal(t, "some-content", resource0)
	assert.NoError(t, err0)

	// stage resource
	resource1, err1 := resourceClient.GetResource(context.TODO(), "my-project", "my-stage", "", "test/resource.yaml")
	assert.Equal(t, "some-content", resource1)
	assert.NoError(t, err1)

	// project resource
	resource2, err2 := resourceClient.GetResource(context.TODO(), "my-project", "", "", "test/resource.yaml")
	assert.Equal(t, "some-content", resource2)
	assert.NoError(t, err2)
}

func TestGetResource_EmptyResourceContent(t *testing.T) {
	resourcesHandler := &fake.ResourcesInterfaceMock{
		GetResourceFunc: func(ctx context.Context, scope v2.ResourceScope, opts v2.ResourcesGetResourceOptions) (*models.Resource, error) {
			return &models.Resource{
				Metadata: &models.Version{
					Branch:      "branch",
					UpstreamURL: "http://some.url.com",
					Version:     "1.0",
				},
				ResourceContent: "", // <-- it's empty :-O
				ResourceURI:     strutils.Stringp("test/resource.yaml"),
			}, nil
		},
	}

	//service resource
	resourceClient := NewResourceHelper(resourcesHandler)
	resource0, err0 := resourceClient.GetResource(context.TODO(), "my-project", "my-stage", "my-service", "test/resource.yaml")
	var errExp0 *ResourceEmptyError
	assert.Equal(t, "", resource0)
	assert.ErrorAs(t, err0, &errExp0)

	// stage resource
	resource1, err1 := resourceClient.GetResource(context.TODO(), "my-project", "my-stage", "", "test/resource.yaml")
	var errExp1 *ResourceEmptyError
	assert.Equal(t, "", resource1)
	assert.ErrorAs(t, err1, &errExp1)

	// project resource
	resource2, err2 := resourceClient.GetResource(context.TODO(), "my-project", "", "", "test/resource.yaml")
	var errExp2 *ResourceEmptyError
	assert.Equal(t, "", resource2)
	assert.ErrorAs(t, err2, &errExp2)
}

func TestUploadResource_Error(t *testing.T) {
	resourcesHandler := &fake.ResourcesInterfaceMock{
		CreateResourcesFunc: func(ctx context.Context, project string, stage string, service string, resources []*models.Resource, opts v2.ResourcesCreateResourcesOptions) (*models.EventContext, *models.Error) {
			return nil, &models.Error{}
		},
	}

	resourceClient := NewResourceHelper(resourcesHandler)
	err := resourceClient.UploadResource(context.TODO(), []byte{}, "test/resource.yaml", "my-project", "my-stage", "my-service")
	var errExp *ResourceUploadFailedError
	assert.ErrorAs(t, err, &errExp)

}

func TestUploadResource(t *testing.T) {
	resourcesHandler := &fake.ResourcesInterfaceMock{
		CreateResourcesFunc: func(ctx context.Context, project string, stage string, service string, resources []*models.Resource, opts v2.ResourcesCreateResourcesOptions) (*models.EventContext, *models.Error) {
			return &models.EventContext{KeptnContext: strutils.Stringp("abcde")}, nil
		},
	}

	resourceClient := NewResourceHelper(resourcesHandler)
	err := resourceClient.UploadResource(context.TODO(), []byte{}, "test/resource.yaml", "my-project", "my-stage", "my-service")
	assert.NoError(t, err)

}
