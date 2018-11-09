package collector

import "github.com/ajbosco/statboard/pkg/statboard"

// Collector collects a metric data point
type Collector interface {
	Collect(metricName string, monthsBack int) ([]statboard.Metric, error)
}
