package cache

import (
	"sync"

	"github.com/albrow/zoom"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

var (
	reposSingleton *zoom.Collection
	once           sync.Once
)

// UpdateCache takes an array of Repos to save in the cache
func UpdateCache(repos []models.Repo) error {
	reposCollection, err := GetRepos()
	if err != nil {
		return err
	}
	for _, r := range repos {
		// Convert to Zoom model
		repo := data.Repo(r)
		err = reposCollection.Save(&repo)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetRepos returns the Repos Zoom collection
func GetRepos() (*zoom.Collection, error) {
	var err error
	once.Do(func() {
		reposSingleton, err = config.GetRedisPool().NewCollectionWithOptions(&data.Repo{}, zoom.DefaultCollectionOptions.WithIndex(true))
	})
	return reposSingleton, err
}
