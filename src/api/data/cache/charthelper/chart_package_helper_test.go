package charthelper

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

var repoURL = "http://storage.googleapis.com/kubernetes-charts/"

func TestDownloadAndExtractChartTarballOk(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	assert.NoErr(t, DownloadAndExtractChartTarball(chart, repoURL))

	// Relative URLs
	chart, err = getTestChart()
	assert.NoErr(t, err)
	chart.Urls[0] = "drupal-0.3.0.tgz"
	assert.NoErr(t, DownloadAndExtractChartTarball(chart, repoURL))
}

func TestDownloadAndExtractChartTarballErrorCantWrite(t *testing.T) {
	ensureChartDataDirOrig := ensureChartDataDir
	ensureChartDataDir = func(chart *models.ChartPackage) error { return errors.New("Can't create dir") }
	defer func() { ensureChartDataDir = ensureChartDataDirOrig }()

	chart, err := getTestChart()
	assert.NoErr(t, err)
	err = DownloadAndExtractChartTarball(chart, repoURL)
	assert.ExistsErr(t, err, "trying to create a non valid directory")
}

func TestDownloadAndExtractChartTarballErrorDownload(t *testing.T) {
	// Stubs
	downloadTarballOrig := downloadTarball
	defer func() { downloadTarball = downloadTarballOrig }()
	downloadTarball = func(chart *models.ChartPackage, url string) error { return errors.New("Can't download") }
	tarballExistsOrig := tarballExists
	defer func() { tarballExists = tarballExistsOrig }()
	tarballExists = func(chart *models.ChartPackage) bool { return false }
	// EO stubs
	chart, err := getTestChart()
	assert.NoErr(t, err)
	err = DownloadAndExtractChartTarball(chart, repoURL)
	assert.ExistsErr(t, err, "Error downloading the tar file")
	_, err = os.Stat(chartDataDir(chart))
	assert.ExistsErr(t, err, "chart data dir has been removed")
}

func TestDownloadAndExtractChartTarballErrorExtract(t *testing.T) {
	extractFilesFromTarballOrig := extractFilesFromTarball
	defer func() { extractFilesFromTarball = extractFilesFromTarballOrig }()
	extractFilesFromTarball = func(chart *models.ChartPackage) error { return errors.New("Can't download") }
	chart, err := getTestChart()
	assert.NoErr(t, err)
	err = DownloadAndExtractChartTarball(chart, repoURL)
	assert.ExistsErr(t, err, "Error extracting tar file content")
	_, err = os.Stat(chartDataDir(chart))
	assert.ExistsErr(t, err, "chart data dir has been removed")
}

// It creates the tar file in local filesystem
func TestDownloadTarballCreatesFileInDestination(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	// Stubs
	// Disable remote URL download
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://storage.googleapis.com/kubernetes-charts/drupal-0.3.0.tgz",
		httpmock.NewStringResponder(200, "Mocked Response"))

	// Mock download path
	chartDataDirOrig := chartDataDir
	defer func() { chartDataDir = chartDataDirOrig }()
	pathExists, _ := ioutil.TempDir(os.TempDir(), "chart")
	chartDataDir = func(c *models.ChartPackage) string {
		return pathExists
	}

	err = downloadTarball(chart, repoURL)
	assert.NoErr(t, err)
	_, err = os.Stat(filepath.Join(chartDataDir(chart), "chart.tgz"))
	assert.NoErr(t, err)
}

func TestDownloadTarballErrorDownloading(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	randomPath, _ := ioutil.TempDir(os.TempDir(), "test")
	TarballPathOrig := TarballPath
	defer func() { TarballPath = TarballPathOrig }()
	TarballPath = func(chart *models.ChartPackage) string {
		return filepath.Join(randomPath, "myFile.tar.gz")
	}
	// 404
	chart.Urls[0] = "http://localhost/bogusUrl"
	err = downloadTarball(chart, repoURL)
	assert.ExistsErr(t, err, "Cant copy")
}

