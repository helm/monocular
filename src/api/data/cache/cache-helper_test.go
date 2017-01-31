package cache

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/arschles/assert"
	"github.com/helm/monocular/src/api/mocks"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/testutil"
)

func TestDownloadAndExtractChartTarballOk(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	assert.NoErr(t, downloadAndExtractChartTarball(chart))
}

func TestDownloadAndExtractChartTarballErrorCantWrite(t *testing.T) {
	ensureChartDataDirOrig := ensureChartDataDir
	ensureChartDataDir = func(chart *models.ChartPackage) error { return errors.New("Can't create dir") }
	defer func() { ensureChartDataDir = ensureChartDataDirOrig }()

	chart, err := getTestChart()
	assert.NoErr(t, err)
	err = downloadAndExtractChartTarball(chart)
	assert.ExistsErr(t, err, "trying to create a non valid directory")
}

func TestDownloadAndExtractChartTarballErrorDownload(t *testing.T) {
	// Stubs
	downloadTarballOrig := downloadTarball
	defer func() { downloadTarball = downloadTarballOrig }()
	downloadTarball = func(chart *models.ChartPackage) error { return errors.New("Can't download") }
	tarballExistsOrig := tarballExists
	defer func() { tarballExists = tarballExistsOrig }()
	tarballExists = func(chart *models.ChartPackage) bool { return false }
	// EO stubs
	chart, err := getTestChart()
	assert.NoErr(t, err)
	err = downloadAndExtractChartTarball(chart)
	assert.ExistsErr(t, err, "Error downloading the tar file")
}

func TestDownloadAndExtractChartTarballErrorExtract(t *testing.T) {
	extractFilesFromTarballOrig := extractFilesFromTarball
	defer func() { extractFilesFromTarball = extractFilesFromTarballOrig }()
	extractFilesFromTarball = func(chart *models.ChartPackage) error { return errors.New("Can't download") }
	chart, err := getTestChart()
	assert.NoErr(t, err)
	err = downloadAndExtractChartTarball(chart)
	assert.ExistsErr(t, err, "Error extracting tar file content")
}

// It creates the tar file in local filesystem
func TestDownloadTarballCreatesFileInDestination(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	// Stubs
	// Disable remote URL download
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://kubernetes-charts.storage.googleapis.com/drupal-0.4.1.tgz",
		httpmock.NewStringResponder(200, "Mocked Response"))

	// Mock download path
	randomPath, _ := ioutil.TempDir(os.TempDir(), "test")
	tarballTmpPathOrig := tarballTmpPath
	defer func() { tarballTmpPath = tarballTmpPathOrig }()
	tarballTmpPath = func(chart *models.ChartPackage) string {
		return filepath.Join(randomPath, "myFile.tar.gz")
	}
	err = downloadTarball(chart)
	assert.NoErr(t, err)
	_, err = os.Stat(tarballTmpPath(chart))
	assert.NoErr(t, err)
}

func TestExtractFilesFromTarballOk(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	// Stubs
	tarballTmpPathOrig := tarballTmpPath
	defer func() { tarballTmpPath = tarballTmpPathOrig }()
	tarballTmpPath = func(chart *models.ChartPackage) string {
		path, _ := mocks.MockedtarballTmpPath()
		return path
	}
	err = extractFilesFromTarball(chart)
	assert.NoErr(t, err)
	// Saves all the files in filesToKeep
	for _, fileName := range filesToKeep {
		file := filepath.Join(chartDataDir(chart), fileName)
		_, err = os.Stat(file)
		assert.NoErr(t, err)
	}
	files, _ := ioutil.ReadDir(chartDataDir(chart))
	assert.Equal(t, len(files), len(filesToKeep), "only contain the wanted files")
}

func TestExtractFilesFromTarballNotFound(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	// Stubs
	tarballTmpPathOrig := tarballTmpPath
	defer func() { tarballTmpPath = tarballTmpPathOrig }()
	tarballTmpPath = func(chart *models.ChartPackage) string {
		return "/does-not-exist.tar.gz"
	}
	err = extractFilesFromTarball(chart)
	assert.ExistsErr(t, err, "tar file does not exist")
}

func TestExtractFilesFromTarballCantCopy(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	// Stubs
	tarballTmpPathOrig := tarballTmpPath
	copyFileOrig := copyFile
	defer func() { tarballTmpPath = tarballTmpPathOrig; copyFile = copyFileOrig }()
	tarballTmpPath = func(chart *models.ChartPackage) string {
		path, _ := mocks.MockedtarballTmpPath()
		return path
	}

	copyFile = func(dst, src string) error { return errors.New("Can't copy") }
	err = extractFilesFromTarball(chart)
	assert.Err(t, err, errors.New("Can't copy"))
}

// Can't access to files from the cache
func TestReadFromCacheOk(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	// Stubs
	tarballTmpPathOrig := tarballTmpPath
	defer func() { tarballTmpPath = tarballTmpPathOrig }()
	tarballTmpPath = func(chart *models.ChartPackage) string {
		path, _ := mocks.MockedtarballTmpPath()
		return path
	}
	err = extractFilesFromTarball(chart)
	assert.NoErr(t, err)
	for _, fileName := range filesToKeep {
		res, err := ReadFromCache(chart, fileName)
		assert.NoErr(t, err)
		assert.NotNil(t, res, "string content")
	}
}

func TestReadFromCacheNotFound(t *testing.T) {
	chart, _ := getTestChart()
	res, err := ReadFromCache(chart, "not-found")
	assert.ExistsErr(t, err, "Can't find local file not-found")
	assert.Equal(t, res, "", "empty string")
}

func TestEnsureChartDataDir(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	randomPath, _ := ioutil.TempDir(os.TempDir(), "chart")
	dataDirBaseOrig := dataDirBase
	defer func() { dataDirBase = dataDirBaseOrig }()
	dataDirBase = func() string {
		return randomPath
	}
	chartPath := chartDataDir(chart)

	_, err = os.Stat(chartPath)
	assert.ExistsErr(t, err, "File does not exist")
	ensureChartDataDir(chart)
	_, err = os.Stat(chartPath)
	assert.NoErr(t, err)
}

func TestChartDataExist(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	chartDataDirOrig := chartDataDir
	defer func() { chartDataDir = chartDataDirOrig }()
	pathExists, _ := ioutil.TempDir(os.TempDir(), "chart")
	chartDataDir = func(c *models.ChartPackage) string {
		return pathExists
	}
	// Directory exists
	exists, _ := chartDataExists(chart)
	assert.Equal(t, exists, true, "the directory exists")
	chartDataDir = func(c *models.ChartPackage) string {
		return "/does-not-exist"
	}
	// Directory does not exist
	exists, _ = chartDataExists(chart)
	assert.Equal(t, exists, false, "the directory does not exist")
}

// Auxiliar
func getTestChart() (*models.ChartPackage, error) {
	return chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
}
