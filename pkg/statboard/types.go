package statboard

import "time"

// Metric contains information about metrics
type Metric struct {
	Name  string
	Date  time.Time
	Value float64
}
