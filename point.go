package main

import "time"

// Point is a time-series data point
type Point struct {
	series    string
	timestamp time.Time
	value     float64
}

// TimeSeries representation
type TimeSeries struct {
	name string
	tss  []time.Time
	xss  []float64
}

// NewTS creates new time-series object
func NewTS(p Point) *TimeSeries {
	var ts TimeSeries
	ts.name = p.series
	ts.Append(p.timestamp, p.value)
	return &ts
}

// Append value to time-series
func (t *TimeSeries) Append(ts time.Time, value float64) {
	t.tss = append(t.tss, ts)
	t.xss = append(t.xss, value)
}

// TimeSeriesStorage is an in-memory time-series storage
type TimeSeriesStorage struct {
	tss map[string]*TimeSeries
}

// NewTSS creates new time-series storage
func NewTSS() *TimeSeriesStorage {
	var tss TimeSeriesStorage
	tss.tss = make(map[string]*TimeSeries)
	return &tss
}

// Append datapoint to storage
func (t *TimeSeriesStorage) Append(p Point) {
	if tss, ok := t.tss[p.series]; ok {
		tss.Append(p.timestamp, p.value)
	} else {
		ts := NewTS(p)
		t.tss[p.series] = ts
	}
}
