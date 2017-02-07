package cache

import "github.com/helm/monocular/src/api/swagger/models"

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
