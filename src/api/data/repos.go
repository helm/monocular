package data

import (
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// Repo is a Zoom Model for storing repositories
type Repo models.Repo

// ModelId returns the unique name of the Repo
func (r *Repo) ModelId() string {
	if r.Name == nil {
		return "<nil>"
	}
	return *r.Name
}

// SetModelId sets the unique name of the Repo
func (r *Repo) SetModelId(name string) {
	r.Name = &name
}
