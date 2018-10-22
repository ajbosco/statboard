package statboard

import "time"

// Metric contains information about each metric
type Metric struct {
	ID    string `storm:"id"`
	Name  string
	Date  time.Time
	Value float64
}
