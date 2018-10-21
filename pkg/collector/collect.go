package collector

import "github.com/ajbosco/statboard/pkg/metric"

// Collector collects a metric data point
type Collector interface {
	Collect(metricName string, daysBack int) ([]metric.Metric, error)
}
