package chartpackagesort

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/mocks"
)

func TestSortedByName(t *testing.T) {
	chartsImplementation := mocks.NewMockCharts(mocks.MockedMethods{})
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	// Randomize slice before sorting
	for i := range charts {
		j := rand.Intn(i + 1)
		charts[i], charts[j] = charts[j], charts[i]
	}
	sort.Sort(ByName(charts))
	for i := 0; i < len(charts)-1; i++ {
		if *charts[i].Name > *charts[i+1].Name {
			t.Errorf("Array not sorted by name %s > %s", *charts[i].Name, *charts[i+1].Name)
		}
	}
}

func TestSortedBySemver(t *testing.T) {
	chartsImplementation := mocks.NewMockCharts(mocks.MockedMethods{})
	charts, err := chartsImplementation.All()
	chart := charts[0]
	assert.NoErr(t, err)
	versions, err := chartsImplementation.ChartVersionsFromRepo(chart.Repo, *chart.Name)
	assert.NoErr(t, err)
	// Randomize slice before sorting
	for i := range versions {
		j := rand.Intn(i + 1)
		versions[i], versions[j] = versions[j], versions[i]
	}
	sort.Sort(BySemver(versions))
	for i := 0; i < len(versions)-1; i++ {
		v1, _ := semver.NewVersion(*versions[i].Version)
		v2, _ := semver.NewVersion(*versions[i+1].Version)
		if v2.LessThan(v1) {
			t.Errorf("Array not sorted by semver %s > %s", *versions[i].Version, *versions[i+1].Version)
		}
	}
}

// If it is not a valid semver, it still sorts
func TestSortedBySemverWrongVersion(t *testing.T) {
	chartsImplementation := mocks.NewMockCharts(mocks.MockedMethods{})
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	// Bogus versions
	*charts[0].Version = "not-valid"
	sort.Sort(BySemver(charts))
	*charts[len(charts)-1].Version = "not-valid"
	sort.Sort(BySemver(charts))
}