func TestExtractFilesFromTarballOk(t *testing.T) {
	chart, err := getTestChart()
	ensureChartDataDir(chart)
	assert.NoErr(t, err)
	// Stubs
	TarballPathOrig := TarballPath
	defer func() { TarballPath = TarballPathOrig }()
	TarballPath = func(chart *models.ChartPackage) string {
		path := MockedTarballPath()
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
	TarballPathOrig := TarballPath
	defer func() { TarballPath = TarballPathOrig }()
	TarballPath = func(chart *models.ChartPackage) string {
		return "/does-not-exist.tar.gz"
	}
	err = extractFilesFromTarball(chart)
	assert.ExistsErr(t, err, "tar file does not exist")
}

func TestExtractFilesFromTarballCantCopy(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	// Stubs
	TarballPathOrig := TarballPath
	copyFileOrig := copyFile
	defer func() { TarballPath = TarballPathOrig; copyFile = copyFileOrig }()
	TarballPath = func(chart *models.ChartPackage) string {
		path := MockedTarballPath()
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
	TarballPathOrig := TarballPath
	defer func() { TarballPath = TarballPathOrig }()
	TarballPath = func(chart *models.ChartPackage) string {
		path := MockedTarballPath()
		return path
	}
	ensureChartDataDir(chart)
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
	DataDirBaseOrig := DataDirBase
	defer func() { DataDirBase = DataDirBaseOrig }()
	DataDirBase = func() string {
		return randomPath
	}
	chartPath := chartDataDir(chart)

	_, err = os.Stat(chartPath)
	assert.ExistsErr(t, err, "File does not exist")
	ensureChartDataDir(chart)
	_, err = os.Stat(chartPath)
	assert.NoErr(t, err)
}

func TestCleanChartDataDir(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	randomPath, _ := ioutil.TempDir(os.TempDir(), "chart")
	DataDirBaseOrig := DataDirBase
	defer func() { DataDirBase = DataDirBaseOrig }()
	DataDirBase = func() string {
		return randomPath
	}
	chartPath := chartDataDir(chart)
	ensureChartDataDir(chart)
	_, err = os.Stat(chartPath)
	assert.NoErr(t, err)

	err = cleanChartDataDir(chart)
	assert.NoErr(t, err)

	_, err = os.Stat(chartPath)
	assert.ExistsErr(t, err, "chart dir removed")
}

// Required because DataDirBase has been overriden
var origDataDirBase = DataDirBase()

func TestDataDirBase(t *testing.T) {
	path := filepath.Join(config.BaseDir(), "repo-data")
	assert.Equal(t, origDataDirBase, path, "dataDirbase uses BaseDir")
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
	exists, _ := ChartDataExists(chart)
	assert.Equal(t, exists, true, "the directory exists")
	chartDataDir = func(c *models.ChartPackage) string {
		return "/does-not-exist"
	}
	// Directory does not exist
	exists, _ = ChartDataExists(chart)
	assert.Equal(t, exists, false, "the directory does not exist")
}

func TestCopyFile(t *testing.T) {
	src, _ := ioutil.TempFile(os.TempDir(), "")
	dest, _ := ioutil.TempFile(os.TempDir(), "")
	err := copyFile(dest.Name(), src.Name())
	assert.NoErr(t, err)
	// Src does not exist
	err = copyFile(dest.Name(), "/does-not-exist")
	assert.ExistsErr(t, err, "SRC does not exist")
	// Dest is not valid
	err = copyFile(os.TempDir(), src.Name())
	assert.ExistsErr(t, err, "Destination not valid")
}

func TestUntar(t *testing.T) {
	src := MockedTarballPath()
	dest, _ := ioutil.TempDir(os.TempDir(), "")
	err := untar(src, dest)
	assert.NoErr(t, err)
	files, _ := ioutil.ReadDir(dest)
	assert.Equal(t, len(files), 1, "only contains the parent dir")

	// Cant read the tar file
	err = untar("does-not-exist.tar.gz", dest)
	assert.ExistsErr(t, err, "SRC does not exist")

	// It is not valid gzip file
	notValidFile, _ := ioutil.TempFile(os.TempDir(), "")
	err = untar(notValidFile.Name(), dest)
	assert.ExistsErr(t, err, "SRC does not exist")
}

func TestReadmeStaticURL(t *testing.T) {
	chart, err := getTestChart()
	assert.NoErr(t, err)
	res := ReadmeStaticUrl(chart, "prefix")
	readmePath := filepath.Join(chartDataDir(chart), "README.md")
	assert.Equal(t, res, staticUrl(readmePath, "prefix"), "Readme file ")
}

// Auxiliar
func getTestChart() (*models.ChartPackage, error) {
	randomPath, _ := ioutil.TempDir(os.TempDir(), "chart")
	DataDirBase = func() string {
		return randomPath
	}
	url := "http://storage.googleapis.com/kubernetes-charts/drupal-0.3.0.tgz"
	version := "0.3.0"
	name := "drupal"
	return &models.ChartPackage{
		Urls:    []string{url},
		Version: &version,
		Name:    &name,
		Repo:    "stable",
	}, nil
}

// Returns the test tarball path
func MockedTarballPath() string {
	return filepath.Join("testdata", "drupal-0.3.0.tgz")
}
