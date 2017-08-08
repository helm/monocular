package cache

import (
	"errors"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/config/repos"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/cache/charthelper"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/data/util"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/charts"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

var chartsImplementation = getChartsImplementation()

func TestCachedChartsChartFromRepo(t *testing.T) {
	setupTestRepoCache(nil)
	defer teardownTestRepoCache()
	// TODO: validate chart data
	_, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	_, err = chartsImplementation.ChartFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to Charts.ChartFromRepo()")
	_, err = chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus chart name to Charts.ChartFromRepo()")
}

func TestCachedChartsChartVersionFromRepo(t *testing.T) {
	setupTestRepoCache(nil)
	defer teardownTestRepoCache()
	chart, err := chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, testutil.ChartVersionString)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.Name, testutil.ChartName, "chart name")
	assert.Equal(t, *chart.Version, testutil.ChartVersionString, "chart version")
	_, err = chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, "99.99.99")
	assert.ExistsErr(t, err, "sent bogus chart version to ChartVersionFromRepo")
	_, err = chartsImplementation.ChartVersionFromRepo(testutil.BogusRepo, testutil.ChartName, testutil.ChartVersionString)
	assert.ExistsErr(t, err, "sent bogus repo name to Charts.ChartFromRepo()")
	_, err = chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.BogusRepo, testutil.ChartVersionString)
	assert.ExistsErr(t, err, "sent bogus chart name to Charts.ChartFromRepo()")
}

func TestCachedChartsChartVersionsFromRepo(t *testing.T) {
	setupTestRepoCache(nil)
	defer teardownTestRepoCache()
	charts, err := chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := chartsImplementation.ChartVersionsFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}

func TestCachedChartsAll(t *testing.T) {
	setupTestRepoCache(nil)
	defer teardownTestRepoCache()
	_, err := chartsImplementation.All()
	assert.NoErr(t, err)
}

func TestCachedChartsSearch(t *testing.T) {
	setupTestRepoCache(nil)
	defer teardownTestRepoCache()
	params := charts.SearchChartsParams{
		Name: "drupal",
	}
	charts, err := chartsImplementation.Search(params)
	assert.NoErr(t, err)
	// flatten chart+version results into a chart resource array
	resources := helpers.MakeChartResources(charts)
	assert.Equal(t, len(resources), 1, "number of unique chart results")
}

func TestCachedChartsAllFromRepo(t *testing.T) {
	setupTestRepoCache(nil)
	defer teardownTestRepoCache()
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := chartsImplementation.AllFromRepo(testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}

func TestCachedChartsRefresh(t *testing.T) {
	setupTestRepoCache(nil)
	defer teardownTestRepoCache()
	// Stubs Download and processing
	DownloadAndExtractChartTarballOrig := charthelper.DownloadAndExtractChartTarball
	defer func() { charthelper.DownloadAndExtractChartTarball = DownloadAndExtractChartTarballOrig }()
	charthelper.DownloadAndExtractChartTarball = func(chart *models.ChartPackage) error { return nil }

	DownloadAndProcessChartIconOrig := charthelper.DownloadAndProcessChartIcon
	defer func() { charthelper.DownloadAndProcessChartIcon = DownloadAndProcessChartIconOrig }()
	charthelper.DownloadAndProcessChartIcon = func(chart *models.ChartPackage) error { return nil }

	// EO stubs

	err := chartsImplementation.Refresh()
	assert.NoErr(t, err)
}

func TestCachedChartsRefreshErrorPropagation(t *testing.T) {
	// Invalid repo URL
	rep := []models.Repo{
		models.Repo{
			Name: util.StrToPtr("stable"),
			URL:  util.StrToPtr("./localhost"),
		},
	}
	setupTestRepoCache(&rep)
	chImplementation := NewCachedCharts()
	err := chImplementation.Refresh()
	assert.ExistsErr(t, err, "Invalid Repo URL")

	teardownTestRepoCache()
	// Repo does not exist
	rep = repos.Repos{
		models.Repo{
			Name: util.StrToPtr("stable"),
			URL:  util.StrToPtr("http://localhost"),
		},
	}
	setupTestRepoCache(&rep)
	defer teardownTestRepoCache()
	chImplementation = NewCachedCharts()
	err = chImplementation.Refresh()
	assert.ExistsErr(t, err, "Repo does not exist")
}

func TestCachedChartsRefreshErrorDownloadingPackage(t *testing.T) {
	ChartDataExistsOrig := charthelper.ChartDataExists
	defer func() { charthelper.ChartDataExists = ChartDataExistsOrig }()
	charthelper.ChartDataExists = func(chart *models.ChartPackage) (bool, error) { return false, nil }

	DownloadAndExtractChartTarballOrig := charthelper.DownloadAndExtractChartTarball
	defer func() { charthelper.DownloadAndExtractChartTarball = DownloadAndExtractChartTarballOrig }()
	knownError := errors.New("error on DownloadAndExtractChartTarball")
	charthelper.DownloadAndExtractChartTarball = func(chart *models.ChartPackage) error {
		return knownError
	}

	repos := []models.Repo{
		models.Repo{
			Name: util.StrToPtr("stable"),
			URL:  util.StrToPtr("http://storage.googleapis.com/kubernetes-charts"),
		},
	}
	setupTestRepoCache(&repos)
	defer teardownTestRepoCache()
	chImplementation := NewCachedCharts()
	// It does not return error
	err := chImplementation.Refresh()
	assert.NoErr(t, err)
}

func getChartsImplementation() data.Charts {
	// Stub ChartDataExists to avoid downloading extra data
	ChartDataExistsOrig := charthelper.ChartDataExists
	defer func() { charthelper.ChartDataExists = ChartDataExistsOrig }()
	charthelper.ChartDataExists = func(chart *models.ChartPackage) (bool, error) {
		return true, nil
	}

	// configure the api here
	chImplementation := NewCachedCharts()
	return chImplementation
}

func setupTestRepoCache(repos *[]models.Repo) {
	config.NewRedisPool()
	if repos == nil {
		repos = &[]models.Repo{
			models.Repo{
				Name: util.StrToPtr("stable"),
				URL:  util.StrToPtr("http://storage.googleapis.com/kubernetes-charts"),
			},
			models.Repo{
				Name: util.StrToPtr("incubator"),
				URL:  util.StrToPtr("http://storage.googleapis.com/kubernetes-charts-incubator"),
			},
		}
	}
	NewCachedRepos(*repos)
	chartsImplementation.Refresh()
}

func teardownTestRepoCache() {
	if _, err := Repos.DeleteAll(); err != nil {
		log.Fatal("could not clear cache")
	}
	if err := config.Pool.Close(); err != nil {
		log.Fatal("could not close redis pool")
	}
}
