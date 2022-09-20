package keptn

import (
	"context"
	"errors"
	"fmt"
	keptnapi "github.com/keptn/go-utils/pkg/lib/keptn"
	"gopkg.in/yaml.v3"
)

// SLI struct for SLI.yaml
type SLI struct {
	SpecVersion string            `yaml:"spec_version"`
	Indicators  map[string]string `yaml:"indicators"`
}

// GetSLIOptions specifies the project, stage, dervice and SLI file name to be fetched
type GetSLIOptions struct {
	Project     string
	Stage       string
	Service     string
	SLIFileName string
}

// SLIHelper is the default implementation of SLIReader used for retrieving SLI content
type SLIHelper struct {
	client ResourceClientInterface
}

// NewSLIHelper creates a new SLIHelper with a Keptn resource handler for the configuration service.
func NewSLIHelper(client ResourceClientInterface) *SLIHelper {
	return &SLIHelper{
		client: client,
	}
}

type sliMap map[string]string

func (m sliMap) insertOrUpdateMany(x map[string]string) {
	for key, value := range x {
		m[key] = value
	}
}

// GetSLIs gets the SLIs stored for the specified project, stage and service.
// First, the configuration of project-level is retrieved, which is then overridden by configuration on stage level, and then overridden by configuration on service level.
func (rc *SLIHelper) GetSLIs(ctx context.Context, options GetSLIOptions) (map[string]string, error) {
	slis := make(sliMap)

	// try to get SLI config from project
	if options.Project != "" {
		projectSLIs, err := getSLIsFromResource(func() (string, error) { return rc.client.GetProjectResource(ctx, options.Project, options.SLIFileName) })
		if err != nil {
			return nil, err
		}

		slis.insertOrUpdateMany(projectSLIs)
	}

	// try to get SLI config from stage
	if options.Project != "" && options.Stage != "" {
		stageSLIs, err := getSLIsFromResource(func() (string, error) {
			return rc.client.GetStageResource(ctx, options.Project, options.Stage, options.SLIFileName)
		})
		if err != nil {
			return nil, err
		}

		slis.insertOrUpdateMany(stageSLIs)
	}

	// try to get SLI config from service
	if options.Project != "" && options.Stage != "" && options.Service != "" {
		serviceSLIs, err := getSLIsFromResource(func() (string, error) {
			return rc.client.GetServiceResource(ctx, options.Project, options.Stage, options.Service, options.SLIFileName)
		})
		if err != nil {
			return nil, err
		}

		slis.insertOrUpdateMany(serviceSLIs)
	}

	return slis, nil
}

type resourceGetterFunc func() (string, error)

// getSLIsFromResource uses the specified function to get a resource and returns the SLIs as a map.
// If is is not possible to get the resource for any other reason than it is not found, or it is not possible to unmarshal the file or it doesn't contain any indicators, an error is returned.
func getSLIsFromResource(resourceGetter resourceGetterFunc) (map[string]string, error) {
	resource, err := resourceGetter()
	if err != nil {
		var rnfErrorType *ResourceNotFoundError
		if errors.As(err, &rnfErrorType) {
			return nil, nil
		}

		return nil, err
	}

	return readSLIsFromResource(resource)
}

// readSLIsFromResource unmarshals a resource as a SLIConfig and returns the SLIs as a map.
// If it is not possible to unmarshal the file or it doesn't contain any indicators, an error is returned.
func readSLIsFromResource(resource string) (map[string]string, error) {
	sliConfig := keptnapi.SLIConfig{}
	err := yaml.Unmarshal([]byte(resource), &sliConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to unrmarshal sli content: %v", err)
	}

	if len(sliConfig.Indicators) == 0 {
		return nil, errors.New("missing required field: indicators")
	}

	return sliConfig.Indicators, nil
}
