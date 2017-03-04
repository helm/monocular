package handlers

import (
	"github.com/helm/monocular/src/api/swagger/models"
)

// DataResourceBody returns an data encapsulated version of a resource
func DataResourceBody(resource *models.Resource) *models.ResourceData {
	return &models.ResourceData{
		Data: resource,
	}
}

// DataResourcesBody returns an data encapsulated version of an array of resources
func DataResourcesBody(resources []*models.Resource) *models.ResourceArrayData {
	return &models.ResourceArrayData{
		Data: resources,
	}
}
