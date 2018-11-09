package collector

import (
	"testing"
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/stretchr/testify/assert"
)

func TestFitbitAggregateEvents(t *testing.T) {
	testActivityDate := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)

	tt := []struct {
		name     string
		steps    []FitbitSteps
		metrics  []statboard.Metric
		expected []statboard.Metric
	}{
		{
			name: "single activity",
			steps: []FitbitSteps{
				{
					ActivityDate: "2018-01-01",
					Steps:        "100",
				},
			},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testActivityDate.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  time.Date(testActivityDate.Year(), testActivityDate.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testActivityDate.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  time.Date(testActivityDate.Year(), testActivityDate.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 100.0,
				},
			},
		},
		{
			name: "multiple events",
			steps: []FitbitSteps{
				{
					ActivityDate: "2018-01-01",
					Steps:        "100",
				},
				{
					ActivityDate: "2018-01-01",
					Steps:        "125",
				},
				{
					ActivityDate: "2018-01-01",
					Steps:        "50",
				},
			},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testActivityDate.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  time.Date(testActivityDate.Year(), testActivityDate.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testActivityDate.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  time.Date(testActivityDate.Year(), testActivityDate.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 275.0,
				},
			},
		},
		{
			name:  "no events",
			steps: []FitbitSteps{},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  time.Date(testActivityDate.Year(), testActivityDate.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  time.Date(testActivityDate.Year(), testActivityDate.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			actual, err := aggregateSteps(ts.steps, ts.metrics)
			assert.Equal(t, ts.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestFitbitCollect_InvalidMetric(t *testing.T) {
	c := FitbitCollector{}

	_, err := c.Collect("fake_metric", 1)
	assert.Error(t, err)
}
