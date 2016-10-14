package cache

import (
	"time"

	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/jobs"
)

type refreshChartsData struct {
	chartsImplementation data.Charts
	frequency            time.Duration
	name                 string
}

// NewRefreshChartsData creates a new Periodic implementation that refreshes charts data.
func NewRefreshChartsData(
	chartsImplementation data.Charts,
	frequency time.Duration,
	name string,
) jobs.Periodic {
	return &refreshChartsData{
		chartsImplementation: chartsImplementation,
		frequency:            frequency,
		name:                 name,
	}
}

// Do implements the Periodic interface
func (r *refreshChartsData) Do() error {
	if err := r.chartsImplementation.Refresh(); err != nil {
		return err
	}
	return nil
}

// Frequency implements the Periodic interface
func (r *refreshChartsData) Frequency() time.Duration {
	return r.frequency
}

// Name implements the Periodic interface
func (r *refreshChartsData) Name() string {
	return r.name
}
