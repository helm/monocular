package cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/cache/charthelper"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/data/repos"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/swagger/restapi/operations"
)

type cachedCharts struct {
	// knownRepos is a slice of maps, each of which looks like this: "reponame": "https://repo.url/"
	knownRepos repos.Repos
	allCharts  map[string][]*models.ChartPackage
	rwm        *sync.RWMutex
}

// NewCachedCharts returns a new data.Charts implementation
func NewCachedCharts(repos repos.Repos) data.Charts {
	return &cachedCharts{
		knownRepos: repos,
		rwm:        new(sync.RWMutex),
		allCharts:  make(map[string][]*models.ChartPackage),
	}
}

// ChartFromRepo is the interface implementation for data.Charts
// It returns the reference to a single versioned chart (the most recently published version)
func (c *cachedCharts) ChartFromRepo(repo, name string) (*models.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	if c.allCharts[repo] != nil {
		chart, err := helpers.GetLatestChartVersion(c.allCharts[repo], name)
		if err != nil {
			return nil, err
		}
		return chart, nil
	}
	return nil, fmt.Errorf("no charts found for repo %s\n", repo)
}

// ChartVersionFromRepo is the interface implementation for data.Charts
// It returns the reference to a single versioned chart
func (c *cachedCharts) ChartVersionFromRepo(repo, name, version string) (*models.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	if c.allCharts[repo] != nil {
		chart, err := helpers.GetChartVersion(c.allCharts[repo], name, version)
		if err != nil {
			return nil, err
		}
		return chart, nil
	}
	return nil, fmt.Errorf("no charts found for repo %s\n", repo)
}

// ChartVersionsFromRepo is the interface implementation for data.Charts
// It returns the reference to a slice of all versions of a particular chart in a repo
func (c *cachedCharts) ChartVersionsFromRepo(repo, name string) ([]*models.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	if c.allCharts[repo] != nil {
		charts, err := helpers.GetChartVersions(c.allCharts[repo], name)
		if err != nil {
			return nil, err
		}
		return charts, nil
	}
	return nil, fmt.Errorf("no charts found for repo %s\n", repo)
}

// AllFromRepo is the interface implementation for data.Charts
// It returns the reference to a slice of all versions of all charts in a repo
func (c *cachedCharts) AllFromRepo(repo string) ([]*models.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	if c.allCharts[repo] != nil {
		return c.allCharts[repo], nil
	}
	return nil, fmt.Errorf("no charts found for repo %s\n", repo)
}

// All is the interface implementation for data.Charts
// It returns the reference to a slice of all versions of all charts in all repos
func (c *cachedCharts) All() ([]*models.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	var allCharts []*models.ChartPackage
	// TODO: parallellize this, it won't scale well with lots of repos
	for _, repo := range c.knownRepos {
		var charts []*models.ChartPackage
		for _, chart := range c.allCharts[repo.Name] {
			charts = append(charts, chart)
		}
		allCharts = append(allCharts, charts...)
	}
	return allCharts, nil
}

func (c *cachedCharts) Search(params operations.SearchChartsParams) ([]*models.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	var ret []*models.ChartPackage
	charts, err := c.All()
	if err != nil {
		return nil, err
	}
	for _, chart := range charts {
		if strings.Contains(*chart.Name, params.Name) {
			ret = append(ret, chart)
		}
	}
	return ret, nil
}

// Refresh is the interface implementation for data.Charts
// It refreshes cached data for all authoritative repository+chart data
func (c *cachedCharts) Refresh() error {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	fmt.Printf("Using cache directory %s\n", charthelper.DataDirBase())
	for _, repo := range c.knownRepos {
		// Append index.yaml
		u, _ := url.Parse(repo.URL)
		u.Path = path.Join(u.Path, "index.yaml")

		resp, err := http.Get(u.String())
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		charts, err := helpers.ParseYAMLRepo(body, repo.Name)
		if err != nil {
			return err
		}
		c.allCharts[repo.Name] = []*models.ChartPackage{}
		for _, chart := range charts {
			// Extra files. Skipped if the directory exists
			dataExists, err := charthelper.ChartDataExists(chart)
			if err != nil {
				return err
			}
			if !dataExists {
				fmt.Printf("Local cache missing for %s-%s\n", *chart.Name, *chart.Version)

				err := charthelper.DownloadAndExtractChartTarball(chart)
				if err != nil {
					fmt.Printf("Error on DownloadAndExtractChartTarball: %v\n", err)
					// Skip chart if error extracting the tarball
					continue
				}
				// If we have a problem processing an image it will fallback to the default one
				charthelper.DownloadAndProcessChartIcon(chart)
			}
			c.allCharts[repo.Name] = append(c.allCharts[repo.Name], chart)
		}
	}
	return nil
}
