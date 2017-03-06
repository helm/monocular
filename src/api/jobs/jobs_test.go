package jobs

import (
	"errors"
	"testing"
	"time"

	"github.com/arschles/assert"
)

type testPeriodic struct {
	t        *testing.T
	err      error
	i        int
	freq     time.Duration
	name     string
	firstRun bool
}

func (t *testPeriodic) Do() error {
	t.t.Logf("testPeriodic Do at %s", time.Now())
	t.i++
	return t.err
}

func (t testPeriodic) Frequency() time.Duration {
	return t.freq
}

func (t testPeriodic) Name() string {
	return t.name
}

func (t testPeriodic) FirstRun() bool {
	return t.firstRun
}

func TestDoPeriodic(t *testing.T) {
	interval := time.Duration(500) * time.Millisecond
	p := &testPeriodic{t: t, err: nil, freq: interval, name: "test-periodic-job", firstRun: true}
	canceller := DoPeriodic([]Periodic{p})
	time.Sleep(interval / 2) // wait a little while for the goroutine to call the job once
	assert.True(t, p.i == 1, "the periodic wasn't called once")
	time.Sleep(interval)
	assert.True(t, p.i == 2, "the periodic wasn't called twice")
	time.Sleep(interval)
	assert.True(t, p.i == 3, "the periodic wasn't called thrice")
	canceller()
}

func TestDoPeriodicWithError(t *testing.T) {
	interval := time.Duration(500) * time.Millisecond
	p := &testPeriodic{t: t, err: errors.New("Do() crashes"), freq: interval, name: "test-periodic-job", firstRun: true}
	canceller := DoPeriodic([]Periodic{p})
	time.Sleep(interval / 2) // wait a little while for the goroutine to call the job once
	assert.True(t, p.i == 1, "the periodic wasn't called once")
	time.Sleep(interval)
	assert.True(t, p.i == 2, "the periodic wasn't called twice")
	canceller()
}

func TestDoPeriodicNoFirstRun(t *testing.T) {
	interval := time.Duration(500) * time.Millisecond
	p := &testPeriodic{t: t, err: nil, freq: interval, name: "test-periodic-job2", firstRun: false}
	canceller := DoPeriodic([]Periodic{p})
	time.Sleep(interval / 2) // wait a little while for the goroutine to call the job once
	assert.Equal(t, p.i, 0, "the periodic has not being called")
	time.Sleep(interval)
	assert.Equal(t, p.i, 1, "the periodic has being called once")
	canceller()
}
