package main

import "time"

// Point is a time-series data point
type Point struct {
	timestamp time.Time
	series    string
	value     float64
}
