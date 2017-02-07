package cache

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/arschles/assert"
)

func TestSortedByName(t *testing.T) {
	chartsImplementation := getChartsImplementation()
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
