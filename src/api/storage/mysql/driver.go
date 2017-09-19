package mysql

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"

	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

const defaultHost = "monocular:monocular@tcp(localhost:3306)/monocular"

const mergeQuery = `insert into repo values(?, ?, ?)
on duplicate key update url=?, source=?`

type Driver struct {
	db *sql.DB
}

func New(host string) (*Driver, error) {
	if host == "" {
		host = defaultHost
	}
	db, err := sql.Open("mysql", host)
	if err != nil {
		return nil, err
	}
	return &Driver{db}, nil
}

func (d *Driver) GetRepo(name string) (*data.Repo, bool, error) {
	row := d.db.QueryRow("select * from repo where name=?", name)
	var existingName string
	var existingUrl string
	var existingSource string
	err := row.Scan(&existingName, &existingUrl, &existingSource)
	switch {
	case err == sql.ErrNoRows:
		return nil, false, nil
	case err != nil:
		return nil, false, err
	default:
		repo := data.Repo{&existingUrl, &existingName, existingSource}
		return &repo, true, nil
	}
}

func (d *Driver) GetRepos() ([]*data.Repo, error) {
	log.Info("GETTING REPOS MYSQL")
	rows, err := d.db.Query("select * from repo")
	if err != nil {
		log.Info("MYSQL ERR GETTING REPOS")
		return nil, err
	}
	defer rows.Close()

	var repos []*data.Repo
	for rows.Next() {
		log.Info("IN REPO ITER")
		var name string
		var url string
		var source string
		err = rows.Scan(&name, &url, &source)
		if err != nil {
			return nil, err
		}
		repo := data.Repo{&url, &name, source}
		repos = append(repos, &repo)
	}
	log.Info("DONE REPO ITER")
	return repos, rows.Err()
}

func (d *Driver) DeleteRepos() (int64, error) {
	res, err := d.db.Exec("delete from repo")
	if err != nil {
		return 0, err
	}
	numAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return numAffected, nil
}

func (d *Driver) DeleteRepo(name string) (bool, error) {
	_, found, err := d.GetRepo(name)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	_, err = d.db.Exec("delete from repo where name=?", name)
	return true, err
}

func (d *Driver) CreateRepo(repo *data.Repo) error {
	lowerLevelRepo := models.Repo{repo.URL, repo.Name, repo.Source}
	return d.mergeRepo(lowerLevelRepo)
}

func (d *Driver) mergeRepo(repo models.Repo) error {
	_, err := d.db.Exec(mergeQuery, *repo.Name, *repo.URL,
		repo.Source, *repo.URL, repo.Source)
	return err
}

// MergeRepos takes an array of Repos to save in the cache
func (d *Driver) MergeRepos(repos []models.Repo) error {
	for _, repo := range repos {
		err := d.mergeRepo(repo)
		if err != nil {
			return err
		}
	}
	return nil
}
