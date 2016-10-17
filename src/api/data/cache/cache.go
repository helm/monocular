package cache

import (
	"fmt"
	"sync"

	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/mocks"
	"github.com/helm/monocular/src/api/swagger/models"
)

var mockCharts = mocks.NewMockCharts()

type cachedCharts struct {
	knownRepos []map[string]string
	allCharts  map[string][]*models.ChartVersion
	rwm        *sync.RWMutex
}

// NewCachedCharts returns a new data.Charts implementation
func NewCachedCharts(repos []map[string]string) data.Charts {
	return &cachedCharts{
		knownRepos: repos,
		rwm:        new(sync.RWMutex),
		allCharts:  make(map[string][]*models.ChartVersion),
	}
}

// ChartFromRepo is the interface implementation for data.Charts
// It returns the reference to a single versioned chart (the most recently published version)
func (c *cachedCharts) ChartFromRepo(repo, name string) (*models.ChartVersion, error) {
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

// AllFromRepo is the interface implementation for data.Charts
// It returns the reference to a slice of all versions of all charts in a repo
func (c *cachedCharts) AllFromRepo(repo string) ([]*models.ChartVersion, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	if c.allCharts[repo] != nil {
		return c.allCharts[repo], nil
	}
	return nil, fmt.Errorf("no charts found for repo %s\n", repo)
}

// All is the interface implementation for data.Charts
// It returns the reference to a slice of all versions of all charts in all repos
func (c *cachedCharts) All() ([]*models.Resource, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	var allCharts []*models.Resource
	// TODO: parallellize this, it won't scale well with lots of repos
	for _, repo := range c.knownRepos {
		for repoName := range repo {
			for _, chart := range c.allCharts[repoName] {
				resource := helpers.MakeChartResource(chart, repoName)
				allCharts = append(allCharts, resource)
			}
		}
	}
	return allCharts, nil
}

// Refresh is the interface implementation for data.Charts
// It refreshes cached data for all authoritative repository+chart data
func (c *cachedCharts) Refresh() error {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	for _, repo := range c.knownRepos {
		for repoName := range repo {
			charts, err := mockCharts.AllFromRepo(repoName)
			if err != nil {
				return err
			}
			c.allCharts[repoName] = []*models.ChartVersion{}
			for _, chart := range charts {
				c.allCharts[repoName] = append(c.allCharts[repoName], chart)
			}
		}
	}
	return nil
}
