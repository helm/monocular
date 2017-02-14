package cache

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/cache/charthelper"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/data/repos"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/swagger/restapi/operations"
	"github.com/helm/monocular/src/api/testutil"
)

var chartsImplementation = getChartsImplementation()

func TestCachedChartsChartFromRepo(t *testing.T) {
	// TODO: validate chart data
	_, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	_, err = chartsImplementation.ChartFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to Charts.ChartFromRepo()")
	_, err = chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus chart name to Charts.ChartFromRepo()")
}

func TestCachedChartsChartVersionFromRepo(t *testing.T) {
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
	charts, err := chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := chartsImplementation.ChartVersionsFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}

func TestCachedChartsAll(t *testing.T) {
	_, err := chartsImplementation.All()
	assert.NoErr(t, err)
}

func TestCachedChartsSearch(t *testing.T) {
	params := operations.SearchChartsParams{
		Name: "drupal",
	}
	charts, err := chartsImplementation.Search(params)
	assert.NoErr(t, err)
	// flatten chart+version results into a chart resource array
	resources := helpers.MakeChartResources(charts)
	assert.Equal(t, len(resources), 1, "number of unique chart results")
}

func TestCachedChartsAllFromRepo(t *testing.T) {

	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := chartsImplementation.AllFromRepo(testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}

func TestCachedChartsRefresh(t *testing.T) {
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
	rep := repos.Repos{
		repos.Repo{
			Name: "stable",
			URL:  "./localhost",
		},
	}
	chImplementation := NewCachedCharts(rep)
	err := chImplementation.Refresh()
	assert.ExistsErr(t, err, "Invalid Repo URL")
	// Repo does not exist
	rep = repos.Repos{
		repos.Repo{
			Name: "stable",
			URL:  "http://localhost",
		},
	}
	chImplementation = NewCachedCharts(rep)
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

	repos := repos.Repos{
		repos.Repo{
			Name: "stable",
			URL:  "http://storage.googleapis.com/kubernetes-charts",
		},
	}
	chImplementation := NewCachedCharts(repos)
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
	repos := repos.Repos{
		repos.Repo{
			Name: "stable",
			URL:  "http://storage.googleapis.com/kubernetes-charts",
		},
		repos.Repo{
			Name: "incubator",
			URL:  "http://storage.googleapis.com/kubernetes-charts-incubator",
		},
	}
	chImplementation := NewCachedCharts(repos)
	chImplementation.Refresh()
	return chImplementation
}
