package cache

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/cache/charthelper"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/charts"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

var dbSession = models.NewMockSession(models.MockDBConfig{})
var chartsImplementation = getChartsImplementation(dbSession)

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
	params := charts.SearchChartsParams{
		Name: "drupal",
	}
	charts, err := chartsImplementation.Search(params)
	assert.NoErr(t, err)
	// flatten chart+version results into a chart resource array
	db, _ := dbSession.DB()
	resources := helpers.MakeChartResources(db, charts)
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
	charthelper.DownloadAndExtractChartTarball = func(chart *swaggermodels.ChartPackage, repoURL string) error { return nil }

	DownloadAndProcessChartIconOrig := charthelper.DownloadAndProcessChartIcon
	defer func() { charthelper.DownloadAndProcessChartIcon = DownloadAndProcessChartIconOrig }()
	charthelper.DownloadAndProcessChartIcon = func(chart *swaggermodels.ChartPackage) error { return nil }

	// EO stubs

	err := chartsImplementation.Refresh()
	assert.NoErr(t, err)
}

func TestCachedChartsRefreshErrorPropagation(t *testing.T) {
	tests := []struct {
		name  string
		repos []*models.Repo
	}{
		{"invalid repo url", []*models.Repo{{Name: "stable", URL: "./localhost"}}},
		{"inexistant repo", []*models.Repo{{Name: "stable", URL: "http://localhost"}}},
	}

	defer func() { models.MockRepos = models.OfficialRepos }()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chImplementation := NewCachedCharts(models.NewMockSession(models.MockDBConfig{}))
			models.MockRepos = tt.repos
			err := chImplementation.Refresh()
			assert.ExistsErr(t, err, tt.name)
		})
	}
}

func TestCachedChartsRefreshErrorDownloadingPackage(t *testing.T) {
	ChartDataExistsOrig := charthelper.ChartDataExists
	defer func() { charthelper.ChartDataExists = ChartDataExistsOrig }()
	charthelper.ChartDataExists = func(chart *swaggermodels.ChartPackage) (bool, error) { return false, nil }

	DownloadAndExtractChartTarballOrig := charthelper.DownloadAndExtractChartTarball
	defer func() { charthelper.DownloadAndExtractChartTarball = DownloadAndExtractChartTarballOrig }()
	knownError := errors.New("error on DownloadAndExtractChartTarball")
	charthelper.DownloadAndExtractChartTarball = func(chart *swaggermodels.ChartPackage, repoURL string) error {
		return knownError
	}

	chImplementation := NewCachedCharts(dbSession)
	// It does not return error
	err := chImplementation.Refresh()
	assert.NoErr(t, err)
}

func getChartsImplementation(dbSession datastore.Session) data.Charts {
	// Stub ChartDataExists to avoid downloading extra data
	ChartDataExistsOrig := charthelper.ChartDataExists
	defer func() { charthelper.ChartDataExists = ChartDataExistsOrig }()
	charthelper.ChartDataExists = func(chart *swaggermodels.ChartPackage) (bool, error) {
		return true, nil
	}

	// configure the api here
	chImplementation := NewCachedCharts(dbSession)
	chImplementation.Refresh()
	return chImplementation
}
