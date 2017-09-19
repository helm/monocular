package redis

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/albrow/zoom"

	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

const defaultHost = "localhost:6379"

var (
	pool           *zoom.Pool
	redisPoolOnce  sync.Once
	getReposOnce   sync.Once
	reposSingleton *zoom.Collection
)

type Driver struct {
	Host string
}

func New(host string) (*Driver, error) {
	if host == "" {
		host = defaultHost
	}
	return &Driver{host}, nil
}

// getRedisPool returns a pool of Zoom connections
func (d *Driver) getRedisPool() *zoom.Pool {
	redisPoolOnce.Do(func() {
		pool = d.newRedisPool()
	})
	return pool
}

func (d *Driver) newRedisPool() *zoom.Pool {
	return zoom.NewPool(d.Host)
}

// getReposCollection returns the Repos Zoom collection
func (d *Driver) getReposCollection() (*zoom.Collection, error) {
	var err error
	getReposOnce.Do(func() {
		reposSingleton, err = d.getRedisPool().NewCollectionWithOptions(&data.Repo{}, zoom.DefaultCollectionOptions.WithIndex(true))
	})
	return reposSingleton, err
}

func (d *Driver) GetRepo(name string) (*data.Repo, bool, error) {
	reposCollection, err := d.getReposCollection()
	if err != nil {
		return nil, false, err
	}
	repo := data.Repo{}
	err = reposCollection.Find(name, &repo)

	log.Info("GETTING REPO")
	if err != nil {
		log.Info("ERR GETTING REPO:", err)
		if _, ok := err.(zoom.ModelNotFoundError); ok {
			log.Info("MODEL NOT FOUND ERR")
			return nil, false, nil
		}
		log.Info("OTHER ERR")
		return nil, false, err
	}

	return &repo, true, err
}

func (d *Driver) GetRepos() ([]*data.Repo, error) {
	reposCollection, err := d.getReposCollection()
	if err != nil {
		return nil, err
	}
	var repos []*data.Repo
	// TODO osoriano - should be just repos?
	err = reposCollection.FindAll(&repos)
	return repos, err
}

func (d *Driver) DeleteRepos() (int64, error) {
	reposCollection, err := d.getReposCollection()
	if err != nil {
		return 0, err
	}
	numDeleted, err := reposCollection.DeleteAll()
	return int64(numDeleted), err
}

func (d *Driver) DeleteRepo(name string) (bool, error) {
	reposCollection, err := d.getReposCollection()
	if err != nil {
		return false, err
	}
	return reposCollection.Delete(name)
}

func (d *Driver) CreateRepo(repo *data.Repo) error {
	reposCollection, err := d.getReposCollection()
	if err != nil {
		return err
	}
	return reposCollection.Save(repo)
}

// MergeRepos takes an array of Repos to save in the cache
func (d *Driver) MergeRepos(repos []models.Repo) error {
	reposCollection, err := d.getReposCollection()
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
