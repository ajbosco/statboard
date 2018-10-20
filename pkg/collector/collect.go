package collector

// Collector collects a metric data point
type Collector interface {
	Collect(metricName string, daysBack int) ([]Metric, error)
}
