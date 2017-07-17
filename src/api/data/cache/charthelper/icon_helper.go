package charthelper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/disintegration/imaging"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

const originalFormat = "original"

// IconOutput defines the output used by the consumer
type IconOutput struct {
	Name string
	Path string
}

// Strategy resizing method type
type strategyType uint8

/*
fit: scale down srcImage to fit the bounding box
fill: resize and crop the srcImage to fill the box
*/
const (
	fit strategyType = iota
	fill
)

type iconFormat struct {
	Name     string
	Width    int
	Height   int
	Strategy strategyType
}

var availableFormats = []iconFormat{
	{"160x160-fit", 160, 160, fit},
}

// DownloadAndProcessChartIcon the chart icon and process it
var DownloadAndProcessChartIcon = func(chart *models.ChartPackage) error {
	// There is no icon to download
	if chart.Icon == "" {
		return nil
	}

	if err := ensureChartDataDir(chart); err != nil {
		return err
	}

	if err := downloadIcon(chart); err != nil {
		return err
	}

	if err := processIcon(chart); err != nil {
		return err
	}

	return nil
}

// Downloads the icon and saves it in [chartDataDir]/original.[ext]
func downloadIcon(chart *models.ChartPackage) error {
	// Chart does not have an icon
	if chart.Icon == "" {
		return fmt.Errorf("Chart %s does not have icon property", *chart.Name)
	}

	// Use the same extension
	dest, _ := iconPath(chart, originalFormat)

	log.WithFields(log.Fields{
		"source": chart.Icon,
		"dest":   dest,
	}).Info("Downloading icon")

	// Download
	c := &http.Client{
		Timeout: defaultTimeout,
	}
	resp, err := c.Get(chart.Icon)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error downloading %s, %d\n", chart.Icon, resp.StatusCode)
	}

	// Create output
	out, err := os.Create(dest)
	defer out.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Process icons
func processIcon(chart *models.ChartPackage) error {
	exist, _ := iconExist(chart, originalFormat)
	if !exist {
		return fmt.Errorf("Can't find icon to process")
	}

	origPath, _ := iconPath(chart, originalFormat)
	for _, format := range availableFormats {
		destPath, _ := iconPath(chart, format.Name)
		log.WithFields(log.Fields{
			"source": origPath,
			"dest":   destPath,
		}).Info("Processing icon")
		orig, err := imaging.Open(origPath)
		if err != nil {
			return err
		}

		if format.Strategy == fill {
			processed := imaging.Fill(orig, format.Width, format.Height, imaging.Center, imaging.Lanczos)
			imaging.Save(processed, destPath)
		} else if format.Strategy == fit {
			processed := imaging.Fit(orig, format.Width, format.Height, imaging.Lanczos)
			imaging.Save(processed, destPath)
		}
	}
	return nil
}

func iconPath(chart *models.ChartPackage, format string) (string, error) {
	if chart.Icon == "" {
		return "", fmt.Errorf("Chart %s does not have icon property", *chart.Name)
	}
	return filepath.Join(chartDataDir(chart), fmt.Sprintf("logo-%s%s", format, filepath.Ext(chart.Icon))), nil
}

// Checks if the icon in specified format exist in chartDataDir
var iconExist = func(chart *models.ChartPackage, format string) (bool, error) {
	path, err := iconPath(chart, format)
	if err != nil {
		log.WithField("path", path).WithError(err).Error("IconExist error")
		return false, err
	}
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		log.WithField("path", path).Info("IconExist notFound")
		return false, nil
	}
	log.WithField("path", path).WithError(err).Error("IconExist error")
	return false, err
}

// AvailableIcons returns the list of supported icons and existing in the FS
var AvailableIcons = func(chart *models.ChartPackage, prefix string) []*IconOutput {
	var res []*IconOutput
	for _, format := range availableFormats {
		exist, err := iconExist(chart, format.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"chart":  *chart.Name,
				"format": format.Name,
			}).Error("Error on IconExists")
			continue
		}
		if !exist {
			continue
		}
		path, _ := iconPath(chart, format.Name)
		res = append(res, &IconOutput{
			Name: format.Name,
			Path: staticUrl(path, prefix),
		})
	}
	return res
}
