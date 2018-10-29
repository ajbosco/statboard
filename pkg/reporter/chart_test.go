package reporter

import (
	"testing"
	"time"

	chartjs "github.com/ajbosco/goChartjs"
	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/stretchr/testify/assert"
)

func TestChartMetricToPoints(t *testing.T) {
	testDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(24 * time.Hour)

	tt := []struct {
		name     string
		metrics  []statboard.Metric
		expected []chartjs.Point
	}{
		{
			name: "metric to point",
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testDate,
					Value: 0.0,
				},
			},
			expected: []chartjs.Point{
				{
					X: testDate.Format("02-Jan-2006"),
					Y: 0.0,
				},
			},
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			actual := metricsToPoints(ts.metrics)
			assert.Equal(t, ts.expected, actual)
		})
	}
}

func TestChartGetChart(t *testing.T) {
	testDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(24 * time.Hour)
	testLineTension := 0
	testPoints := []chartjs.Point{
		{
			X: testDate.Format("02-Jan-2006"),
			Y: 0.0,
		},
	}

	testDataset := []chartjs.Dataset{{
		Label:           "testMetric",
		LineTension:     &testLineTension,
		Data:            testPoints,
		BackgroundColor: "testColor",
	},
	}

	expected := chartjs.Chart{
		Name:      "testMetric",
		ChartType: "line",
		Options: &chartjs.Options{
			Scales: chartjs.Scales{
				XAxes: []chartjs.Axes{
					{
						Time: &chartjs.Time{
							Unit: "day",
						},
						Type: "time",
					},
				},
			},
			MaintainAspectRatio: chartjs.False,
			Responsive:          chartjs.True,
		},
	}

	expected.Data.Datasets = testDataset

	actual := getChart(testPoints, "testMetric", "testColor")

	assert.Equal(t, expected, actual)
}

func TestChartNewChart(t *testing.T) {
	testDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(24 * time.Hour)
	testMetric := []statboard.Metric{
		{
			Name:  "testMetric",
			Date:  testDate,
			Value: 0.0,
		},
	}

	expected := Chart{
		metricName: "testMetric",
		ChartName:  "testChart",
		color:      "rgb(66,134,244)",
		metrics:    testMetric,
	}

	actual, err := NewChart("testMetric", "testChart", "#4286f4", testMetric)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestChartNewChart_InvalidColor(t *testing.T) {
	testDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(24 * time.Hour)
	testMetric := []statboard.Metric{
		{
			Name:  "testMetric",
			Date:  testDate,
			Value: 0.0,
		},
	}

	_, err := NewChart("testMetric", "testChart", "fakeColor", testMetric)
	assert.Error(t, err)
}
