package data

import (
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// Extend the repos model to satisfy the zoom model interface
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
