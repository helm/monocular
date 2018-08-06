package repohelper

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/models"
	swaggermodels "github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/version"
)

// GetRepoIndexFile Get the charts from the index file
var GetChartsFromRepoIndexFile = func(repo *models.Repo) ([]*swaggermodels.ChartPackage, error) {
	u, _ := url.Parse(repo.URL)
	u.Path = path.Join(u.Path, "index.yaml")

	// 1 - Download repo index
	var client http.Client
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", version.GetUserAgent())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 2 - Parse repo index
	charts, err := helpers.ParseYAMLRepo(body, repo.Name)
	if err != nil {
		return nil, err
	}

	return charts, nil
}
