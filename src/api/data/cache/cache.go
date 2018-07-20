package cache

import (
	"fmt"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/cache/charthelper"
	"github.com/helm/monocular/src/api/data/cache/repohelper"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/datastore"
	"github.com/helm/monocular/src/api/models"
	swaggermodels "github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/swagger/restapi/operations/charts"
)

type cachedCharts struct {
	allCharts map[string][]*swaggermodels.ChartPackage
	rwm       *sync.RWMutex
	dbSession datastore.Session
}

// NewCachedCharts returns a new data.Charts implementation
func NewCachedCharts(dbSession datastore.Session) data.Charts {
	return &cachedCharts{
		rwm:       new(sync.RWMutex),
		allCharts: make(map[string][]*swaggermodels.ChartPackage),
		dbSession: dbSession,
	}
}

// ChartFromRepo is the interface implementation for data.Charts
// It returns the reference to a single versioned chart (the most recently published version)
func (c *cachedCharts) ChartFromRepo(repo, name string) (*swaggermodels.ChartPackage, error) {
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
func (c *cachedCharts) ChartVersionFromRepo(repo, name, version string) (*swaggermodels.ChartPackage, error) {
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
func (c *cachedCharts) ChartVersionsFromRepo(repo, name string) ([]*swaggermodels.ChartPackage, error) {
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
func (c *cachedCharts) AllFromRepo(repo string) ([]*swaggermodels.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	if c.allCharts[repo] != nil {
		return c.allCharts[repo], nil
	}
	return nil, fmt.Errorf("no charts found for repo %s\n", repo)
}

// All is the interface implementation for data.Charts
// It returns the reference to a slice of all versions of all charts in all repos
func (c *cachedCharts) All() ([]*swaggermodels.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	var allCharts []*swaggermodels.ChartPackage

	db, closer := c.dbSession.DB()
	defer closer()

	repos, err := models.ListRepos(db)
	if err != nil {
		return nil, err
	}
	// TODO: parallellize this, it won't scale well with lots of repos
	for _, repo := range repos {
		var charts []*swaggermodels.ChartPackage
		for _, chart := range c.allCharts[repo.Name] {
			charts = append(charts, chart)
		}
		allCharts = append(allCharts, charts...)
	}
	return allCharts, nil
}

func (c *cachedCharts) Search(params charts.SearchChartsParams) ([]*swaggermodels.ChartPackage, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	var ret []*swaggermodels.ChartPackage
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

// deleteChart is the interface implementation for data.Charts
// It deletes Chart from memory
// NOTE: This method does not delete from filesystem for now. It is used to test single chart refresh
func (c *cachedCharts) DeleteChart(repoName string, chartName string, chartVersion string) error {

	db, closer := c.dbSession.DB()
	defer closer()

	repo, err := models.GetRepo(db, repoName)
	if err != nil {
		return err
	}

	c.rwm.Lock()

	for k, chart := range c.allCharts[repoName] {

		if *chart.Name == chartName && *chart.Version == chartVersion {
			log.WithFields(log.Fields{
				"path": charthelper.DataDirBase(),
			}).Info(chartName + " - " + chartVersion + " deleted")

			c.allCharts[repo.Name] = append(c.allCharts[repo.Name][:k], c.allCharts[repo.Name][k+1:]...)
			break
		}
	}

	c.rwm.Unlock()

	return nil
}

// Refresh is the interface implementation for data.Charts
// It refreshes cached data for a specific repo and chart
func (c *cachedCharts) RefreshChart(repoName string, chartName string) error {

	log.WithFields(log.Fields{
		"path": charthelper.DataDirBase(),
	}).Info("Using cache directory")

	db, closer := c.dbSession.DB()
	defer closer()

	repo, err := models.GetRepo(db, repoName)
	if err != nil {
		return err
	}
	charts, err := repohelper.GetChartsFromRepoIndexFile(repo)
	if err != nil {
		return err
	}

	didUpdate := false
	var alreadyExists bool
	for _, chartFromIndex := range charts {
		if *chartFromIndex.Name == chartName {
			didUpdate = true
			ch := make(chan chanItem, len(charts))
			defer close(ch)
			go processChartMetadata(chartFromIndex, repo.URL, ch)

			it := <-ch
			if it.err == nil {
				c.rwm.Lock()
				// find the key
				alreadyExists = false
				for k, chart := range c.allCharts[repo.Name] {
					if *chart.Name == *it.chart.Name && *chart.Version == *it.chart.Version {
						c.allCharts[repo.Name][k] = it.chart
						alreadyExists = true
						break
					}
				}
				if alreadyExists == false {
					c.allCharts[repo.Name] = append(c.allCharts[repo.Name], it.chart)
				}
				c.rwm.Unlock()
			} else {
				return it.err
			}

		}
	}

	if didUpdate == false {
		return fmt.Errorf("no chart \"%s\" found for repo %s\n", chartName, repo.Name)
	} else {
		return nil
	}
}

// Refresh is the interface implementation for data.Charts
// It refreshes cached data for all authoritative repository+chart data
func (c *cachedCharts) Refresh() error {
	// New list of charts that will replace cached charts
	var updatedCharts = make(map[string][]*swaggermodels.ChartPackage)

	log.WithFields(log.Fields{
		"path": charthelper.DataDirBase(),
	}).Info("Using cache directory")

	db, closer := c.dbSession.DB()
	defer closer()

	repos, err := models.ListRepos(db)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		charts, err := repohelper.GetChartsFromRepoIndexFile(repo)
		if err != nil {
			log.WithFields(log.Fields{
				"repo": repo,
				"error": err,
			}).Error("error on refresh charts from repo")
			continue
		}

		// 3 - Process elements in index
		var chartsWithData []*swaggermodels.ChartPackage
		// Buffered channel
		ch := make(chan chanItem, len(charts))
		defer close(ch)

		// 3.1 - parallellize processing
		for _, chart := range charts {
			go processChartMetadata(chart, repo.URL, ch)
		}
		// 3.2 Channel drain
		for range charts {
			it := <-ch
			// Only append the ones that have not failed
			if it.err == nil {
				chartsWithData = append(chartsWithData, it.chart)
			}
		}
		updatedCharts[repo.Name] = chartsWithData
	}

	// 4 - Update the stored cache with the new elements if everything went well
	c.rwm.Lock()
	c.allCharts = updatedCharts
	c.rwm.Unlock()
	return nil
}

// Represents every element processed in paralell
type chanItem struct {
	chart *swaggermodels.ChartPackage
	err   error
}

// Counting semaphore, 25 downloads max in paralell
var tokens = make(chan struct{}, 15)

func processChartMetadata(chart *swaggermodels.ChartPackage, repoURL string, out chan<- chanItem) {
	tokens <- struct{}{}
	// Semaphore control channel
	defer func() {
		<-tokens
	}()

	var it chanItem
	it.chart = chart

	// Extra files. Skipped if the directory exists
	dataExists, err := charthelper.ChartDataExists(chart)
	if err != nil {
		it.err = err
		out <- it
		return
	}

	if !dataExists {
		log.WithFields(log.Fields{
			"name":    *chart.Name,
			"version": *chart.Version,
		}).Info("Local cache missing")

		err := charthelper.DownloadAndExtractChartTarball(chart, repoURL)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error on DownloadAndExtractChartTarball")
			// Skip chart if error extracting the tarball
			it.err = err
			out <- it
			return
		}
		// If we have a problem processing an image it will fallback to the default one
		err = charthelper.DownloadAndProcessChartIcon(chart)
		if err != nil {
			log.WithFields(log.Fields{
				"chart": *chart.Name,
				"error": err,
			}).Error("Error on DownloadAndProcessChartIcon")
		}
	}
	out <- it
}
