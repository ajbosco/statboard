package metric

import "time"

// Metric contains information about metrics
type Metric struct {
	ID    string `storm:"id"`
	Name  string
	Date  time.Time
	Value float64
}
