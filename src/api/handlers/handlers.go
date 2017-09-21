package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// Params a key-value map of path params
type Params map[string]string

// WithParams can be used to wrap handlers to take an extra arg for path params
type WithParams func(http.ResponseWriter, *http.Request, Params)

func (h WithParams) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	h(w, req, vars)
}

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
