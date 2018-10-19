package collector

import "time"

type metric struct {
	Source string
	Date   time.Time
	Value  float64
}
