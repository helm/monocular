package charthelper

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

const defaultTimeout time.Duration = 10 * time.Second

// DownloadAndExtractChartTarball the chart tar file linked by metadata.Urls and store
// the wanted files (i.e README.md) under chartDataDir
var DownloadAndExtractChartTarball = func(chart *models.ChartPackage, repoURL string) (err error) {
	if err := ensureChartDataDir(chart); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			cleanChartDataDir(chart)
		}
	}()

	if !tarballExists(chart) {
		if err := downloadTarball(chart, repoURL); err != nil {
			return err
		}
	}

	if err := extractFilesFromTarball(chart); err != nil {
		return err
	}

	return nil
}

var tarballExists = func(chart *models.ChartPackage) bool {
	_, err := os.Stat(TarballPath(chart))
	return err == nil
}

// Downloads the tar.gz file associated with the chart version exposed by the index
// in order to extract specific files for caching
var downloadTarball = func(chart *models.ChartPackage, repoURL string) error {
	source := chart.Urls[0]
	if _, err := url.ParseRequestURI(source); err != nil {
		// If the chart URL is not absolute, join with repo URL. It's fine if the
		// URL we build here is invalid as we can catch this error when actually
		// making the request
		u, _ := url.Parse(repoURL)
		u.Path = path.Join(u.Path, source)
		source = u.String()
	}

	destination := TarballPath(chart)

	// Create output
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	log.WithFields(log.Fields{
		"source": source,
		"dest":   destination,
	}).Info("Downloading metadata")
	// Download tarball
	c := &http.Client{
		Timeout: defaultTimeout,
	}
	resp, err := c.Get(source)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"source":     source,
			"statusCode": resp.StatusCode,
		}).Error("Error downloading tarball")
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
	TarballPath := TarballPath(chart)
	tarballExpandedPath, err := ioutil.TempDir(os.TempDir(), "chart")
	if err != nil {
		return err
	}

	// Extract
	if _, err := os.Stat(TarballPath); err != nil {
		return fmt.Errorf("Can not find file to extract %s", TarballPath)
	}

	if err := untar(TarballPath, tarballExpandedPath); err != nil {
		return err
	}

	// Save specific files defined by filesToKeep
	// Include /[chartName] namespace
	chartPath := filepath.Join(tarballExpandedPath, *chart.Name)
	for _, fileName := range filesToKeep {
		src := filepath.Join(chartPath, fileName)
		dest := filepath.Join(chartDataDir(chart), fileName)

		log.WithFields(log.Fields{
			"path": dest,
		}).Info("Storing in cache")

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

func cleanChartDataDir(chart *models.ChartPackage) error {
	return os.RemoveAll(chartDataDir(chart))
}

// TarballPath returns the location of the chart package in the local cache
var TarballPath = func(chart *models.ChartPackage) string {
	return filepath.Join(chartDataDir(chart), "chart.tgz")
}

// DataDirBase is the directory used to store cached data like readme files
// Variable so it can be mocked
var DataDirBase = func() string {
	return filepath.Join(config.BaseDir(), "repo-data")
}

// Data directory with cached content for the current ChartPackage
var chartDataDir = func(chart *models.ChartPackage) string {
	return filepath.Join(DataDirBase(), chart.Repo, *chart.Name, *chart.Version)
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

// ChartDataExists checks if the chart cache directory is present
var ChartDataExists = func(chart *models.ChartPackage) (bool, error) {
	_, err := os.Stat(chartDataDir(chart))
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ReadmeStaticUrl returns the static path for the README.md file
func ReadmeStaticUrl(chart *models.ChartPackage, prefix string) string {
	path := filepath.Join(chartDataDir(chart), "README.md")
	return staticUrl(path, prefix)
}

func staticUrl(path, prefix string) string {
	return strings.Replace(path, DataDirBase(), prefix, 1)
}
