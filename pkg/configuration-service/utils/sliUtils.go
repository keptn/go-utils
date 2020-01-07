package utils

import (
	"strings"

	"github.com/keptn/go-utils/pkg/configuration-service/models"
	"gopkg.in/yaml.v2"
)

// SLIConfig represents the struct of a SLI file
type SLIConfig struct {
	Indicators map[string]string `json:"indicators" yaml:"indicators"`
}

// GetSLIConfiguration retrieves the SLI configuration for a service considering SLI configuration on stage and project level.
// First, the configuration of project-level is retrieved, which is then overridden by configuration on stage level,
// overridden by configuration on service level.
func (r *ResourceHandler) GetSLIConfiguration(project string, stage string, service string, resourceURI string) (map[string]string, error) {
	var res *models.Resource
	var err error
	SLIs := make(map[string]string)

	// get sli config from project
	if project != "" {
		res, err = r.GetProjectResource(project, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(err.Error(), "resource not found") {
				return nil, err
			}
		}
		SLIs, err = addResourceContentToSLIMap(SLIs, res)
		if err != nil {
			return nil, err
		}
	}

	// get sli config from stage
	if project != "" && stage != "" {
		res, err = r.GetStageResource(project, stage, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(err.Error(), "resource not found") {
				return nil, err
			}
		}
		SLIs, err = addResourceContentToSLIMap(SLIs, res)
		if err != nil {
			return nil, err
		}
	}

	// get sli config from service
	if project != "" && stage != "" && service != "" {
		res, err = r.GetServiceResource(project, stage, service, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(err.Error(), "resource not found") {
				return nil, err
			}
		}
		SLIs, err = addResourceContentToSLIMap(SLIs, res)
		if err != nil {
			return nil, err
		}
	}

	return SLIs, nil
}

func addResourceContentToSLIMap(SLIs map[string]string, resource *models.Resource) (map[string]string, error) {
	if resource != nil {
		sliConfig := SLIConfig{}
		err := yaml.Unmarshal([]byte(resource.ResourceContent), &sliConfig)
		if err != nil {
			return nil, err
		}

		for key, value := range sliConfig.Indicators {
			SLIs[key] = value
		}
	}
	return SLIs, nil
}