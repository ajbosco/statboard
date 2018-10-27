package collector

import (
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
)

func generateEmptyMetrics(metricName string, start time.Time, end time.Time) []statboard.Metric {
	var metrics []statboard.Metric

	// Generate statboard.Metric for each date in range
	for d := start.Truncate(24 * time.Hour); d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		met := statboard.Metric{Date: d, Name: metricName, Value: 0}
		metrics = append(metrics, met)
	}

	return metrics
}
