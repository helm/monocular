package jobs

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// Periodic is an interface for managing periodic job invocation
type Periodic interface {
	// Do begins the periodic job. It starts the first execution of the job, and then is
	// responsible for executing it every Frequency() thereafter
	Do() error
	Frequency() time.Duration
	Name() string
	FirstRun() bool
}

// PeriodicCanceller will cancel one or more Periodic jobs
type PeriodicCanceller func()

// DoPeriodic calls p.Do() once, and then again every p.Frequency() on each element p in pSlice.
// For each p in pSlice, a new goroutine is started, and the returned channel can be closed
// to stop all of the goroutines.
func DoPeriodic(pSlice []Periodic) PeriodicCanceller {
	doneCh := make(chan struct{})
	for _, p := range pSlice {
		go func(p Periodic) {
			// execute once at the beginning
			if p.FirstRun() {
				err := p.Do()
				if err != nil {
					log.Printf("periodic job ran and returned error (%s)", err)
				}
				log.Printf("periodic job %s ran", p.Name())
			}
			ticker := time.NewTicker(p.Frequency())
			for {
				select {
				case <-ticker.C:
					err := p.Do()
					if err != nil {
						log.Printf("periodic job ran and returned error (%s)", err)
					}
				case <-doneCh:
					ticker.Stop()
					return
				}
			}
		}(p)
	}
	return func() { close(doneCh) }
}
