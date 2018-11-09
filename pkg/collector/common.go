package collector

import (
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
)

func generateEmptyMetrics(metricName string, start time.Time, end time.Time) []statboard.Metric {
	var metrics []statboard.Metric

	// Generate statboard.Metric for each month in range
	start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour)
	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 1, 0) {
		met := statboard.Metric{Date: d, Name: metricName, Value: 0}
		metrics = append(metrics, met)
	}

	return metrics
}
