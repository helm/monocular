package cache

import (
	log "github.com/Sirupsen/logrus"
	"github.com/albrow/zoom"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// Repos is a Zoom Collection for the Repo model
var Repos *zoom.Collection

// NewCachedRepos returns a data.Repos object to manage repositories
func NewCachedRepos(repos []models.Repo) {
	log.Info("setting up Repos collection")
	var err error
	Repos, err = config.Pool.NewCollectionWithOptions(&data.Repo{}, zoom.DefaultCollectionOptions.WithIndex(true))
	if err != nil {
		log.Fatal("unable to create new Repo collection: ", err)
	}
	for _, r := range repos {
		// Convert to Zoom model
		repo := data.Repo(r)
		err = Repos.Save(&repo)
		if err != nil {
			log.Error("unable to save Repo: ", err)
		}
	}
}
