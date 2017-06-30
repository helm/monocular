package chartpackagesort

import (
	"github.com/Masterminds/semver"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

/*
This interface implementation will be
extended to support sorting by repostory first
*/

// ByName sorts by Name property
type ByName []*models.ChartPackage

// Sorting by name
func (c ByName) Len() int {
	return len(c)
}
func (c ByName) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ByName) Less(i, j int) bool {
	return *c[i].Name < *c[j].Name
}

// BySemver sorts by Semantic Versioning
type BySemver []*models.ChartPackage

// Sorting by name
func (c BySemver) Len() int {
	return len(c)
}
func (c BySemver) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c BySemver) Less(i, j int) bool {
	vi, err := semver.NewVersion(*c[i].Version)
	if err != nil {
		return true
	}
	vj, err := semver.NewVersion(*c[j].Version)
	if err != nil {
		return true
	}
	return vi.LessThan(vj)
}
