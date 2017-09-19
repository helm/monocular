package storage

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/storage/mysql"
	"github.com/kubernetes-helm/monocular/src/api/storage/redis"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

type Storage interface {
	MergeRepos(repos []models.Repo) error
	GetRepo(name string) (*data.Repo, bool, error)
	GetRepos() ([]*data.Repo, error)
	DeleteRepos() (int64, error)
	DeleteRepo(name string) (bool, error)
	CreateRepo(repo *data.Repo) error
}

var Driver Storage

func Init(storageConfig config.StorageConfig) error {
	var err error
	switch storageConfig.Driver {
	case "redis":
		Driver, err = redis.New(storageConfig.Host)
	case "mysql":
		Driver, err = mysql.New(storageConfig.Host)
	case "":
		log.Info("No storage driver specified. Defaulting to Redis")
		Driver, err = redis.New(storageConfig.Host)
	default:
		return errors.New(fmt.Sprintf("Invalid storage.Driver: %s", storageConfig.Driver))
	}
	return err
}
