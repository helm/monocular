package helpers

import (
	"fmt"
	"log"
	"strings"

	"github.com/helm/monocular/src/api/pkg/swagger/models"
	"gopkg.in/yaml.v2"
)

// IsYAML checks for valid YAML
func IsYAML(b []byte) bool {
	var yml map[string]interface{}
	return yaml.Unmarshal(b, &yml) == nil
}

// ParseYAMLRepo converts a YAML representation of a repo
// to a slice of versioned charts
func ParseYAMLRepo(rawYAML []byte) ([]models.ChartVersion, error) {
	repo := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(rawYAML, &repo); err != nil {
		return nil, err
	}
	var charts []models.ChartVersion
	for chartVersion := range repo {
		cV := repo[chartVersion]
		c, err := yaml.Marshal(&cV)
		if err != nil {
			log.Fatalf("couldn't parse repo chart: %v", err)
		}
		var chart models.ChartVersion
		if err := yaml.Unmarshal(c, &chart); err != nil {
			log.Fatalf("couldn't parse repo chart: %v", err)
		}
		charts = append(charts, chart)
	}
	return charts, nil
}

// ParseYAMLChartVersion converts a YAML representation of a versioned chart
// to a ChartVersion type
func ParseYAMLChartVersion(rawYAML []byte) (models.ChartVersion, error) {
	var chart models.ChartVersion
	if err := yaml.Unmarshal(rawYAML, &chart); err != nil {
		return models.ChartVersion{}, err
	}
	return chart, nil
}

// MakeChartResource composes a Resource type that represents a repo+chart
func MakeChartResource(chart models.ChartVersion, repo, version string) models.Resource {
	var ret models.Resource
	ret.Type = "chart"
	ret.ID = fmt.Sprintf("stable/%s", chart.Name)
	ret.Links = &models.ChartResourceLinks{
		Latest: fmt.Sprintf("/v1/charts/%s/%s/%s", repo, chart.Name, version),
	}
	ret.Attributes = &models.ChartResourceAttributes{
		Repo:        repo,
		Name:        chart.Name,
		Description: chart.Description,
		Created:     chart.Created,
		Home:        chart.Home,
	}
	return ret
}

// GetLatestChartVersion returns the most recent version from a slice of versioned charts
func GetLatestChartVersion(charts []models.ChartVersion, name string) (models.ChartVersion, error) {
	var latest string
	var ret models.ChartVersion
	for _, chart := range charts {
		if chart.Name == name {
			if latest == "" {
				latest = chart.Version
				ret = chart
			} else {
				newest, err := newestSemVer(latest, chart.Version)
				if err != nil {
					return models.ChartVersion{}, err
				}
				latest = newest
				if latest == chart.Version {
					ret = chart
				}
			}
		}
	}
	if latest == "" {
		return ret, fmt.Errorf("unable to determine latest version")
	}
	return ret, nil
}

// newestSemVer returns the newest (largest) semver string
func newestSemVer(v1 string, v2 string) (string, error) {
	v1Slice := strings.Split(v1, ".")
	v2Slice := strings.Split(v2, ".")
	for i, subVer1 := range v1Slice {
		if v2Slice[i] > subVer1 {
			return v2, nil
		} else if subVer1 > v2Slice[i] {
			return v1, nil
		}
	}
	return v1, nil
}
