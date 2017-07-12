package charthelper

import (
	"errors"
	"io/ioutil"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

func TestDownloadAndProcessChartIconNoIcon(t *testing.T) {
	chart, _ := getTestChart()
	assert.NoErr(t, DownloadAndProcessChartIcon(chart))
}

func TestDownloadAndProcessChartIconOk(t *testing.T) {
	chart, _ := getTestChartWithIcon()
	defer httpmock.DeactivateAndReset()
	assert.NoErr(t, DownloadAndProcessChartIcon(chart))
}

func TestDownloadAndProcessChartIconErrorCantCreateDIR(t *testing.T) {
	ensureChartDataDirOrig := ensureChartDataDir
	ensureChartDataDir = func(chart *models.ChartPackage) error { return errors.New("Can't create dir") }
	defer func() { ensureChartDataDir = ensureChartDataDirOrig }()
	chart, _ := getTestChartWithIcon()
	defer httpmock.DeactivateAndReset()
	assert.ExistsErr(t, DownloadAndProcessChartIcon(chart), "Can't create destination path")
}
func TestDownloadAndProcessChartIconErrorDownload(t *testing.T) {
	chart, _ := getTestChartWithIcon()
	httpmock.DeactivateAndReset()
	assert.ExistsErr(t, DownloadAndProcessChartIcon(chart), "Can't download")
}

func TestDownloadAndProcessChartIconErrorProcessing(t *testing.T) {
	chart, _ := getTestChartWithIcon()
	defer httpmock.DeactivateAndReset()
	icon, _ := ioutil.ReadFile("testdata/noticon.txt")
	httpmock.RegisterResponder("GET", chart.Icon,
		httpmock.NewBytesResponder(200, icon))
	assert.ExistsErr(t, DownloadAndProcessChartIcon(chart), "Can't process")
}

// Download Icon
func TestDownloadIconError(t *testing.T) {
	chart, _ := getTestChart()

	chart.Icon = ""
	err := downloadIcon(chart)
	assert.ExistsErr(t, err, "No Icon")

	// Invald protocol
	chart.Icon = "./invalid-protocol/bogusUrl"
	err = downloadIcon(chart)
	assert.ExistsErr(t, err, "Invalid protocol")

	// 404
	chart.Icon = "http://localhost/bogusUrl"
	err = downloadIcon(chart)
	assert.ExistsErr(t, err, "It does not exist")
}

func TestProcessIconNoOriginalIcon(t *testing.T) {
	chart, _ := getTestChart()
	err := processIcon(chart)
	assert.ExistsErr(t, err, "Original icon can not be processed")
}

func TestProcessIconOk(t *testing.T) {
	chart, _ := getTestChartWithIcon()
	defer httpmock.DeactivateAndReset()
	availableFormats = []iconFormat{
		{"160x160-fit", 160, 160, fit},
		{"160x160-fill", 160, 160, fill},
	}

	ensureChartDataDir(chart)
	err := downloadIcon(chart)
	assert.NoErr(t, err)
	files, _ := ioutil.ReadDir(chartDataDir(chart))
	assert.Equal(t, len(files), 1, "Only the original")
	assert.NoErr(t, processIcon(chart))
	files, _ = ioutil.ReadDir(chartDataDir(chart))
	assert.Equal(t, len(files), len(availableFormats)+1, "Original + generated")
}

func TestIconExist(t *testing.T) {
	chart, _ := getTestChartWithIcon()
	defer httpmock.DeactivateAndReset()
	ensureChartDataDir(chart)
	err := downloadIcon(chart)
	assert.NoErr(t, err)
	exist, _ := iconExist(chart, "original")
	assert.Equal(t, exist, true, "It has the original icon")

	exist, _ = iconExist(chart, "bogus")
	assert.Equal(t, exist, false, "Can't find the icon format")
}

func TestAvailableIcons(t *testing.T) {
	chart, _ := getTestChartWithIcon()
	defer httpmock.DeactivateAndReset()
	ensureChartDataDir(chart)
	err := downloadIcon(chart)
	assert.NoErr(t, err)
	icons := AvailableIcons(chart, "myPrefix")
	assert.Equal(t, len(icons), 0, "Empty array if no files present")

	// Generate icons
	assert.NoErr(t, processIcon(chart))
	icons = AvailableIcons(chart, "myPrefix")
	assert.Equal(t, len(icons), len(availableFormats), "Include all the formats")
	for i, format := range availableFormats {
		assert.Equal(t, icons[i].Name, format.Name, "Same name")
		path, _ := iconPath(chart, format.Name)
		assert.Equal(t, icons[i].Path, staticUrl(path, "myPrefix"), "Same path")
	}
	// Error retrieving errors
	iconExistOrig := iconExist
	defer func() { iconExist = iconExistOrig }()
	iconExist = func(chart *models.ChartPackage, format string) (bool, error) {
		return false, errors.New("Error raised")
	}
	icons = AvailableIcons(chart, "myPrefix")
	assert.Equal(t, len(icons), 0, "It skipped the ones that failed")
}

func getTestChartWithIcon() (*models.ChartPackage, error) {
	chart, _ := getTestChart()
	chart.Icon = "http://localhost/mock-icon.png"
	httpmock.Activate()
	icon, _ := ioutil.ReadFile("testdata/icon.png")
	httpmock.RegisterResponder("GET", chart.Icon,
		httpmock.NewBytesResponder(200, icon))
	return chart, nil
}
