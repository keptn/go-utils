package utils

import (
	"github.com/keptn/go-utils/pkg/configuration-service/utils"
	"github.com/keptn/go-utils/pkg/models"
)

// KeptnHandler provides an interface to keptn resources
type KeptnHandler struct {
	ResourceHandler *utils.ResourceHandler
}

// NewKeptnHandler returns a new KeptnHandler instance
func NewKeptnHandler(rh *utils.ResourceHandler) *KeptnHandler {
	return &KeptnHandler{
		ResourceHandler: rh,
	}
}

// GetShipyard returns the shipyard definition of a project
func (k *KeptnHandler) GetShipyard(project string) (*models.Shipyard, error) {
	shipyardResource, err := k.ResourceHandler.GetProjectResource(project, "shipyard.yaml")
	if err != nil {
		return nil, err
	}

	shipyard := models.Shipyard{}
	err = yaml.Unmarshal([]byte(shipyardResource.ResourceContent), &shipyard)
	if err != nil {
		return nil, err
	}
	return &shipyard, nil
}
