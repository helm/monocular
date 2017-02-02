package cache

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/helm/monocular/src/api/swagger/models"
)

// Downloads the chart tar file linked by metadata.Urls and store
// the wanted files (i.e README.md) under chartDataDir
var downloadAndExtractChartTarball = func(chart *models.ChartPackage) error {
	if err := ensureChartDataDir(chart); err != nil {
		return err
	}

	if !tarballExists(chart) {
		if err := downloadTarball(chart); err != nil {
			return err
		}
	}

	if err := extractFilesFromTarball(chart); err != nil {
		return err
	}

	return nil
}

var tarballExists = func(chart *models.ChartPackage) bool {
	_, err := os.Stat(tarballTmpPath(chart))
	return err == nil
}

// Downloads the tar.gz file associated with the chart version exposed by the index
// in order to extract specific files for caching
var downloadTarball = func(chart *models.ChartPackage) error {
	source := chart.Urls[0]
	destination := tarballTmpPath(chart)

	// Create output
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Printf("Downloading metadata from %s\n", source)
	// Download tarball
	resp, err := http.Get(source)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error downloading %s, %d", source, resp.StatusCode)
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

var filesToKeep = []string{"README.md"}

var extractFilesFromTarball = func(chart *models.ChartPackage) error {
	tarballPath := tarballTmpPath(chart)
	tarballExpandedPath, err := ioutil.TempDir(os.TempDir(), "chart")
	if err != nil {
		return err
	}

	// Extract
	if _, err := os.Stat(tarballPath); err != nil {
		return fmt.Errorf("Can not find file to extract %s", tarballPath)
	}

	if err := untar(tarballPath, tarballExpandedPath); err != nil {
		return err
	}

	// Save specific files defined by filesToKeep
	// Include /[chartName] namespace
	chartPath := filepath.Join(tarballExpandedPath, *chart.Name)
	for _, fileName := range filesToKeep {
		src := filepath.Join(chartPath, fileName)
		dest := filepath.Join(chartDataDir(chart), fileName)
		fmt.Printf("Storing %s\n", dest)
		if err := copyFile(dest, src); err != nil {
			return err
		}
	}
	return nil
}

var ensureChartDataDir = func(chart *models.ChartPackage) error {
	dir := chartDataDir(chart)
	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("Could not create %s: %s", dir, err)
		}
	}
	return nil
}

// Temporary path used for downloaded tarball
var tarballTmpPath = func(chart *models.ChartPackage) string {
	splitTarURL := strings.Split(chart.Urls[0], "/")
	return filepath.Join("/tmp", splitTarURL[len(splitTarURL)-1])
}

// Directory used to store cached data like readme files
// Variable so it can be mocked
var dataDirBase = func() string {
	return filepath.Join(os.Getenv("HOME"), "repo-data")
}

// Data directory with cached content for the current ChartPackage
var chartDataDir = func(chart *models.ChartPackage) string {
	return filepath.Join(dataDirBase(), chart.Repo, *chart.Name, *chart.Version)
}

// Credit https://github.com/kubernetes/helm/blob/8120508c1bd73f6f49d1f16f7a6eacbb4e707655/pkg/chartutil/expand.go#L28
func untar(tarball, dir string) error {
	r, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer r.Close()

	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		//split header name and create missing directories
		d, _ := filepath.Split(header.Name)
		fullDir := filepath.Join(dir, d)
		_, err = os.Stat(fullDir)
		if err != nil && d != "" {
			if err = os.MkdirAll(fullDir, 0700); err != nil {
				return err
			}
		}

		path := filepath.Clean(filepath.Join(dir, header.Name))
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}
	}
	return nil
}

var copyFile = func(dst, src string) error {
	i, err := os.Open(src)
	if err != nil {
		return err
	}
	defer i.Close()

	o, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = io.Copy(o, i)

	return err
}

// ReadFromCache reads a file from the local cache based on the chart and filename provided
var ReadFromCache = func(chart *models.ChartPackage, filename string) (string, error) {
	fileToLoad := filepath.Join(chartDataDir(chart), filename)
	dat, err := ioutil.ReadFile(fileToLoad)

	if err != nil {
		return "", fmt.Errorf("Can't find local file %s", fileToLoad)
	}

	return string(dat), nil
}

var chartDataExists = func(chart *models.ChartPackage) (bool, error) {
	_, err := os.Stat(chartDataDir(chart))
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
