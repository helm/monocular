package storage

import (
	"errors"
	"fmt"

	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/storage/redis"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

type Storage interface {
	MergeRepos(repos []models.Repo) error
	GetRepo(name string) (*data.Repo, bool, error)
	GetRepos() ([]*data.Repo, error)
	DeleteRepos() (int, error)
	DeleteRepo(name string) (bool, error)
	CreateRepo(repo *data.Repo) error
}

var Driver Storage

func Init(storageConfig config.StorageConfig) error {
	switch storageConfig.Driver {
	case "redis":
		Driver = redis.New(storageConfig.Host)
	//case "mysql":
	//	Driver = mysql.NewDriver()
	case "":
		// TODO add log statement... No storage driver specified. Defaulting to Redis
		Driver = redis.New(storageConfig.Host)
	default:
		return errors.New(fmt.Sprintf("Invalid storage.Driver: %s", storageConfig.Driver))
	}
	return nil
}
